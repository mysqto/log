// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package log

import (
	"syscall"
)

var (
	kernel32DLL                 = syscall.NewLazyDLL("kernel32.dll")
	setConsoleTextAttributeProc = kernel32DLL.NewProc("SetConsoleTextAttribute")
)

// Character attributes
// Note:
// -- The attributes are combined to produce various colors (e.g., Blue + Green will create Cyan).
//    Clearing all foreground or background colors results in black; setting all creates white.
// See https://msdn.microsoft.com/en-us/library/windows/desktop/ms682088(v=vs.85).aspx#_win32_character_attributes.
const (
	fgBlack     = 0x0000
	fgBlue      = 0x0001
	fgGreen     = 0x0002
	fgCyan      = 0x0003
	fgRed       = 0x0004
	fgMagenta   = 0x0005
	fgYellow    = 0x0006
	fgWhite     = 0x0007
	fgIntensity = 0x0008
	fgMask      = 0x000F
)

var (
	colorsW = []uint16{
		FATAL: fgMagenta,
		ERROR: fgRed,
		WARN:  fgYellow,
		INFO:  fgGreen,
		DEBUG: fgCyan,
		none:  fgWhite,
	}

	boldColorsW = []uint16{
		FATAL: fgMagenta | fgIntensity,
		ERROR: fgRed | fgIntensity,
		WARN:  fgYellow | fgIntensity,
		INFO:  fgGreen | fgIntensity,
		DEBUG: fgCyan | fgIntensity,
		none:  fgWhite | fgIntensity,
	}

	colorResetW uint16 = fgWhite
)

// colorful writes the message with console colors under windows
func (r *Record) colorful(backend Backend, bold bool) {

	var col uint16
	if bold {
		col = boldColorsW[r.level]
	} else {
		col = colorsW[r.level]
	}

	setConsoleTextAttribute(backend.fd(), col)
	_ = backend.write(r.msgBuf())
	setConsoleTextAttribute(backend.fd(), colorResetW)
}

// setConsoleTextAttribute sets the attributes of characters written to the
// console screen buf by the WriteFile or WriteConsole function.
// See http://msdn.microsoft.com/en-us/library/windows/desktop/ms686047(v=vs.85).aspx.
func setConsoleTextAttribute(handler Handler, attribute uint16) bool {
	ok, _, _ := setConsoleTextAttributeProc.Call(handler.Fd(), uintptr(attribute), 0)
	return ok != 0
}
