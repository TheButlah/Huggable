package runner

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/thebutlah/huggable.us/handlers"
	"golang.org/x/crypto/acme/autocert"
)

var pathMap = map[string]http.Handler{
	"/": handlers.NewStaticContent("web/static"),
}

//// Start and option config ////

// config is used to aggregate configured options for `Start()`
type config struct {
	httpPort, httpsPort string
	certPath            string
	useSelfSigned       bool // TODO: Actually allow this to be True
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

// CertCache configures the location of the TLS/SSL Certificate Cache for the
// HTTPS Listener. Path should be specified unix-style.
func CertCache(path string) Option {
	return func(c *config) error {
		path = filepath.FromSlash(path)
		// No need for error check - DirCache will create directory
		c.certPath = path
		return nil
	}
}

// UseSelfSignedCerts Forces the app to use self-signed certificates.
var UseSelfSignedCerts = func(c *config) error {
	c.useSelfSigned = true
	return nil
}

// Start starts the server using the given options to determine the port.
// `domains` should be a list of domains registered with the certificate
// authority that point to our IP address. Handles localhost as a special case.
func Start(domains []string, options ...Option) error {
	////// Configure the options for `Start()` //////
	cfg := new(config)
	{
		// Default values for the arguments
		errs := [...]error{
			HTTPPort("http")(cfg),
			HTTPSPort("https")(cfg),
			CertCache("certs/")(cfg),
		}

		for _, err := range errs {
			if err != nil {
				return err
			}
		}

		// Mutate config using provided options
		for _, opt := range options {
			if err := opt(cfg); err != nil {
				return err
			}
		}

		// TODO: Make this supported
		if cfg.useSelfSigned {
			return fmt.Errorf("error: using self-signed certs is not yet supported")
		}

		log.Printf("runner config: %+v\n", cfg)
	}

	log.Println("Starting the server...")

	////// Start http listener that redirects to https //////
	{
		mux := http.NewServeMux()
		mux.Handle("/", handlers.NewRedirectHTTP("https"))

		ln, err := getLocalTCPListener(cfg.httpPort)
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

		ln, err := getLocalTCPListener(cfg.httpsPort)
		if err != nil {
			log.Panic(err)
		}
		defer ln.Close()

		// Build a manager for the certificates
		certManager := &autocert.Manager{
			Cache:      autocert.DirCache(cfg.certPath),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
		}
		server := http.Server{
			Addr:      ln.Addr().String(),
			Handler:   mux,
			TLSConfig: certManager.TLSConfig(),

			// Don't hold resources forever
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		log.Printf("Listening for HTTPS requests on \"%s\"", ln.Addr())
		// certFile and keyFile already specified in TLSConfig
		return server.ServeTLS(ln, "", "")
	}
}
