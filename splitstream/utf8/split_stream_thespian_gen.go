// code generated by thespian; DO NOT EDIT

package utf8

import (
	import1 "github.com/djmitche/logrole/common"
	"github.com/djmitche/thespian"
)

// splitStreamBase is embedded in the private actor struct and contains
// common fields as well as default method implementations
type splitStreamBase struct {
	rt *thespian.Runtime
	tx *SplitStreamTx
	rx *SplitStreamRx
}

// handleStart is called when the actor starts.  The default implementation
// does nothing, but users may implement this method to perform startup.
func (a *splitStreamBase) handleStart() {}

// handleStop is called when the actor stops cleanly.  The default
// implementation does nothing, but users may implement this method to perform
// cleanup.
func (a *splitStreamBase) handleStop() {}

// handleSuperEvent is called for supervisory events.  Actors which do not
// supervise need not implement this method.
func (a *splitStreamBase) handleSuperEvent(ev thespian.SuperEvent) {}

// SplitStreamBuilder is used to build new SplitStream actors.
type SplitStreamBuilder struct {
	splitStream
	logData import1.DataMailbox
}

func (bldr SplitStreamBuilder) spawn(rt *thespian.Runtime) *SplitStreamTx {
	reg := rt.Register()
	bldr.logData.ApplyDefaults()

	rx := &SplitStreamRx{
		id:         reg.ID,
		rt:         rt,
		stopChan:   reg.StopChan,
		superChan:  reg.SuperChan,
		healthChan: reg.HealthChan,
		logData:    bldr.logData.Rx(),
	}

	tx := &SplitStreamTx{
		ID:       reg.ID,
		stopChan: reg.StopChan,
		logData:  bldr.logData.Tx(),
	}

	// copy to a new splitStream instance
	pvt := bldr.splitStream
	pvt.rt = rt
	pvt.rx = rx
	pvt.tx = tx

	go pvt.loop()
	return tx
}

// SplitStreamRx contains the Rx sides of the mailboxes, for access from the
// SplitStream implementation.
type SplitStreamRx struct {
	id uint64
	rt *thespian.Runtime

	stopChan   <-chan struct{}
	superChan  <-chan thespian.SuperEvent
	healthChan <-chan struct{}
	logData    import1.DataRx
}

// supervise starts supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Supervize.
func (rx *SplitStreamRx) supervise(otherID uint64) {
	rx.rt.Supervise(rx.id, otherID)
}

// unsupervise stops supervision of the actor identified by otherID.
// It is a shortcut to thespian.Runtime.Unupervize.
func (rx *SplitStreamRx) unsupervise(otherID uint64) {
	rx.rt.Unsupervise(rx.id, otherID)
}

// SplitStreamTx is the public handle for SplitStream actors.
type SplitStreamTx struct {
	// ID is the unique ID of this actor
	ID       uint64
	stopChan chan<- struct{}
	logData  import1.DataTx
}

// Stop sends a message to stop the actor.  This does not wait until
// the actor has stopped.
func (a *SplitStreamTx) Stop() {
	select {
	case a.stopChan <- struct{}{}:
	default:
	}
}

// LogData sends to the actor's logData mailbox.
func (tx *SplitStreamTx) LogData(m []byte) {
	tx.logData.C <- m
}

func (a *splitStream) loop() {
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
		case m := <-rx.logData.Chan():
			a.handleLogData(m)
		}
	}
}
