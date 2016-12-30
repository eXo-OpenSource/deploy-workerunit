package main

import (
	"fmt"
	"os"
)

func main() {
	// Integrity checks
	if os.Getenv("API_SECRET") == "" {
		fmt.Fprintf(os.Stderr, "API_SECRET environment variable not defined\n")
		return
	}

	// Get watchdog enabled state (defaults to 'disabled')
	watchdogEnabled := false
	if os.Getenv("WATCHDOG") != "" {
		watchdogEnabled = true
	}

	// Open config for patching
	fmt.Println("Patching config...")
	patcher, err := NewMTAConfigPatcher("/var/lib/mtasa/mods/deathmatch/mtaserver.conf")
	if err != nil {
		panic(err)
	}

	// Patch some entries
	patcher.Patch("serverport", os.Getenv("MTA_GAME_PORT"))
	patcher.Patch("httpport", os.Getenv("MTA_HTTP_PORT"))
	patcher.Patch("servername", os.Getenv("MTA_SERVER_NAME"))
	patcher.Patch("serverip", os.Getenv("MTA_SERVER_IP"))
	patcher.Save()

	// Create MTA server instance
	fmt.Println("Creating MTAServer...")
	server := NewMTAServer("/var/lib/mtasa/mta-server64", watchdogEnabled)

	// Should we start the server immediately
	if os.Getenv("START_MTA") != "" {
		server.Start()
	}

	// Create api
	fmt.Println("Creating API...")
	api := NewApi(server)

	// Listen for requests on the main goroutine
	fmt.Println("Waiting for commands...")
	api.Listen()
}
