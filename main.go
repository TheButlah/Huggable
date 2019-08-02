package main

import (
	"log"
	"os"

	"github.com/thebutlah/huggable.us/runner"
)

func main() {
	log.Println("Starting server...")
	// passing https gets mapped to the default port for https
	log.Println(os.Getwd())
	err := runner.Start(
		[]string{"www.huggable.us", "huggable.us"},
	)
	if err != nil {
		log.Fatal(err)
	}
}
