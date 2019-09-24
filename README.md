# golog

A full functional log package for golang with following features

* enhanced but compatible with official `log` package. Simply using import alias to replace official `log`

``` go
package main

import (
    "bytes"
    "fmt"
    "github.com/mysqto/log"
)

func main() {
    var (
        buf    bytes.Buffer
        logger = log.New(&buf, "INFO: ", log.Lshortfile)

        infof = func(info string) {
            logger.Output(2, info)
        }
    )

    infof("Hello world")

    fmt.Print(&buf)
}
```
