package runner

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/thebutlah/huggable.us/handlers"
)

var pathMap = map[string]http.Handler{
	"/": handlers.NewStaticContent("web/static"),
}

//// Start and option config ////

// config is used to aggregate configured options for `Start()`
type config struct {
	httpPort, httpsPort string
	certPath, keyPath   string
}

// Option is the type alias for configuring options for `Start()`
type Option func(*config) error

// HTTPPort configures the HTTP listener for the server. Will attempt to map
// `port` to a valid number using `net.LookupPort("tcp", port)`
func HTTPPort(port string) Option {
	return func(c *config) error {
		p, err := net.LookupPort("tcp", port)
		if err != nil {
			return fmt.Errorf(
				"runner: port `%s` invalid for `HTTPPort`", port,
			)
		}
		c.httpPort = strconv.Itoa(p)
		return nil
	}
}

// HTTPSPort configures the HTTP listener for the server. Will attempt to map
// `port` to a valid number using `net.LookupPort("tcp", port)`
func HTTPSPort(port string) Option {
	return func(c *config) error {
		p, err := net.LookupPort("tcp", port)
		if err != nil {
			return fmt.Errorf(
				"runner: port `%s` invalid for `HTTPSPort`", port,
			)
		}
		c.httpsPort = strconv.Itoa(p)
		return nil
	}
}

// CertPaths configures the location of the TLS/SSL Certificate for the HTTPS
// Listener. Each argument should be a path to the corresponding certificate
// or private key.
func CertPaths(cert, key string) Option {
	return func(c *config) error {
		args := [2]string{cert, key}
		for _, arg := range args {
			if _, err := os.Stat(arg); err != nil {
				// arg may or may not exist, but we error anyway
				return fmt.Errorf(
					"runner: could not determine if file `%s` exists", arg,
				)
			}
		}
		c.certPath = cert
		c.keyPath = key
		return nil
	}
}

// Start starts the server using the given options to determine the port.
func Start(options ...Option) error {
	////// Configure the options for `Start()` //////
	cfg := new(config)
	{
		// Default values for the arguments
		HTTPPort("http")(cfg)
		HTTPSPort("https")(cfg)
		CertPaths(
			// Letsencrypt Cert install locations on unix
			"/etc/letsencrypt/live/huggable.us/fullchain.pem",
			"/etc/letsencrypt/live/huggable.us/privkey.pem",
		)(cfg)

		// Mutate config using provided options
		for _, opt := range options {
			if err := opt(cfg); err != nil {
				return err
			}
		}
		log.Println("`runner.Start()` config:\n", cfg)
	}

	////// Start http listener that redirects to https //////
	{
		mux := http.NewServeMux()
		mux.Handle("/", handlers.NewRedirectHTTP("https"))

		ln, err := net.Listen("tcp", ":"+cfg.httpPort)
		if err != nil {
			// TODO: If Listen fails, try to bind to systemd provided socket
			log.Panic(err)
		} else {
			defer ln.Close()
			log.Printf("Listening for HTTP requests on \"%s\"", ln.Addr())
			go func() { log.Println(http.Serve(ln, mux)) }()
		}
	}

	////// Start main https listener //////
	{
		mux := http.NewServeMux()
		for p, h := range pathMap {
			mux.Handle(p, h)
		}

		ln, err := net.Listen("tcp", ":"+cfg.httpsPort)
		if err != nil {
			// TODO: If Listen fails, try to bind to systemd provided socket
			log.Panic(err)
		} else {
			defer ln.Close()
			log.Printf("Listening for HTTPS requests on \"%s\"", ln.Addr())
			return http.ServeTLS(ln, mux, cfg.certPath, cfg.keyPath)
		}
	}
	return fmt.Errorf("runner: was never supposed to reach this")
}
