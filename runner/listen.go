// We will try to bind to ports directly instead of using sockets provided by
// systemd on all platforms other than linux
// +build !linux

package runner

import "net"

// getLocalTCPListener gets a `net.Listener` in a way that is independent of 
// platform or whether the program is running as a service. `port` should be 
// without ":"
func getLocalTCPListener(port string) (net.Listener, error) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	return ln.(*net.TCPListener), nil
}
