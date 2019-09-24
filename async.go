// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"io"
	"sync"
	"sync/atomic"
)

// AsyncLog an async logger
type AsyncLog struct {
	stopped uint32
	out     io.Writer
	mu      sync.Mutex
	handler Handler
	istty   bool
	queue   chan *Record
	stop    chan struct{} // Notify closing
}

// NewAsyncLogger creates a new async logger with a io.Writer
func NewAsyncLogger(w io.Writer, prefix string, flag int) *Logger {

	l := &Logger{
		level:   DEBUG,
		mu:      sync.Mutex{},
		prefix:  prefix,
		flag:    flag,
		backend: NewAsyncBackend(w),
		index:   0,
		name:    procName(),
	}

	go l.backend.start()
	logger.Store(l)
	return l
}

// NewAsyncBackend creates a new async backend
func NewAsyncBackend(w io.Writer) Backend {
	handler, ok := w.(Handler)
	return &AsyncLog{
		out:     w,
		handler: handler,
		istty:   ok && isatty(handler.Fd()),
		queue:   make(chan *Record, 1024),
		stop:    make(chan struct{}),
	}
}

// Writer returns the io.Writer of current Syslog
func (l *AsyncLog) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// SetWriter set the io.Writer of current logger
func (l *AsyncLog) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if fd, ok := w.(Handler); ok {
		l.handler = fd
		l.istty = isatty(fd.Fd())
	}
	l.out = w
}

func (l *AsyncLog) write(data []byte) error {
	_, err := l.out.Write(data)
	return err
}

func (l *AsyncLog) log(r *Record) {
	if atomic.LoadUint32(&l.stopped) == 0 {
		l.queue <- r
	}
}

func (l *AsyncLog) start() {
	for r := range l.queue {
		writeLog(l, r)
	}
	close(l.stop)
}

// Flush the current log backend
func (l *AsyncLog) Flush() {
	atomic.StoreUint32(&l.stopped, 1)
	close(l.queue)
	<-l.stop
}

func (l *AsyncLog) isatty() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.istty
}

func (l *AsyncLog) fd() Handler {
	return l.handler
}
