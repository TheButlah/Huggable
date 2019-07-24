package main

import (
	"log"

	"github.com/thebutlah/huggable.us/runner"
)

func main() {
	log.Println("Starting server...")
	// passing https gets mapped to the default port for https
	err := runner.Start()
	if err != nil {
		log.Fatal(err)
	}
}
