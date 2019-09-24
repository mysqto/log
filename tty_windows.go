// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package log

import (
	tty "github.com/mysqto/isatty"
)

func isatty(fd uintptr) bool {
	return tty.IsCygwinTerminal(fd)
}
