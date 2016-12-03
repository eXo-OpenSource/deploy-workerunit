package main

func main() {
	// Create MTA server instance
	server := NewMTAServer("/var/lib/mtasa/mta-server64")

	// Create api
	api := NewApi(server)

	// Listen for requests on the main goroutine
	api.Listen()
}
