package runner

import (
	"log"
	"net"

	"github.com/coreos/go-systemd/activation"
	sysdUtil "github.com/coreos/go-systemd/util"
	"github.com/pkg/errors"
)

// getLocalTCPListener gets a `net.Listener` in a way that is independent of
// platform or whether the program is running as a service. `port` should be
// without ":"
func getLocalTCPListener(port string) (net.Listener, error) {
	isSystemd, err := sysdUtil.RunningFromSystemService()
	if err != nil {
		return nil, err
	}
	if isSystemd {
		return getLocalTCPListenerSystemd(port)
	}
	return net.Listen("tcp", ":"+port)
}

var listeners []net.Listener

// getListener gets a `net.Listener` in a way that is independent of platform
// or whether the program is running as a service. `port` should be without ":"
func getLocalTCPListenerSystemd(port string) (net.Listener, error) {
	log.Printf("Attempting to get listener on port \"%s\"\n", port)
	p, err := net.LookupPort("tcp", port)
	if err != nil {
		return nil, err
	}

	////// Cache listeners for later //////
	if listeners == nil {
		// This call only works once
		lns, err := activation.Listeners()
		switch {
		case err != nil:
			return nil, err
		case len(lns) < 2: // 2 listeners, 1 for each port
			return nil, errors.New(ErrMissingServiceSockets)
		}
		listeners = lns
		log.Printf("Cached %d listeners...", len(listeners))
	}

	////// Determine which, if any, of the listeners match what we want //////
	var result net.Listener
	for _, ln := range listeners {
		// Skip any that aren't tcp, otherwise convert them
		addr, ok := ln.Addr().(*net.TCPAddr)
		if !ok || addr == nil || addr.Network() != "tcp" {
			continue
		}

		// TODO: Determine if we need to filter by localhost or not

		// Select the one with a matching port
		if addr.Port == p {
			result = ln
		}
	}

	if result == nil {
		return nil, errors.New(ErrNoSuchListener)
	}
	return result, nil
}
