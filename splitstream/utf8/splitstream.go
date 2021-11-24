// Package utf8 implements a splitstream that assumes utf-8 formatted data
// and splits on newlines.
package utf8

import (
	"bytes"

	"github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
)

//go:generate go run github.com/djmitche/thespian/cmd/thespian generate

type splitStream struct {
	splitStreamBase

	sink        common.RawMessageSink
	unprocessed []byte
}

func NewSplitStream(rt *thespian.Runtime, sink common.RawMessageSink) *SplitStreamTx {
	return SplitStreamBuilder{
		splitStream: splitStream{
			sink: sink,
		},
	}.spawn(rt)
}

func (ss *splitStream) handleLogData(data []byte) {
	unprocessed := append(ss.unprocessed, data...)

	for {
		i := bytes.IndexByte(unprocessed, byte('\n'))
		if i == -1 {
			break
		}
		rawMsg := string(unprocessed[:i])
		unprocessed = unprocessed[i+1:]

		ss.sink.RawLogMessage(rawMsg)
	}

	ss.unprocessed = unprocessed
}
