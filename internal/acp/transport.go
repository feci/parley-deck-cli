package acp

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"sync"
)

// Transport reads and writes JSON-RPC 2.0 messages as newline-delimited JSON
// (NDJSON). Each direction is independent — Read blocks until a line arrives,
// Write serializes one message at a time.
type Transport struct {
	reader *bufio.Reader
	writer io.Writer
	writeM sync.Mutex
	// maxLineBytes caps the size of a single NDJSON line so a malformed
	// agent cannot exhaust memory by emitting an unterminated megabyte.
	maxLineBytes int
}

// NewTransport wraps an stdio pair. The reader is buffered with a 1 MiB
// initial buffer because ACP agents stream content chunks that can exceed
// the default 4 KiB bufio buffer.
func NewTransport(reader io.Reader, writer io.Writer) *Transport {
	const initialBuf = 1 << 20
	br := bufio.NewReaderSize(reader, initialBuf)
	return &Transport{reader: br, writer: writer, maxLineBytes: 8 << 20}
}

// Read returns the next decoded message or io.EOF when the agent closes
// stdout. Lines that fail to parse as JSON are skipped silently — AionUi
// does the same since some agents print non-JSON banners before opening
// their NDJSON stream.
func (t *Transport) Read() (Message, error) {
	for {
		line, err := t.readLine()
		if err != nil {
			return Message{}, err
		}
		if len(line) == 0 {
			continue
		}
		var msg Message
		if jsonErr := json.Unmarshal(line, &msg); jsonErr != nil {
			continue
		}
		return msg, nil
	}
}

func (t *Transport) readLine() ([]byte, error) {
	var collected []byte
	for {
		chunk, isPrefix, err := t.reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) && len(collected) > 0 {
				return collected, nil
			}
			return nil, err
		}
		collected = append(collected, chunk...)
		if !isPrefix {
			return collected, nil
		}
		if len(collected) > t.maxLineBytes {
			return nil, errors.New("acp: NDJSON line exceeded max size")
		}
	}
}

// Write serializes a message and appends a newline. Concurrent writers are
// serialized through a mutex so JSON-RPC framing stays intact.
func (t *Transport) Write(msg Message) error {
	t.writeM.Lock()
	defer t.writeM.Unlock()
	encoded, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	encoded = append(encoded, '\n')
	_, err = t.writer.Write(encoded)
	return err
}
