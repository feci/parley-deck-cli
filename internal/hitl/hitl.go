package hitl

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/store"
)

const (
	RiskLow    = "low"
	RiskNormal = "normal"
	RiskHigh   = "high"

	StatusOpen         = "open"
	StatusAnswered     = "answered"
	StatusAutoAnswered = "auto_answered"
)

type Question struct {
	ID            string    `json:"id"`
	Agent         string    `json:"agent"`
	Prompt        string    `json:"prompt"`
	Details       string    `json:"details"`
	DefaultAnswer string    `json:"default_answer"`
	Risk          string    `json:"risk"`
	Status        string    `json:"status"`
	Answer        string    `json:"answer"`
	CreatedAt     time.Time `json:"created_at"`
	AnsweredAt    time.Time `json:"answered_at"`
}

type Store struct {
	runDir string
}

func New(runDir string) Store {
	return Store{runDir: runDir}
}

func (s Store) QuestionsDir() string {
	return filepath.Join(s.runDir, "questions")
}

func (s Store) Create(question Question) (Question, error) {
	now := time.Now().UTC()
	if question.Agent == "" {
		question.Agent = "runner"
	}
	if question.ID == "" {
		question.ID = NewQuestionID(question.Agent, now)
	}
	if question.Risk == "" {
		question.Risk = RiskNormal
	}
	if question.Status == "" {
		question.Status = StatusOpen
	}
	if question.CreatedAt.IsZero() {
		question.CreatedAt = now
	}
	if question.Prompt == "" {
		return Question{}, errors.New("question prompt is required")
	}
	if err := s.write(question); err != nil {
		return Question{}, err
	}
	if err := s.appendQuestionEvent(question); err != nil {
		return Question{}, err
	}
	return question, nil
}

func (s Store) List() ([]Question, error) {
	entries, err := os.ReadDir(s.QuestionsDir())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var questions []Question
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		question, err := s.read(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	sort.SliceStable(questions, func(i, j int) bool {
		if !questions[i].CreatedAt.Equal(questions[j].CreatedAt) {
			return questions[i].CreatedAt.Before(questions[j].CreatedAt)
		}
		return questions[i].ID < questions[j].ID
	})
	return questions, nil
}

func (s Store) Answer(id, answer string, auto bool) (Question, error) {
	question, err := s.read(id)
	if err != nil {
		return Question{}, err
	}
	if question.Status != StatusOpen {
		return question, nil
	}
	if auto {
		if question.Risk != RiskLow || strings.TrimSpace(question.DefaultAnswer) == "" {
			return Question{}, fmt.Errorf("question %s cannot be auto-answered", id)
		}
		answer = question.DefaultAnswer
		question.Status = StatusAutoAnswered
	} else {
		if strings.TrimSpace(answer) == "" {
			return Question{}, errors.New("answer is required")
		}
		question.Status = StatusAnswered
	}
	question.Answer = answer
	question.AnsweredAt = time.Now().UTC()
	if err := s.write(question); err != nil {
		return Question{}, err
	}
	if err := s.appendAnsweredEvent(question); err != nil {
		return Question{}, err
	}
	return question, nil
}

func (s Store) AutoAnswerOpen() ([]Question, error) {
	questions, err := s.List()
	if err != nil {
		return nil, err
	}
	var answered []Question
	for _, question := range questions {
		if question.Status != StatusOpen || question.Risk != RiskLow || strings.TrimSpace(question.DefaultAnswer) == "" {
			continue
		}
		updated, err := s.Answer(question.ID, question.DefaultAnswer, true)
		if err != nil {
			return nil, err
		}
		answered = append(answered, updated)
	}
	return answered, nil
}

func (s Store) read(id string) (Question, error) {
	path, err := s.questionPath(id)
	if err != nil {
		return Question{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Question{}, err
	}
	var question Question
	if err := json.Unmarshal(data, &question); err != nil {
		return Question{}, err
	}
	return question, nil
}

func (s Store) write(question Question) error {
	path, err := s.questionPath(question.ID)
	if err != nil {
		return err
	}
	if err := fsutil.MkdirAllResilient(s.QuestionsDir(), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(question, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(s.QuestionsDir(), "."+question.ID+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

func (s Store) questionPath(id string) (string, error) {
	if id == "" || id == "." || id == ".." || strings.ContainsAny(id, `/\`) {
		return "", fmt.Errorf("invalid question id %q", id)
	}
	return filepath.Join(s.QuestionsDir(), id+".json"), nil
}

func (s Store) appendQuestionEvent(question Question) error {
	return store.New(s.runDir).Append(store.Event{
		Time: time.Now().UTC(),
		Type: "hitl.question",
		Data: map[string]any{
			"question_id": question.ID,
			"agent":       question.Agent,
			"risk":        question.Risk,
			"status":      question.Status,
		},
	})
}

func (s Store) appendAnsweredEvent(question Question) error {
	return store.New(s.runDir).Append(store.Event{
		Time: time.Now().UTC(),
		Type: "hitl.answered",
		Data: map[string]any{
			"question_id": question.ID,
			"agent":       question.Agent,
			"status":      question.Status,
		},
	})
}

func NewQuestionID(agent string, now time.Time) string {
	return fmt.Sprintf("%s-%s-%s", now.UTC().Format("20060102T150405.000000000Z"), slugify(agent), randomSuffix())
}

func randomSuffix() string {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "000000"
	}
	return hex.EncodeToString(b[:])
}

func slugify(value string) string {
	value = strings.ToLower(value)
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case !lastDash:
			b.WriteByte('-')
			lastDash = true
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "agent"
	}
	return slug
}
