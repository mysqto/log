// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"runtime"
)

// getRuntimeInfo returns function name, file name, file line of current call stack
// if cannot get those info from runtime, will return a default value ??? for function
// and file along with 0 for line
func getRuntimeInfo(depth int) (string, string, int) {
	pc, fn, ln, ok := runtime.Caller(depth)
	if !ok {
		fn = "???"
		ln = 0
	}
	function := "???"
	caller := runtime.FuncForPC(pc)
	if caller != nil {
		function = caller.Name()
	}
	return function, fn, ln
}
