package batcher

import (
	"time"

	"github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
)

//go:generate go run github.com/djmitche/thespian/cmd/thespian generate

const (
	flushInterval = 2 * time.Second
)

type batcher struct {
	batcherBase

	messages []common.Message
	sink     common.BatchSink
}

func New(rt *thespian.Runtime, sink common.BatchSink) *BatcherTx {
	return BatcherBuilder{
		batcher: batcher{
			sink: sink,
		},
	}.spawn(rt)
}

func (b *batcher) handleStart() {
	b.rx.flush.Reset(flushInterval)
}

func (b *batcher) handleLogMessage(msg common.Message) {
	b.messages = append(b.messages, msg)
}

func (b *batcher) handleFlush(_ time.Time) {
	if b.messages != nil {
		b.sink.Batch(b.messages)
		b.messages = nil
	}
}
