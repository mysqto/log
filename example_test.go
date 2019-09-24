// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"fmt"
)

func ExampleLogger() {
	var (
		buf    bytes.Buffer
		logger = New(&buf, "logger: ", Lshortfile)
	)

	logger.Print("Hello, log file!")

	fmt.Print(&buf)
	// Output:
	// logger: example_test.go:18: Hello, log file!
}

func ExampleLogger_Output() {
	var (
		buf    bytes.Buffer
		logger = New(&buf, "INFO: ", Lshortfile)

		infoF = func(info string) {
			_ = logger.Output(3, info)
		}
	)

	infoF("Hello world")

	fmt.Print(&buf)
	// Output:
	// INFO: example_test.go:31: Hello world
}
