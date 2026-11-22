package compute

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"time"
)

type LogLine struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type RingBuffer struct {
	mu          sync.RWMutex
	maxLines    int
	lines       []LogLine
	subscribers map[chan LogLine]struct{}
}

func NewRingBuffer(maxLines int) *RingBuffer {
	if maxLines <= 0 {
		maxLines = 1000
	}
	return &RingBuffer{
		maxLines:    maxLines,
		lines:       make([]LogLine, 0, maxLines),
		subscribers: make(map[chan LogLine]struct{}),
	}
}

// Write implements io.Writer, scanning incoming bytes line-by-line.
func (rb *RingBuffer) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		rb.Append(scanner.Text())
	}
	return len(p), nil
}

// Append appends a log message and notifies live subscribers.
func (rb *RingBuffer) Append(message string) {
	line := LogLine{
		Timestamp: time.Now().UTC(),
		Message:   message,
	}

	rb.mu.Lock()
	if len(rb.lines) >= rb.maxLines {
		// Evict oldest (ring slide)
		copy(rb.lines, rb.lines[1:])
		rb.lines[len(rb.lines)-1] = line
	} else {
		rb.lines = append(rb.lines, line)
	}

	// Fan out to active subscribers
	for ch := range rb.subscribers {
		select {
		case ch <- line:
		default:
			// Non-blocking drop if consumer is slow to prevent stalling the writer
		}
	}
	rb.mu.Unlock()
}

// GetTail returns up to `n` recent log lines.
func (rb *RingBuffer) GetTail(n int) []LogLine {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	total := len(rb.lines)
	if n <= 0 || n > total {
		n = total
	}

	start := total - n
	result := make([]LogLine, n)
	copy(result, rb.lines[start:])
	return result
}

// Subscribe registers a live stream channel.
func (rb *RingBuffer) Subscribe() (<-chan LogLine, func()) {
	ch := make(chan LogLine, 100)

	rb.mu.Lock()
	rb.subscribers[ch] = struct{}{}
	rb.mu.Unlock()

	unsubscribe := func() {
		rb.mu.Lock()
		delete(rb.subscribers, ch)
		close(ch)
		rb.mu.Unlock()
	}

	return ch, unsubscribe
}

var _ io.Writer = (*RingBuffer)(nil)
