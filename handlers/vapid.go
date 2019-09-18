package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	push "github.com/SherClockHolmes/webpush-go"
)

// VAPIDEndpoint is a `http.Handler` that listens for get requests and returns the VAPID public key
type VAPIDEndpoint struct {
	PublicKey  string
	privateKey string
}

// NewVAPIDEndpoint constructs a new VAPIDEndpoint with a random keypair
func NewVAPIDEndpoint() (e VAPIDEndpoint) {
	e = VAPIDEndpoint{}
	err := e.RegenerateKeys()
	if err != nil {
		log.Printf("Error when generating VAPID keys for first time: %s\n", err)
		return
	}
	return e
}

// RegenerateKeys generates a new set of VAPID keys, updates the endpoint, and saves them in the db
func (e *VAPIDEndpoint) RegenerateKeys() (err error) {
	e.privateKey, e.PublicKey, err = push.GenerateVAPIDKeys()
	if err != nil {
		return err
	}
	// TODO: Store keypair in database to ensure that we can still communicate with subscribed clients
	return nil
}

func (e VAPIDEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoded, err := json.Marshal(e)
	if err != nil {
		log.Printf("Error while marshalling JSON in (VAPIDEndpoint).ServeHTTP(): %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(encoded)
	if err != nil {
		log.Printf("Error while writing response in (VAPIDEndpoint).ServeHTTP(): %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//Compile-time check that `VAPIDEndpoint` implements `http.Handler`
var _ http.Handler = (*VAPIDEndpoint)(nil)
