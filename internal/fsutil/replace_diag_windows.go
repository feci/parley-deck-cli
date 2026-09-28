//go:build windows

package fsutil

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// TEMPORARY replace-failure diagnostic (kimi-1 corrected battery, claude-1
// recurrence consult; FINAL §D.6 untouched). REMOVAL PLAN: this batch is
// declared temporary — once the concurrent-replacement class is identified
// from hosted evidence, the probes are removed (or narrowed per review) and
// only the structured errno in the wrapped error may remain. Everything here
// is FAILURE-PATH-ONLY (unreachable when MoveFileEx succeeds — the probes
// are themselves opens, exactly what the row-30 work removed from the
// success path), REPORT-NEVER-ACT (no retry, no fallback, no swallow: the
// wrapped error still fails the operation), and GRADE-NOT-CLASSIFY: these
// are probe-time samples at t1 > t0 — a succeeding probe does not clear the
// failure-time holder, a failing probe does not prove the t0 cause, and no
// sample licenses any retry or §D.6 amendment.

// DiagInFlightReaders is the per-path in-flight reader registry (D2′): the
// budget reader increments on entry and decrements on exit, keyed by cleaned
// path so a policy.json rename block is not conflated with readers of other
// files. Sampled by the replace failure path BEFORE the probes (closest to
// t0). Product readers only; a false zero is possible when a t0 reader
// closed before the sample.
var DiagInFlightReaders = struct {
	counts atomic.Pointer[map[string]int]
}{}
var diagReadersInit = func() bool {
	m := map[string]int{}
	DiagInFlightReaders.counts.Store(&m)
	return true
}()

// DiagReaderEnter registers an in-flight read of path.
func DiagReaderEnter(path string, at time.Time) {
	m := *DiagInFlightReaders.counts.Load()
	n := map[string]int{}
	for k, v := range m {
		n[k] = v
	}
	n[path]++
	DiagInFlightReaders.counts.Store(&n)
	_ = at
}

// DiagReaderExit unregisters an in-flight read of path.
func DiagReaderExit(path string) {
	m := *DiagInFlightReaders.counts.Load()
	n := map[string]int{}
	for k, v := range m {
		n[k] = v
	}
	if n[path] > 0 {
		n[path]--
	}
	DiagInFlightReaders.counts.Store(&n)
}

// diagReaderSample reports the in-flight count for path at sample time.
func diagReaderSample(path string) int {
	return (*DiagInFlightReaders.counts.Load())[path]
}

// replaceDiag grades a MoveFileEx failure (t1 sampling; grading only):
// raw errno, per-side existence+attributes (P1), a DELETE-access full-share
// probe (P2 — fails only when some handle denies delete-sharing, the exact
// rename-blocking condition, or delete-pending), a READ_ATTRIBUTES full-share
// control (P3 — fails on delete-pending/hard deny), and the D2′ reader
// sample. Unclassified outcomes are reported as unclassified.
func replaceDiag(op string, staged, target string, errno error) string {
	at := time.Now().UnixNano()
	sample := diagReaderSample(target)
	p := func(side, path string) string {
		var attrWord uint32
		exists := "yes"
		attrsPtr, err := windows.UTF16PtrFromString(path)
		if err != nil {
			return fmt.Sprintf("%s: path-error %v", side, err)
		}
		var data windows.Win32FileAttributeData
		if err := windows.GetFileAttributesEx(attrsPtr, windows.GetFileExInfoStandard, (*byte)(unsafe.Pointer(&data))); err != nil {
			if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) || errors.Is(err, windows.ERROR_PATH_NOT_FOUND) {
				exists = "no"
			} else {
				exists = "stat-err:" + err.Error()
			}
			attrWord = 0
		} else {
			attrWord = data.FileAttributes
		}
		// P2: DELETE access with FULL share — mirrors the rename's own
		// requirement; succeeds over benign share-delete readers.
		p2 := "ok"
		if h, err := windows.CreateFile(attrsPtr, windows.DELETE, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
			nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0); err != nil {
			p2 = fmt.Sprintf("errno-%d", errnoValue(err))
		} else {
			windows.CloseHandle(h)
		}
		// P3 control: READ_ATTRIBUTES with full share — fails on
		// delete-pending or a hard deny.
		p3 := "ok"
		if h, err := windows.CreateFile(attrsPtr, windows.FILE_READ_ATTRIBUTES, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
			nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0); err != nil {
			p3 = fmt.Sprintf("errno-%d", errnoValue(err))
		} else {
			windows.CloseHandle(h)
		}
		graded := "unclassified"
		if p2 != "ok" && p3 == "ok" {
			graded = "delete-sharing-denied-handle-present(t1)"
		} else if p2 != "ok" && p3 != "ok" {
			graded = "delete-pending-or-hard-deny(t1)"
		} else if p2 == "ok" && p3 == "ok" {
			graded = "no-handle-blocker-at-t1(filter-level|attribute-at-t0|gone-transient;unclassified)"
		}
		return fmt.Sprintf("%s: exists=%s attrs=%#x readonly-bit=%t P2=%s P3=%s grade=%s",
			side, exists, attrWord, attrWord&windows.FILE_ATTRIBUTE_READONLY != 0, p2, p3, graded)
	}
	return fmt.Sprintf("[replace-diag t=%d op=%s errno=%d readers(%s)=%d | %s | %s]",
		at, op, errnoValue(errno), target, sample, p("staged", staged), p("target", target))
}

func errnoValue(err error) uint32 {
	var e windows.Errno
	if errors.As(err, &e) {
		return uint32(e)
	}
	return 0
}
