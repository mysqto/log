// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"sync/atomic"
	"time"
)

type writeMode int

// write mode
const (
	writeModeLog writeMode = iota
	writeModeLogf
	writeModeLogln
)

// Record represents a log record and contains the timestamp when the record
// was created, an increasing id, filename and line and finally the actual
// formatted log line.
type Record struct {
	// index and formatted are the 1st two field to keep memory aligned since on 32-bit machine
	// if the atomic value is not aligned, cmpxchg will cause coredump
	index     uint64 // current log index
	formatted uint32 // flag to identify log buffer is formatted
	time      time.Time
	prefix    *string
	module    *string
	level     Level
	file      *string
	line      int
	function  *string
	fmt       *string
	args      []interface{}
	buf       []byte
	flag      int
	mode      writeMode
	color     string
	newline   bool
}

func itoa(buf *[]byte, i, wid int) {
	utoa(buf, uint64(i), wid)
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func utoa(buf *[]byte, i uint64, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// splitLast return the 0
func splitLast(long string, sep byte) string {
	for i := len(long) - 1; i > 0; i-- {
		if long[i] == sep {
			return long[i+1:]
		}
	}
	return long
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * sequence number (if corresponding flags are provided),
//   * level string
//   * logger name (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided),
//   * function name (if corresponding flags are provided).
func (r *Record) formatHeader(buf *[]byte) {
	if r.prefix != nil {
		*buf = append(*buf, *r.prefix...)
	}
	if r.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if r.flag&LUTC != 0 {
			r.time = r.time.UTC()
		}
		if r.flag&Ldate != 0 {
			year, month, day := r.time.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if r.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := r.time.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if r.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, r.time.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	if r.flag&Lsequence != 0 {
		*buf = append(*buf, '[')
		utoa(buf, atomic.LoadUint64(&r.index), 10)
		*buf = append(*buf, "] "...)
	}

	// none is for Print/Printf/Println
	if r.level > none {
		*buf = append(*buf, '[')
		*buf = append(*buf, levels[r.level]...)
		*buf = append(*buf, "] "...)
	}

	if r.flag&Lloggername > 0 && r.module != nil && len(*r.module) > 0 {
		*buf = append(*buf, *r.module...)

		if r.flag&Lgoroutineid != 0 {
			*buf = append(*buf, '-')
			*buf = append(*buf, id()...)
		}
		*buf = append(*buf, ' ')
	}

	if r.flag&(Lshortfile|Llongfile) != 0 {
		if r.flag&Lshortfile != 0 {
			*r.file = splitLast(*r.file, '/')
		}
		*buf = append(*buf, *r.file...)
		*buf = append(*buf, ':')
		itoa(buf, r.line, -1)
		*buf = append(*buf, ':')
		if r.flag&(Lshortfunc|Llongfunc) == 0 {
			*buf = append(*buf, ' ')
		}
	}
	if r.flag&(Lshortfunc|Llongfunc) != 0 {
		if r.flag&Lshortfunc != 0 {
			*r.function = splitLast(*r.function, '.')
		}
		*buf = append(*buf, *r.function...)
		*buf = append(*buf, ": "...)
	}
}

// print the message
func (r *Record) print() string {
	switch r.mode {
	case writeModeLog:
		return fmt.Sprint(r.args...)
	case writeModeLogf:
		return fmt.Sprintf(*r.fmt, r.args...)
	case writeModeLogln:
		return fmt.Sprintln(r.args...)
	default:
		return fmt.Sprintln(r.args...)
	}
}

// msgBuf format the message into Record.buf
func (r *Record) msgBuf() []byte {
	if atomic.LoadUint32(&r.formatted) == 0 {
		r.buf = r.buf[:0]
		if len(r.color) > 0 {
			r.buf = append(r.buf, r.color...)
		}
		r.formatHeader(&r.buf)
		message := r.print()
		r.buf = append(r.buf, message...)
		if r.newline && (len(message) == 0 || message[len(message)-1] != '\n') {
			r.buf = append(r.buf, '\n')
		}
		if len(r.color) > 0 {
			r.buf = append(r.buf, colorReset...)
		}
		atomic.StoreUint32(&r.formatted, 1)
	}
	return r.buf
}

func (r *Record) message() string {
	return string(r.msgBuf())
}
