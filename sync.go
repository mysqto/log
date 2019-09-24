// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"io"
	"sync"
)

// SyncLog write log synchronously
type SyncLog struct {
	out     io.Writer
	mu      sync.Mutex
	handler Handler
	istty   bool
}

// NewSyncBackend create a new sync backend
func NewSyncBackend(w io.Writer) Backend {
	handler, ok := w.(Handler)
	return &SyncLog{
		out:     w,
		handler: handler,
		istty:   ok && isatty(handler.Fd()),
	}
}

// Writer returns the io.Writer of current logger
func (l *SyncLog) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// SetWriter set the io.Writer of current logger
func (l *SyncLog) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if handler, ok := w.(Handler); ok {
		l.handler = handler
		l.istty = isatty(handler.Fd())
	}
	l.out = w
}

func (l *SyncLog) write(data []byte) error {
	_, err := l.out.Write(data)
	return err
}

func (l *SyncLog) log(r *Record) {
	writeLog(l, r)
}

func (l *SyncLog) start() {
}

// Flush the current log backend
func (l *SyncLog) Flush() {
	// nothing to do with current logger
}

func (l *SyncLog) isatty() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.istty
}

func (l *SyncLog) fd() Handler {
	return l.handler
}
