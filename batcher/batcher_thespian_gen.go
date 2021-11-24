// code generated by thespian; DO NOT EDIT

package batcher

import (
	import2 "github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
	import1 "github.com/djmitche/thespian/mailbox"
)

// batcherBase is embedded in the private actor struct and contains
// common fields as well as default method implementations
type batcherBase struct {
	rt *thespian.Runtime
	tx *BatcherTx
	rx *BatcherRx
}

// handleStart is called when the actor starts.  The default implementation
// does nothing, but users may implement this method to perform startup.
func (a *batcherBase) handleStart() {}

// handleStop is called when the actor stops cleanly.  The default
// implementation does nothing, but users may implement this method to perform
// cleanup.
func (a *batcherBase) handleStop() {}

// handleSuperEvent is called for supervisory events.  Actors which do not
// supervise need not implement this method.
func (a *batcherBase) handleSuperEvent(ev thespian.SuperEvent) {}

// BatcherBuilder is used to build new Batcher actors.
type BatcherBuilder struct {
	batcher

	logMessage import2.MessageMailbox
}

func (bldr BatcherBuilder) spawn(rt *thespian.Runtime) *BatcherTx {
	reg := rt.Register()

	bldr.logMessage.ApplyDefaults()

	rx := &BatcherRx{
		id:         reg.ID,
		rt:         rt,
		stopChan:   reg.StopChan,
		superChan:  reg.SuperChan,
		healthChan: reg.HealthChan,
		flush:      import1.NewTickerRx(rt),
		logMessage: bldr.logMessage.Rx(),
	}

	tx := &BatcherTx{
		ID:       reg.ID,
		stopChan: reg.StopChan,

		logMessage: bldr.logMessage.Tx(),
	}

	// copy to a new batcher instance
	pvt := bldr.batcher
	pvt.rt = rt
	pvt.rx = rx
	pvt.tx = tx

	go pvt.loop()
	return tx
}

// BatcherRx contains the Rx sides of the mailboxes, for access from the
// Batcher implementation.
type BatcherRx struct {
	id uint64
	rt *thespian.Runtime

	stopChan   <-chan struct{}
	superChan  <-chan thespian.SuperEvent
	healthChan <-chan struct{}
	flush      import1.TickerRx
	logMessage import2.MessageRx
}

// supervise starts supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Supervize.
func (rx *BatcherRx) supervise(otherID uint64) {
	rx.rt.Supervise(rx.id, otherID)
}

// unsupervise stops supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Unupervize.
func (rx *BatcherRx) unsupervise(otherID uint64) {
	rx.rt.Unsupervise(rx.id, otherID)
}

// BatcherTx is the public handle for Batcher actors.
type BatcherTx struct {
	// ID is the unique ID of this actor
	ID       uint64
	stopChan chan<- struct{}

	logMessage import2.MessageTx
}

// Stop sends a message to stop the actor.  This does not wait until
// the actor has stopped.
func (a *BatcherTx) Stop() {
	select {
	case a.stopChan <- struct{}{}:
	default:
	}
}

// LogMessage sends to the actor's logMessage mailbox.
func (tx *BatcherTx) LogMessage(m import2.Message) {
	tx.logMessage.C <- m
}

func (a *batcher) loop() {
	rx := a.rx
	defer func() {
		rx.flush.Stop()

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
		case t := <-rx.flush.Chan():
			a.handleFlush(t)
		case m := <-rx.logMessage.Chan():
			a.handleLogMessage(m)
		}
	}
}
