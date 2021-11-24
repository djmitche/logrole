// Package file implements a tailer that will follow a single on-disk file.
package file

//go:generate go run github.com/djmitche/thespian/cmd/thespian generate

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
)

const (
	sleepInterval = time.Second * 1
)

type tailer struct {
	tailerBase

	sink     common.StreamDataSink
	filename string
	file     *os.File
	buf      []byte
}

func NewTailer(rt *thespian.Runtime, filename string, sink common.StreamDataSink) *TailerTx {
	return TailerBuilder{
		tailer{
			filename: filename,
			sink:     sink,
		},
	}.spawn(rt)
}

func (t *tailer) handleStart() {
	t.rx.poll.Reset(sleepInterval)
}

func (t *tailer) handleStop() {
	if t.file != nil {
		t.file.Close()
	}
}

func (t *tailer) handlePoll(_ time.Time) {
	if t.file == nil {
		t.checkForFile()
	}

	if t.file != nil {
		t.tryRead()
	}
}

func (t *tailer) checkForFile() {
	f, err := os.Open(t.filename)
	if err != nil {
		log.Printf("Error opening tailing file: %s (will retry)", err)
		return
	}
	t.file = f
}

func (t *tailer) tryRead() {
	for {
		if t.buf == nil {
			t.buf = make([]byte, 4096)
		}
		n, err := t.file.Read(t.buf)
		if err != nil && err != io.EOF {
			log.Printf("Error reading tailing file: %s (will retry)", err)
		}
		if n == 0 {
			break
		}
		t.sink.LogData(t.buf[:n])
		t.buf = nil
	}
}
