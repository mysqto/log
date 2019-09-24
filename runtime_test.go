// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"bytes"
	"fmt"
	"path/filepath"
)

func ExampleGetRuntimeInfo() {
	var (
		buf    bytes.Buffer
		logger = New(&buf, "", 0)
	)

	function, file, line := getRuntimeInfo(1)

	logger.Printf("%s:%d:%v", filepath.Base(file), line, filepath.Base(function))

	fmt.Print(&buf)
	// Output:
	// runtime_test.go:19:log.ExampleGetRuntimeInfo
}
