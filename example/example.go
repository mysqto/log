// Copyright (c) 2019 Chen Lei <my@mysq.to>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"github.com/mysqto/isatty"
	"github.com/mysqto/log"
	"os"
)

func example() {
	log.Warnln("damn it")
}

func get(file *os.File) {
	info, err := file.Stat()
	if err != nil {
		log.Errorln(err)
		return
	}

	log.Debugln(info.Size())

	data, err := json.Marshal(info)

	if err != nil {
		log.Errorln(err)
		return
	}
	log.Debugln("%s", string(data))
}

func main() {
	defer log.Flush()
	// log.NewSyslog(log.INFO, "", log.Lshortfile|log.Lshortfunc|log.Lsequence)
	// file, err := os.OpenFile("example.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	panic(err)
	//}
	//defer file.Close()
	log.NewRotateLogger(log.DEBUG, "", 32*log.KB, log.GZIP, log.Lfull)

	if isatty.IsTerminal(os.Stdout.Fd()) {
		log.Infoln("IsTerminal")
	}

	if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		log.Infoln("IsCygwinTerminal")
	}

	for i := 0; i < 32; i++ {
		log.Printf("hello world %d", i)
	}

	get(os.Stderr)

	go example()

	func() {
		log.SetLogLevel(log.DEBUG)
		log.Errorf("hello world %v", os.Getpid())
	}()

	func() {
		log.Debugln("Debug")
		log.Infoln("Info")
		log.Warnln("Warn")
		log.Errorln("Error")
		log.Fatalln("Fatal")
	}()

}
