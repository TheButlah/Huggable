package runner

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/thebutlah/huggable.us/handlers"
)

const (
	staticDir = "web/static"
)

var pathMap = map[string]http.Handler{
	"/":               handlers.NewStaticContent(staticDir),
	"/api/push/vapid": handlers.NewVAPIDEndpoint(),
}

//// Start and option config ////

// TLSMode enum
type tlsMode int

const (
	// InvalidTLSMode ensures that not configuring the mode fails
	InvalidTLSMode tlsMode = iota
	// SelfSignedTLSMode uses self-signed certs generated on the fly
	SelfSignedTLSMode
	// AutomaticTLSMode uses the autocert package to get TLS certs
	AutomaticTLSMode
	// ProvidedTLSMode uses certs provided by the user in files
	ProvidedTLSMode
)

func (m tlsMode) String() string {
	return [...]string{"Invalid", "SelfSigned", "Automatic", "Provided"}[m]
}

// config is used to aggregate configured options for `Start()`
type config struct {
	httpPort, httpsPort string
	certsPath           string
	hosts               []string
	mode                tlsMode
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

// AutomaticTLS configures the HTTPS listener to automatically get TLS Certs.
// `hosts` is a list of hosts that belong to this server.
func AutomaticTLS(hosts []string) Option {
	return func(c *config) error {
		c.mode = AutomaticTLSMode
		c.certsPath = "certs"
		// TODO: Do we need to include localhost here or check for length >= 1?
		c.hosts = hosts
		return nil
	}
}

// SelfSignedTLS configures the HTTPS listener to use self-signed TLS Certs
// generated on the fly.
// `hosts` is a list of hosts that belong to this server.
func SelfSignedTLS(hosts []string) Option {
	return func(c *config) error {
		c.mode = SelfSignedTLSMode
		// TODO: Do we need to include localhost here or check for length >= 1?
		c.hosts = hosts
		return nil
	}
}

//TODO: Implement ProvidedTLS

// Start starts the server using the given options to determine the port.
// `hosts` should be a list of hosts registered with the certificate authority
// that point to our IP address. Handles localhost as a special case.
func Start(hosts []string, options ...Option) error {
	////// Configure the options for `Start()` //////
	cfg := new(config)
	{
		// Default values for the arguments
		errs := [...]error{
			HTTPPort("http")(cfg),
			HTTPSPort("https")(cfg),
			AutomaticTLS(hosts)(cfg),
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

		var tlsConfig *tls.Config
		switch cfg.mode {
		case AutomaticTLSMode:
			tlsConfig, err = automaticTLSConfig(cfg.certsPath, cfg.hosts)
		case SelfSignedTLSMode:
			tlsConfig, err = selfSignedTLSConfig(cfg.hosts)
		default:
			err = fmt.Errorf("TLS mode was invalid")
		}
		if err != nil {
			return err
		}

		server := http.Server{
			Addr:      ln.Addr().String(),
			Handler:   mux,
			TLSConfig: tlsConfig,

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
