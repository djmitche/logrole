package text

import (
	"github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
)

//go:generate go run github.com/djmitche/thespian/cmd/thespian generate

type parser struct {
	parserBase

	sink   common.MessageSink
	source string
	tags   []string
}

func NewParser(rt *thespian.Runtime, sink common.MessageSink, source string, tags []string) *ParserTx {
	return ParserBuilder{
		parser: parser{
			sink:   sink,
			source: source,
			tags:   tags,
		},
	}.spawn(rt)
}

func (p *parser) handleRawLogMessage(raw string) {
	message := common.Message{
		Text:   raw,
		Source: p.source,
		Tags:   p.tags,
	}
	p.sink.LogMessage(message)
}
