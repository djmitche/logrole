package main

import (
	"github.com/djmitche/logrole/tailer/file"
	"github.com/djmitche/thespian"
)

func main() {
	rt := thespian.NewRuntime()
	_ = file.NewTailer(rt, "/tmp/test.log")

	rt.Run()
}
