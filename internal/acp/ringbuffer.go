package acp

import "sync"

// RingBuffer is a fixed-capacity append-only sink that keeps only the last
// N bytes. Used to capture child-process stderr without unbounded growth.
// Concurrent writes are serialized.
type RingBuffer struct {
	mu  sync.Mutex
	cap int
	buf []byte
}

// NewRingBuffer constructs a buffer with the given byte capacity.
func NewRingBuffer(capacity int) *RingBuffer {
	if capacity <= 0 {
		capacity = 1
	}
	return &RingBuffer{cap: capacity}
}

// Write appends p, truncating the head when the capacity is exceeded.
func (r *RingBuffer) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf = append(r.buf, p...)
	if len(r.buf) > r.cap {
		r.buf = r.buf[len(r.buf)-r.cap:]
	}
	return len(p), nil
}

// String returns the current buffered contents as a string.
func (r *RingBuffer) String() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return string(r.buf)
}
