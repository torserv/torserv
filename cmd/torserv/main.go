package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"torserve/internal/cloak"
	"torserve/internal/scrub"
	"torserve/internal/server"
	"torserve/internal/tor"
)

func main() {
	noFirejail := flag.Bool("no-firejail", false, "Disable Firejail sandboxing if available")
	newKey := flag.Bool("new-key", false, "Generate a new onion address")
	flag.Parse()

	// Relaunch inside Firejail using a profile, if available and not opted out
	if !*noFirejail && os.Getenv("INSIDE_FIREJAIL") == "" {
		if _, err := exec.LookPath("firejail"); err == nil {
			fmt.Println("[*] Launching inside Firejail sandbox...")
			cwd, _ := os.Getwd()
			profilePath := filepath.Join(cwd, "torserv.profile")
			args := append([]string{
				"--profile=" + profilePath,
				"--private=" + cwd,
				os.Args[0],
			}, os.Args[1:]...)
			cmd := exec.Command("firejail", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Env = append(os.Environ(), "INSIDE_FIREJAIL=1")
			cmd.Run()
			return
		} else {
			fmt.Println("[!] Firejail not found. To enable sandboxing, install it with:")
			fmt.Println("    sudo apt install firejail")
		}
	}

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
