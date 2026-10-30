package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"parley-deck-cli/internal/fsutil"
	"parley-deck-cli/internal/pidlease"
	"parley-deck-cli/internal/quota"
)

const CaptureFileLimit = 256 << 10
const CaptureTotalLimit = 2 << 20
const CaptureRetention = 7 * 24 * time.Hour

type captureSpan struct{ start, end int }

var captureName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{0,100}\.log$`)
var captureGateway = regexp.MustCompile(`^All [A-Za-z0-9._-]{1,32} accounts have exhausted their quota \(reset after [0-9]+[hms]\)$`)

// ScrubZcodeCapture keeps exact safe framing and marks every removed span.
// Unknown or incomplete untruncated grammar writes nothing.
func ScrubZcodeCapture(raw []byte, truncated bool) ([]byte, error) {
	text := string(raw)
	p := &sdkDump{raw: strings.TrimRight(text, "\r\n"), capture: true}
	for {
		p.space()
		if p.pos >= len(p.raw) {
			return nil, fmt.Errorf("capture lacks terminal framing")
		}
		if zcodeTerminal.MatchString(p.raw[p.pos:]) {
			break
		}
		if _, err := p.failure(0); err != nil {
			if truncated {
				lines := strings.Split(strings.TrimRight(text, "\r\n"), "\n")
				last := lines[len(lines)-1]
				if len(last) <= 200 && zcodeTerminal.MatchString(last) {
					return []byte("[TRUNCATED]\n[REDACTED]\n" + last + "\n"), nil
				}
			}
			return nil, fmt.Errorf("capture scrub could not prove framing")
		}
	}
	if len(p.records) == 0 {
		return nil, fmt.Errorf("capture lacks provider record")
	}
	for _, r := range p.records {
		if !zcodeExhaustionMessage(r.message) && !captureGateway.MatchString(r.message) {
			return nil, fmt.Errorf("capture message cannot be safely retained")
		}
		var envelope map[string]json.RawMessage
		var fields map[string]json.RawMessage
		if json.Unmarshal([]byte(r.body), &envelope) != nil || len(envelope) != 1 || json.Unmarshal(envelope["error"], &fields) != nil {
			return nil, fmt.Errorf("capture body cannot be safely retained")
		}
		for k, v := range fields {
			switch k {
			case "message":
			case "reset_at":
				var value string
				if json.Unmarshal(v, &value) != nil {
					return nil, fmt.Errorf("unsafe reset")
				}
				if _, e := time.Parse(time.RFC3339Nano, value); e != nil {
					return nil, fmt.Errorf("unsafe reset")
				}
			case "retry_after", "status", "statusCode", "code":
				if !regexp.MustCompile(`^"?[0-9]{1,12}"?$`).Match(v) {
					return nil, fmt.Errorf("unsafe diagnostic scalar")
				}
			default:
				return nil, fmt.Errorf("unknown diagnostic body field")
			}
		}
	}
	sort.Slice(p.redactions, func(i, j int) bool { return p.redactions[i].start > p.redactions[j].start })
	for _, span := range p.redactions {
		text = text[:span.start] + "[REDACTED]" + text[span.end:]
	}
	if truncated {
		text = "[TRUNCATED]\n" + text
	}
	if len(text) > CaptureFileLimit {
		text = "[TRUNCATED]\n" + text[len(text)-(CaptureFileLimit-len("[TRUNCATED]\n")):]
	}
	return []byte(text), nil
}

type CaptureInventory struct {
	Files           int   `json:"files"`
	Bytes           int64 `json:"bytes"`
	ExpiredRemoved  int   `json:"expired_removed"`
	CapacityRemoved int   `json:"capacity_removed"`
	Purged          int   `json:"purged"`
}

func captureDirectory(dir string, create bool) error {
	for path := dir; path != filepath.Dir(path); path = filepath.Dir(path) {
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return fmt.Errorf("unsafe capture directory")
		}
	}
	if create {
		if err := os.MkdirAll(filepath.Dir(dir), 0700); err != nil {
			return err
		}
		if err := os.Mkdir(dir, 0700); err != nil && !os.IsExist(err) {
			return err
		}
	}
	for _, path := range []string{filepath.Dir(dir), dir} {
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if info.Mode().Perm() != 0700 {
			return fmt.Errorf("capture directories require mode 0700")
		}
	}
	return nil
}

func captureFiles(dir string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []os.FileInfo
	for _, e := range entries {
		if e.Name() == "capture.lease" || e.Name() == "cleanup.json" {
			info, err := os.Lstat(filepath.Join(dir, e.Name()))
			if err != nil {
				return nil, err
			}
			if !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
				return nil, fmt.Errorf("unsafe capture control file")
			}
			continue
		}
		if !captureName.MatchString(e.Name()) {
			return nil, fmt.Errorf("unknown file in capture directory")
		}
		info, err := os.Lstat(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
			return nil, fmt.Errorf("unsafe capture file")
		}
		files = append(files, info)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].ModTime().Before(files[j].ModTime()) })
	return files, nil
}

func captureCleanup(dir string, now time.Time, need int, purge bool) (CaptureInventory, error) {
	files, err := captureFiles(dir)
	if err != nil {
		return CaptureInventory{}, err
	}
	var result CaptureInventory
	var retained []os.FileInfo
	for _, f := range files {
		if purge || now.Sub(f.ModTime()) > CaptureRetention {
			if err = os.Remove(filepath.Join(dir, f.Name())); err != nil {
				return result, err
			}
			if purge {
				result.Purged++
			} else {
				result.ExpiredRemoved++
			}
			continue
		}
		retained = append(retained, f)
		result.Bytes += f.Size()
	}
	for len(retained) > 0 && result.Bytes+int64(need)+2048 > CaptureTotalLimit {
		f := retained[0]
		retained = retained[1:]
		if err = os.Remove(filepath.Join(dir, f.Name())); err != nil {
			return result, err
		}
		result.Bytes -= f.Size()
		result.CapacityRemoved++
	}
	result.Files = len(retained)
	receipt, _ := json.Marshal(struct {
		Time                      time.Time `json:"time"`
		Expired, Capacity, Purged int
	}{now.UTC(), result.ExpiredRemoved, result.CapacityRemoved, result.Purged})
	path := filepath.Join(dir, "cleanup.json")
	if info, e := os.Lstat(path); e == nil && (!info.Mode().IsRegular() || info.Mode().Perm() != 0600) {
		return result, fmt.Errorf("unsafe cleanup receipt")
	} else if e != nil && !os.IsNotExist(e) {
		return result, e
	}
	if err = quota.DurableWrite(path, append(receipt, '\n'), false); err != nil {
		return result, err
	}
	return result, nil
}

func CaptureStatus(dir string, now time.Time, purge bool) (CaptureInventory, error) {
	if err := captureDirectory(dir, false); os.IsNotExist(err) {
		return CaptureInventory{}, nil
	} else if err != nil {
		return CaptureInventory{}, err
	}
	if _, err := captureFiles(dir); err != nil {
		return CaptureInventory{}, err
	}
	l, err := pidlease.TryAcquire(filepath.Join(dir, "capture.lease"), "zcode-capture")
	if err != nil {
		return CaptureInventory{}, err
	}
	defer l.Release()
	return captureCleanup(dir, now, 0, purge)
}

func CaptureZcodeFailure(dir, id string, raw []byte, truncated bool, now time.Time) (bool, error) {
	if !captureName.MatchString(id + ".log") {
		return false, fmt.Errorf("invalid capture identity")
	}
	scrubbed, err := ScrubZcodeCapture(raw, truncated)
	if err != nil {
		return false, err
	}
	if err = captureDirectory(dir, true); err != nil {
		return false, err
	}
	if _, err = captureFiles(dir); err != nil {
		return false, err
	}
	l, err := pidlease.TryAcquire(filepath.Join(dir, "capture.lease"), "zcode-capture")
	if err != nil {
		return false, err
	}
	defer l.Release()
	if _, err = captureCleanup(dir, now, len(scrubbed), false); err != nil {
		return false, err
	}
	path := filepath.Join(dir, id+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if os.IsExist(err) {
		old, e := os.ReadFile(path)
		if e == nil && bytes.Equal(old, scrubbed) {
			return false, nil
		}
		return false, fmt.Errorf("capture identity already exists")
	}
	if err != nil {
		return false, err
	}
	_, err = f.Write(scrubbed)
	if err == nil {
		err = fsutil.SyncFile(f)
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err == nil {
		err = quota.SyncPath(dir)
	}
	return err == nil, err
}
