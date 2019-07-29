// We will try to bind to ports directly instead of using sockets provided by
// systemd on all platforms other than linux
// +build !linux

package runner

import "net"

// getListener gets a `net.Listener` in a way that is independent of platform
// or whether the program is running as a service. Mimics API of `net.Listen()`
func getListener(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}
