package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"torserve/internal/cloak"
	"torserve/internal/scrub"
	"torserve/internal/server"
	"torserve/internal/tor"
)

func main() {
	newKey := flag.Bool("new-key", false, "Generate a new onion address")
	flag.Parse()

	if *newKey {
		fmt.Println("[*] Removing hidden_service/ to generate a new onion key...")
		if err := os.RemoveAll("hidden_service"); err != nil {
			log.Fatalf("Failed to remove hidden_service directory: %v", err)
		}
	}

	cloak.Init()

	fmt.Println("[*] Scrubbing public/ for unsafe or revealing files...")
	if err := scrub.Init(); err != nil {
		log.Fatalf("Scrub error: %v", err)
	}

	if err := server.WatchLive("public"); err != nil {
		log.Fatalf("Live watch error: %v", err)
	}

	fmt.Println("[*] Starting Tor process...")
	if err := tor.Start(); err != nil {
		log.Fatalf("Failed to start Tor: %v", err)
	}

	fmt.Println("[*] Waiting for hidden service hostname...")
	onion, err := tor.WaitForHostname()
	if err != nil {
		log.Fatalf("Failed to get .onion address: %v", err)
	}
	fmt.Printf("[+] Tor hidden service is live at: http://%s\n", onion)

	fmt.Println("[*] Starting local HTTP server on 127.0.0.1:8080...")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
