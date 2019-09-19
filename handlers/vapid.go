package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	push "github.com/SherClockHolmes/webpush-go"
)

// VAPIDEndpoint is a `http.Handler` REST endpoint that listens for get requests and returns the VAPID public key
type VAPIDEndpoint struct {
	PublicKey  string
	privateKey string
}

//Compile-time check that `VAPIDEndpoint` implements `http.Handler`
var _ http.Handler = (*VAPIDEndpoint)(nil)

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
	// This endpoint only accepts GET requests
	if r.Method != "" && r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Only uses exported fields of the struct
	encoded, err := json.Marshal(e)
	if err != nil {
		log.Printf("Error while marshalling JSON in (VAPIDEndpoint).ServeHTTP(): %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setup headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "No-Cache")

	// Write response along with headers
	_, err = w.Write(encoded)
	if err != nil {
		log.Printf("Error while writing response in (VAPIDEndpoint).ServeHTTP(): %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
