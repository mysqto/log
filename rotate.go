// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bufio"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
)

// ByteSize represent file ByteSize in byte
type ByteSize int64

const (
	_ ByteSize = 1 << (iota * 10)
	// KB = 1 kb bytes
	KB
	// MB = 1 mb bytes
	MB
	// GB = 1 gb bytes
	GB
	// TB = 1 tb bytes
	TB
	// PB = 1 pb bytes
	PB
	// EB = 1 eb bytes
	EB
)

// CompressMethod represent the compress method when archive logs
type CompressMethod int

// compress method
const (
	NoCompress CompressMethod = iota
	GZIP
	Zlib
	LZW
)

// String implements the stringer interface
func (c CompressMethod) String() string {
	switch c {
	case GZIP:
		return ".gz"
	case Zlib:
		return ".zlib"
	case LZW:
		return ".lz"
	default:
		return ""
	}
}

// RotateLogger represents an log backend supporting log rotating and compress
type RotateLogger struct {
	stopped     uint32
	maxFiles    int
	maxSize     ByteSize
	writtenSize ByteSize
	filename    string
	fileIndex   int
	out         io.WriteCloser
	mu          sync.Mutex
	queue       chan *Record
	stop        chan struct{} // Notify closing
	compress    CompressMethod
}

// NewRotateLogger creates a rotate logger with given log level and flags
func NewRotateLogger(level Level, prefix string, maxSize ByteSize, compress CompressMethod, flag int) *Logger {

	name := procName()

	l := &Logger{
		level:   level,
		mu:      sync.Mutex{},
		prefix:  prefix,
		flag:    flag,
		backend: NewRotateBackend(fmt.Sprintf("%s.log", name), 32, maxSize, compress),
		index:   0,
		name:    name,
	}
	go l.backend.start()

	logger.Store(l)

	return l
}

// NewRotateBackend creates a rotate logger backend with given parameters
func NewRotateBackend(filename string, maxFiles int, maxSize ByteSize, compress CompressMethod) Backend {

	backend := &RotateLogger{
		maxFiles:    maxFiles,
		maxSize:     maxSize,
		writtenSize: 0,
		filename:    filename,
		mu:          sync.Mutex{},
		queue:       make(chan *Record),
		stop:        make(chan struct{}),
		stopped:     0,
		compress:    compress,
	}

	backend.rotate()

	return backend
}

// Writer returns the io.Writer of current logger
func (l *RotateLogger) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// SetWriter set the io.Writer of current logger
func (l *RotateLogger) SetWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch w.(type) {
	case io.WriteCloser:
		l.out = w.(io.WriteCloser)
	default:
		panic("w must implements io.WriteCloser")
	}
}

// Flush the current log backend
func (l *RotateLogger) Flush() {
	atomic.StoreUint32(&l.stopped, 1)
	close(l.queue)
	<-l.stop
	closeBackend(l)
}

func (l *RotateLogger) log(r *Record) {
	if atomic.LoadUint32(&l.stopped) == 0 {
		l.queue <- r
	}
}

func (l *RotateLogger) write(data []byte) error {
	written, err := l.out.Write(data)
	if err == nil {
		// update written size
		l.writtenSize += ByteSize(written)
	}
	return err
}

func (l *RotateLogger) start() {
	l.rotate()
	for r := range l.queue {
		l.writeLog(r)
	}
	close(l.stop)
}

// since rotating logging always are files(not console or tty) false
func (l *RotateLogger) isatty() bool {
	return false
}

func (l *RotateLogger) fd() Handler {
	if fd, ok := l.out.(Handler); ok {
		return fd
	}
	return nil
}

func (l *RotateLogger) writeLog(r *Record) {
	buf := r.msgBuf()
	size := ByteSize(len(buf))

	if l.writtenSize+size > l.maxSize {
		l.rotate()
	}
	_ = l.write(buf)
}

func (l *RotateLogger) rotate() {

	if l.out != nil {

		if err := l.out.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error closing current writer : %v\n", err)
			return
		}
	}

	if l.compress > NoCompress {
		l.archive()
	}

	// rotate file
	for i := l.maxFiles - 1; i >= -1; i-- {
		fileName := fmt.Sprintf("%s.%d%s", l.filename, i, l.compress)
		newFileName := fmt.Sprintf("%s.%d%s", l.filename, i+1, l.compress)

		if i == -1 {
			fileName = fmt.Sprintf("%s%s", l.filename, l.compress)
		}

		_, err := os.Stat(fileName)

		if err != nil && os.IsNotExist(err) {
			continue
		}

		err = os.Rename(fileName, newFileName)
		if err != nil && !os.IsNotExist(err) {
			_, _ = fmt.Fprintf(os.Stderr, "error moving current file : %v\n", err)
		}
	}
	l.reset()
}

func (l *RotateLogger) reset()  {
	// create a new file
	f, err := os.OpenFile(l.filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	l.out = f
	l.writtenSize = 0
}

// archive the always written log (l.filename) to l.filename.(gz/lz) when compress is requested
func (l *RotateLogger) archive() {

	in, err := os.OpenFile(l.filename, os.O_RDONLY, 0644)

	// fail to open, such as file not exist or some other file system errors
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error opening file %s : %v\n", l.filename, err)
		return
	}

	defer in.Close()

	reader := bufio.NewReader(in)
	content, err := ioutil.ReadAll(reader)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error reading file %s : %v\n", l.filename, err)
		return
	}

	compressedFileName := fmt.Sprintf("%s%s", l.filename, l.compress)

	out, err := os.OpenFile(compressedFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error opening file %s : %v\n", compressedFileName, err)
		return
	}

	// need compress, try to open
	var w io.WriteCloser
	switch l.compress {
	case GZIP:
		w = gzip.NewWriter(out)
	case Zlib:
		w = zlib.NewWriter(out)
	case LZW:
		w = lzw.NewWriter(out, lzw.MSB, 8)
	default:
		w = out
	}

	defer w.Close()
	_, err = w.Write(content)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error commpressing file %s : %v\n", l.filename, err)
	}
}
