package runner

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

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

// option is the type alias for configuring options for `Start()`
type option func(*config) error

// HTTPOptions configures the HTTP listener for the server. If `port` is not a
// valid number, it will be converted to one using `net.LookupPort("tcp", port)`
func HTTPOptions(port string) option {
	return func(c *config) error {
		p, err := net.LookupPort("tcp", port)
		if err != nil {
			return errors.New("Invalid port set for `HTTPOptions`: " + port)
		}
		c.httpPort = strconv.Itoa(p)
		return err
	}
}

// HTTPSOptions configures the HTTPS listener for the server. If `port` is not a
// valid number, it will be converted to one using `net.LookupPort("tcp", port)`
func HTTPSOptions(port string) option {
	return func(c *config) error {
		p, err := net.LookupPort("tcp", port)
		if err != nil {
			return errors.New("Invalid port set for `HTTPSOptions`: " + port)
		}
		c.httpsPort = strconv.Itoa(p)
		return err
	}
}

// CertOptions configures the location of the TLS/SSL Certificate for the HTTPS
// Listener. Each argument should be a path to the corresponding certificate
// or private key.
func CertOptions(cert, key string) option {
	return func(c *config) error {
		args := [2]string{cert, key}
		for _, arg := range args {
			if _, err := os.Stat(arg); err != nil {
				// arg may or may not exist, but we error anyway
				log.Panic(err)
				return err // Never reach this, here for explicitness
			}
		}
		c.certPath = cert
		c.keyPath = key
		return nil
	}
}

// Start starts the server using the given options to determine the port.
func Start(options ...option) error {
	////// Configure the options for `Start()` //////
	cfg := new(config)
	{
		// Default values for the arguments
		HTTPOptions("http")(cfg)
		HTTPSOptions("https")(cfg)
		CertOptions(
			// Letsencrypt Cert install locations on unix
			"/etc/letsencrypt/live/huggable.us/fullchain.pem",
			"/etc/letsencrypt/live/huggable.us/privkey.pem",
		)(cfg)

		if len(options) > 2 {
			return errors.New(
				"`Start()` should be called with at most 2 options",
			)
		}
		// Mutate config using provided options
		for _, opt := range options {
			err := opt(cfg)
			if err != nil {
				return err
			}
		}
		log.Println("`runner.Start()` config:\n", cfg)
	}

	////// Start http listener that redirects to https //////
	{
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var targetURL, start = r.URL, r.URL.String()

			targetURL.Scheme = "https"
			targetURL.Host = strings.Split(r.Host, ":")[0]
			target := targetURL.String()

			log.Printf("Redirecting %s to %s", start, target)
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		})
		log.Printf("Listening for HTTP requests on port \"%s\"", cfg.httpPort)
		go http.ListenAndServe(":"+cfg.httpPort, nil)
	}

	////// Start main https listener //////
	{
		mux := http.NewServeMux()
		for p, h := range pathMap {
			mux.Handle(p, h)
		}

		log.Printf("Listening for HTTPS requests on port \"%s\"", cfg.httpsPort)
		log.Fatal(http.ListenAndServeTLS(
			":"+cfg.httpsPort, cfg.certPath, cfg.keyPath, mux,
		))
	}

	return nil
}
