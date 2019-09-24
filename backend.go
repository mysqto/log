// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"io"
	"os"
)

// Handler a file handler interface
type Handler interface {
	Fd() uintptr
}

// Backend a logger must implement the backend interface
type Backend interface {
	// Writer returns the io.Writer of current backend for compatible with golang's log package
	Writer() io.Writer

	// SetWriter sets the io.Writer of current backend for compatible with golang's log package
	SetWriter(w io.Writer)

	// Flush flushes current logging backend
	Flush()

	// log an log record
	log(r *Record)

	// writes an byte slice to current backend
	write([]byte) error

	// start current backend
	start()

	// returns true if current log backend is a tty for colorful logging
	isatty() bool

	// returns file handler of current backend if it's a file
	fd() Handler
}

func closeBackend(backend Backend) {
	switch backend.Writer().(type) {
	case *os.File:
		file := backend.Writer().(*os.File)
		switch file.Fd() {
		case os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd():
			// do not close this
		default:
			if err := file.Close(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error close file : %v\n", err)
			}
		}
	default:
	}
}
