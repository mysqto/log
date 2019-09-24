// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
)

type color int

// ANSI color for terminal
const (
	colorBlack color = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

func colorSeq(color color) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

func colorSeqBold(color color) string {
	return fmt.Sprintf("\033[%d;1m", int(color))
}

var (
	colors = []string{
		FATAL: colorSeq(colorMagenta),
		ERROR: colorSeq(colorRed),
		WARN:  colorSeq(colorYellow),
		INFO:  colorSeq(colorGreen),
		DEBUG: colorSeq(colorCyan),
		none:  colorSeq(colorWhite),
	}

	boldColors = []string{
		FATAL: colorSeqBold(colorMagenta),
		ERROR: colorSeqBold(colorRed),
		WARN:  colorSeqBold(colorYellow),
		INFO:  colorSeqBold(colorGreen),
		DEBUG: colorSeqBold(colorCyan),
		none:  colorSeqBold(colorWhite),
	}

	colorReset = "\033[0m"
)

// colorfultty writes the message with ANSI colors under windows
func (r *Record) colorfultty(backend Backend, bold bool) {
	var col string
	if bold {
		col = boldColors[r.level]
	} else {
		col = colors[r.level]
	}
	r.color = col
	_ = backend.write(r.msgBuf())
}
