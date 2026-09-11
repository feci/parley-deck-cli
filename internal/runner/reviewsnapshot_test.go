package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReviewSnapshotMoveBackPreservesRecoveryOnPublicationFailure(t *testing.T) {
	snapshot := &ReviewSnapshot{Dir: t.TempDir()}
	body := []byte("reviewer's recovery copy\n")
	if err := os.WriteFile(filepath.Join(snapshot.Dir, "review.md"), body, 0600); err != nil {
		t.Fatal(err)
	}
	live := t.TempDir()
	// A directory cannot be replaced by an artifact. Its existing contents
	// and the snapshot must survive the failed publication.
	destination := filepath.Join(live, "review.md")
	if err := os.Mkdir(destination, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destination, "prior"), []byte("retained"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := snapshot.MoveArtifactBack("review.md", destination); err == nil {
		t.Fatal("invalid replacement passed")
	}
	if got, err := os.ReadFile(filepath.Join(destination, "prior")); err != nil || string(got) != "retained" {
		t.Fatalf("old destination lost: %s %v", got, err)
	}
	if got, err := os.ReadFile(filepath.Join(snapshot.Dir, "review.md")); err != nil || string(got) != string(body) {
		t.Fatalf("recovery copy lost: %s %v", got, err)
	}
	files, err := os.ReadDir(live)
	if err != nil || len(files) != 1 {
		t.Fatalf("publication scratch leaked: %v %v", files, err)
	}
}
