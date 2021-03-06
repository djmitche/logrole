// code generated by thespian; DO NOT EDIT

package file

import (
	"github.com/djmitche/thespian"
	import1 "github.com/djmitche/thespian/mailbox"
)

// tailerBase is embedded in the private actor struct and contains
// common fields as well as default method implementations
type tailerBase struct {
	rt *thespian.Runtime
	tx *TailerTx
	rx *TailerRx
}

// handleStart is called when the actor starts.  The default implementation
// does nothing, but users may implement this method to perform startup.
func (a *tailerBase) handleStart() {}

// handleStop is called when the actor stops cleanly.  The default
// implementation does nothing, but users may implement this method to perform
// cleanup.
func (a *tailerBase) handleStop() {}

// handleSuperEvent is called for supervisory events.  Actors which do not
// supervise need not implement this method.
func (a *tailerBase) handleSuperEvent(ev thespian.SuperEvent) {}

// TailerBuilder is used to build new Tailer actors.
type TailerBuilder struct {
	tailer
}

func (bldr TailerBuilder) spawn(rt *thespian.Runtime) *TailerTx {
	reg := rt.Register()

	rx := &TailerRx{
		id:         reg.ID,
		rt:         rt,
		stopChan:   reg.StopChan,
		superChan:  reg.SuperChan,
		healthChan: reg.HealthChan,
		poll:       import1.NewTickerRx(rt),
	}

	tx := &TailerTx{
		ID:       reg.ID,
		stopChan: reg.StopChan,
	}

	// copy to a new tailer instance
	pvt := bldr.tailer
	pvt.rt = rt
	pvt.rx = rx
	pvt.tx = tx

	go pvt.loop()
	return tx
}

// TailerRx contains the Rx sides of the mailboxes, for access from the
// Tailer implementation.
type TailerRx struct {
	id uint64
	rt *thespian.Runtime

	stopChan   <-chan struct{}
	superChan  <-chan thespian.SuperEvent
	healthChan <-chan struct{}
	poll       import1.TickerRx
}

// supervise starts supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Supervize.
func (rx *TailerRx) supervise(otherID uint64) {
	rx.rt.Supervise(rx.id, otherID)
}

// unsupervise stops supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Unupervize.
func (rx *TailerRx) unsupervise(otherID uint64) {
	rx.rt.Unsupervise(rx.id, otherID)
}

// TailerTx is the public handle for Tailer actors.
type TailerTx struct {
	// ID is the unique ID of this actor
	ID       uint64
	stopChan chan<- struct{}
}

// Stop sends a message to stop the actor.  This does not wait until
// the actor has stopped.
func (a *TailerTx) Stop() {
	select {
	case a.stopChan <- struct{}{}:
	default:
	}
}

func (a *tailer) loop() {
	rx := a.rx
	defer func() {
		rx.poll.Stop()
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
		case t := <-rx.poll.Chan():
			a.handlePoll(t)
		}
	}
}
