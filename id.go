// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
)

// id returns the current goroutine id from runtime stack
func id() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

// procName returns the name of current running process
func procName() string {
	return name(filepath.Base(os.Args[0]))
}

// name returns the file name used by path without extension.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.
func name(path string) string {
	for i := len(path) - 1; i > 0 && path[i] != os.PathSeparator; i-- {
		if path[i] == '.' {
			return path[:i]
		}
	}
	return path
}
