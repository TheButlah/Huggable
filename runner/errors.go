package runner

const (
	// ErrMissingServiceSockets occurs when the runner was started by a service
	// and got fewer (possibly none) sockets than expected.
	ErrMissingServiceSockets = "not enough sockets provided by service"
)
