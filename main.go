package main

import (
	"fmt"
)

func main() {
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
