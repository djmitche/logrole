package main

import (
	"log"

	"github.com/djmitche/logrole/batcher"
	"github.com/djmitche/logrole/common"
	"github.com/djmitche/logrole/parser/text"
	"github.com/djmitche/logrole/splitstream/utf8"
	"github.com/djmitche/logrole/tailer/file"
	"github.com/djmitche/thespian"
)

type sink struct{}

func (s *sink) Batch(batch []common.Message) {
	log.Printf("got messages:")
	for _, msg := range batch {
		log.Printf("  %s from %s", msg.Text, msg.Source)
	}
}

func main() {
	rt := thespian.NewRuntime()

	b := batcher.New(rt, &sink{})
	{
		parser := text.NewParser(rt, b, "test", []string{})
		ss := utf8.NewSplitStream(rt, parser)
		_ = file.NewTailer(rt, "/tmp/test.log", ss)
	}

	{
		parser := text.NewParser(rt, b, "test2", []string{})
		ss := utf8.NewSplitStream(rt, parser)
		_ = file.NewTailer(rt, "/tmp/test2.log", ss)
	}

	rt.Run()
}
