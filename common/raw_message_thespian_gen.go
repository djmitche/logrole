// code generated by thespian; DO NOT EDIT

package common

// RawMessageMailbox is a mailbox for messages of type string.
type RawMessageMailbox struct {
	// C is the bidirectional channel over which messages will be transferred.  If
	// this is not set in the mailbox, a fresh channel will be created automatically.
	C chan string
	// Disabled, if set to true, causes the mailbox to start life disabled.
	Disabled bool
}

// ApplyDefaults applies default settings to this RawMessage, if
// the struct has its zero value.
func (mbox *RawMessageMailbox) ApplyDefaults() {
	if mbox.C == nil {
		mbox.C = make(chan string, 10) // default channel size
	}
}

// Tx creates a RawMessageTx for this mailbox
func (mbox *RawMessageMailbox) Tx() RawMessageTx {
	return RawMessageTx{
		C: mbox.C,
	}
}

// Rx creates a RawMessageRx for this mailbox
func (mbox *RawMessageMailbox) Rx() RawMessageRx {
	return RawMessageRx{
		C:        mbox.C,
		Disabled: mbox.Disabled,
	}
}

// RawMessageTx sends to a mailbox for messages of type string.
type RawMessageTx struct {
	C chan<- string
}

// RawMessageRx receives from a mailbox for messages of type string.
type RawMessageRx struct {
	C <-chan string
	// Disabled, if set to true, will disable receipt of messages from this mailbox.
	Disabled bool
}

// Chan gets a channel for this mailbox, or nil if there is nothing to select from.
func (rx *RawMessageRx) Chan() <-chan string {
	if rx.Disabled {
		return nil
	}
	return rx.C
}