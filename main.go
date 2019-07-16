package main

import (
	"fmt"
	"log"

	"github.com/thebutlah/huggable.us/runner"
)

func main() {
	fmt.Println("Starting server...")
	// passing https gets mapped to the default port for https
	err := runner.Start()
	if err != nil {
		log.Fatal(err)
	}
}
