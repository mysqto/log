// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package log

// colorful writes the message with console colors under windows
func (r *Record) colorful(backend Backend, bold bool) {
	_ = backend.write(r.msgBuf())
}
