package common

type Message struct {
	// Text of the log message
	Text string

	// Source of the log message
	Source string

	// Tags associated with the log message
	Tags []string
}
