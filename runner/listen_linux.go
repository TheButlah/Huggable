package runner

import (
	"fmt"
	"net"

	"github.com/coreos/go-systemd/activation"
	"github.com/pkg/errors"
)

var listeners map[string][]net.Listener

// getListener gets a `net.Listener` in a way that is independent of platform
// or whether the program is running as a service. Mimics API of `net.Listen()`
func getListener(network, address string) (net.Listener, error) {
	if listeners == nil {
		lns, err := activation.ListenersWithNames();
		switch {
		case err != nil:
			return nil, err
		case len(lns) < 2:  // 2 listeners, 1 for each port
			return nil, errors.New(ErrMissingServiceSockets)
		}
	}
	// TODO: Figure out what the heck the keys are for this
	fmt.Println(listeners)
	return nil, nil
}
