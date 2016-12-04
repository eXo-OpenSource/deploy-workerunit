package main

import (
	"fmt"
	"os"
)

func main() {
	// Integrity checks
	if os.Getenv("API_SECRET") == "" {
		fmt.Fprintf(os.Stderr, "API_SECRET environment variable not defined")
		return
	}

	// Create MTA server instance
	fmt.Println("Creating MTAServer...")
	server := NewMTAServer("/var/lib/mtasa/mta-server64")

	// Create api
	fmt.Println("Creating API...")
	api := NewApi(server)

	// Listen for requests on the main goroutine
	fmt.Println("Waiting for commands...")
	api.Listen()
}
