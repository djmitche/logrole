// Package random implements a "tailer" that provides random log lines.
package random

//go:generate go run github.com/djmitche/thespian/cmd/thespian generate

import (
	"math/rand"
	"time"

	"github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
)

type tailer struct {
	tailerBase

	sink common.StreamDataSink
}

func NewTailer(rt *thespian.Runtime, sink common.StreamDataSink) *TailerTx {
	return TailerBuilder{
		tailer: tailer{
			sink: sink,
		},
	}.spawn(rt)
}

func (t *tailer) handleStart() {
	t.rx.tick.Reset(100 * time.Millisecond)
}

var messages = []string{
	"A thing!  It happened!\n",
	"Warp Speed Engaged\n",
	"Service started\n",
	"Service stopped\n",
}

func (t *tailer) handleTick(_ time.Time) {
	if rand.Intn(100) > 10 {
		t.sink.LogData([]byte(messages[rand.Intn(len(messages))]))
	}
}
