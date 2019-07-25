package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// RedirectHTTP is a `http.Handler` that redirects all HTTP traffic to a new
// scheme, typically HTTPS.
type RedirectHTTP struct {
	targetScheme string
}

//Compile-time check that `RedirectHTTP` implements `http.Handler`
var _ http.Handler = (*RedirectHTTP)(nil)

// NewRedirectHTTP constructs a `RedirectHTTP` handler.
// `targetScheme` should not have any colons or slashes (i.e. https vs https://)
func NewRedirectHTTP(targetScheme string) RedirectHTTP {
	invalidRunes := ":/"
	if strings.ContainsAny(targetScheme, invalidRunes) {
		log.Panic(fmt.Errorf(
			"handlers: `targetScheme` must not contain `%s`", invalidRunes,
		))
	}
	return RedirectHTTP{targetScheme: targetScheme}
}

func (rh RedirectHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var targetURL, start = r.URL, r.URL.String()

	targetURL.Scheme = rh.targetScheme
	targetURL.Host = strings.Split(r.Host, ":")[0]
	target := targetURL.String()

	log.Printf("Redirecting %s to %s", start, target)
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
