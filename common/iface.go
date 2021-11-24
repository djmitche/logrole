package common

// StreamDataSink a consumer for streams of log data.
type StreamDataSink interface {
	// Handle a slice of log data.  The caller may not modify the buffer
	// after calling this method.
	LogData([]byte)
}

// RawMessageSink is a consumer for sequences of "raw" log data, separated into
// messages.
type RawMessageSink interface {
	// Handle a raw log message.
	RawLogMessage(string)
}

// MessageSink is a consumer for sequences of parsed log data.
type MessageSink interface {
	// Handle a log message.  The caller may not modify the Message
	// after calling this method.
	LogMessage(Message)
}

// BatchSink is a consumer for batches of log messages.
type BatchSink interface {
	// Handle a batch of log messages.  The caller may not modify the
	// batch after calling this method.
	Batch([]Message)
}
