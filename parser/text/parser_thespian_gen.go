// code generated by thespian; DO NOT EDIT

package text

import (
	import1 "github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
)

// parserBase is embedded in the private actor struct and contains
// common fields as well as default method implementations
type parserBase struct {
	rt *thespian.Runtime
	tx *ParserTx
	rx *ParserRx
}

// handleStart is called when the actor starts.  The default implementation
// does nothing, but users may implement this method to perform startup.
func (a *parserBase) handleStart() {}

// handleStop is called when the actor stops cleanly.  The default
// implementation does nothing, but users may implement this method to perform
// cleanup.
func (a *parserBase) handleStop() {}

// handleSuperEvent is called for supervisory events.  Actors which do not
// supervise need not implement this method.
func (a *parserBase) handleSuperEvent(ev thespian.SuperEvent) {}

// ParserBuilder is used to build new Parser actors.
type ParserBuilder struct {
	parser
	rawLogMessage import1.RawMessageMailbox
}

func (bldr ParserBuilder) spawn(rt *thespian.Runtime) *ParserTx {
	reg := rt.Register()
	bldr.rawLogMessage.ApplyDefaults()

	rx := &ParserRx{
		id:            reg.ID,
		rt:            rt,
		stopChan:      reg.StopChan,
		superChan:     reg.SuperChan,
		healthChan:    reg.HealthChan,
		rawLogMessage: bldr.rawLogMessage.Rx(),
	}

	tx := &ParserTx{
		ID:            reg.ID,
		stopChan:      reg.StopChan,
		rawLogMessage: bldr.rawLogMessage.Tx(),
	}

	// copy to a new parser instance
	pvt := bldr.parser
	pvt.rt = rt
	pvt.rx = rx
	pvt.tx = tx

	go pvt.loop()
	return tx
}

// ParserRx contains the Rx sides of the mailboxes, for access from the
// Parser implementation.
type ParserRx struct {
	id uint64
	rt *thespian.Runtime

	stopChan      <-chan struct{}
	superChan     <-chan thespian.SuperEvent
	healthChan    <-chan struct{}
	rawLogMessage import1.RawMessageRx
}

// supervise starts supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Supervize.
func (rx *ParserRx) supervise(otherID uint64) {
	rx.rt.Supervise(rx.id, otherID)
}

// unsupervise stops supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Unupervize.
func (rx *ParserRx) unsupervise(otherID uint64) {
	rx.rt.Unsupervise(rx.id, otherID)
}

// ParserTx is the public handle for Parser actors.
type ParserTx struct {
	// ID is the unique ID of this actor
	ID            uint64
	stopChan      chan<- struct{}
	rawLogMessage import1.RawMessageTx
}

// Stop sends a message to stop the actor.  This does not wait until
// the actor has stopped.
func (a *ParserTx) Stop() {
	select {
	case a.stopChan <- struct{}{}:
	default:
	}
}

// RawLogMessage sends to the actor's rawLogMessage mailbox.
func (tx *ParserTx) RawLogMessage(m string) {
	tx.rawLogMessage.C <- m
}

func (a *parser) loop() {
	rx := a.rx
	defer func() {

		a.rt.ActorStopped(a.rx.id)
	}()
	a.handleStart()
	for {
		select {
		case <-rx.healthChan:
			// nothing to do
		case ev := <-rx.superChan:
			a.handleSuperEvent(ev)
		case <-rx.stopChan:
			a.handleStop()
			return
		case m := <-rx.rawLogMessage.Chan():
			a.handleRawLogMessage(m)
		}
	}
}