// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!nacl,!plan9

package log

import (
	"io"
	"log/syslog"
	"os"
	"sync"
)

// Syslog is the backend using syslog
type Syslog struct {
	out *syslog.Writer
	mu  sync.Mutex
}

// NewSyslog crete a new logger using syslog backend with level/prefix/flag
func NewSyslog(level Level, prefix string, flag int) *Logger {

	// clear some flags not needed since syslog already provided
	flag &= ^Ldate & ^Ltime &^ Lmicroseconds &^ LUTC

	backend, _ := NewSyslogBackend(level, prefix)

	l := &Logger{
		level:   level,
		mu:      sync.Mutex{},
		prefix:  prefix,
		flag:    flag,
		backend: backend,
		index:   0,
		name:    procName(),
	}

	logger.Store(l)

	return l
}

// NewSyslogBackend creates a syslog backend
func NewSyslogBackend(level Level, prefix string) (Backend, error) {
	syslogLevel := syslog.LOG_INFO
	switch level {
	case FATAL:
		syslogLevel = syslog.LOG_CRIT
	case ERROR:
		syslogLevel = syslog.LOG_ERR
	case WARN:
		syslogLevel = syslog.LOG_WARNING
	case INFO:
		syslogLevel = syslog.LOG_INFO
	case DEBUG:
		syslogLevel = syslog.LOG_DEBUG
	default:
		syslogLevel = syslog.LOG_NOTICE
	}

	out, err := syslog.New(syslogLevel, prefix)

	if err == nil {
		return &Syslog{
			out: out,
		}, nil
	}
	return nil, err
}

// Writer returns the io.Writer of current Syslog
func (l *Syslog) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// SetWriter set the io.Writer of current Syslog
func (l *Syslog) SetWriter(io.Writer) {
	// not supported
}

// log the actual write routine of logging
func (l *Syslog) log(r *Record) {

	// syslog writes log with newline, we don't need extra newline
	r.newline = false
	message := r.message()

	switch r.level {
	case FATAL:
		_ = l.out.Crit(message)
		os.Exit(-1)
	case ERROR:
		_ = l.out.Err(message)
	case WARN:
		_ = l.out.Warning(message)
	case INFO:
		_ = l.out.Info(message)
	case DEBUG:
		_ = l.out.Debug(message)
	default:
		_ = l.out.Notice(message)
	}
}

func (l *Syslog) start() {
}

// Flush the current log backend
func (l *Syslog) Flush() {
}

func (l *Syslog) isatty() bool {
	return false
}

func (l *Syslog) fd() Handler {
	return nil
}

func (l *Syslog) write([]byte) error {
	return nil
}
