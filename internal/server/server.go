package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"torserve/internal/rawhttp"
	"torserve/internal/scrub"

	"github.com/fsnotify/fsnotify"
)

// headerFilter wraps an http.ResponseWriter to filter out specific headers
// added automatically by Go before sending the response.
type headerFilter struct {
	http.ResponseWriter
	headersWritten bool
}

// WriteHeader overrides the default header writing behavior to remove
// Date, Last-Modified, ETag, and Accept-Ranges headers for privacy/security.
func (h *headerFilter) WriteHeader(code int) {
	if !h.headersWritten {
		h.Header().Del("Date")
		h.Header().Del("Last-Modified")
		h.Header().Del("ETag")
		h.Header().Del("Accept-Ranges")
		h.headersWritten = true
	}
	h.ResponseWriter.WriteHeader(code)
}

// WatchLive monitors a directory for new file creations.
// When a new file appears, it is scrubbed automatically using scrub.ScrubFile.
func WatchLive(dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Track recently seen files to avoid duplicate processing
	seen := make(map[string]time.Time)

	go func() {
		defer watcher.Close()

		for {
			select {
			case ev := <-watcher.Events:
				// On new file creation
				if ev.Op&fsnotify.Create != 0 {
					// Debounce repeated events on the same file
					if t, ok := seen[ev.Name]; ok && time.Since(t) < time.Second {
						continue
					}
					seen[ev.Name] = time.Now()

					// Delay a bit to ensure file is fully written before scrubbing
					go func(name string) {
						time.Sleep(500 * time.Millisecond)
						ext := strings.ToLower(filepath.Ext(name))
						if err := scrub.ScrubFile(name, ext); err != nil {
							log.Printf("Failed to scrub new file: %v", err)
						}
					}(ev.Name)
				}

			case err := <-watcher.Errors:
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	return watcher.Add(dir)
}

// Start checks if port 8080 is exposed to the clearnet.
// If so, it prints a warning and exits; otherwise, it starts the raw HTTP server.
func Start() error {
	if exposed, err := isPortExternallyAccessible(8080); err != nil {
		fmt.Println("[!] WARNING: Could not verify port exposure:", err)
	} else if exposed {
		fmt.Println("[!!!] SECURITY ERROR: Port 8080 is accessible from the clearnet.")
		fmt.Println("      This can expose your files and defeat Tor anonymity.")
		fmt.Println("      TorServe will NOT run until this is fixed.\n")
		fmt.Println("🔒 To fix this, run:")
		fmt.Println("    sudo iptables -A INPUT -i lo -p tcp --dport 8080 -j ACCEPT")
		fmt.Println("    sudo iptables -A INPUT -p tcp --dport 8080 -j DROP")
		fmt.Println("\n💡 Then restart TorServe.")
		os.Exit(1)
	}

	return rawhttp.Start()
}

// headerSanitizer wraps http.ResponseWriter and removes or replaces
// Go's default headers with more privacy-conscious and secure headers.
type headerSanitizer struct {
	http.ResponseWriter
	headersWritten bool
}

// WriteHeader removes default headers and sets security headers before writing the response.
func (hs *headerSanitizer) WriteHeader(code int) {
	if !hs.headersWritten {
		// Remove unwanted default headers
		hs.Header().Del("Date")
		hs.Header().Del("Last-Modified")
		hs.Header().Del("ETag")
		hs.Header().Del("Accept-Ranges")

		// Set security and privacy headers
		hs.Header().Set("Server", "")
		hs.Header().Set("Cache-Control", "no-store")
		hs.Header().Set("Pragma", "no-cache")
		hs.Header().Set("Expires", "0")
		hs.Header().Set("X-Content-Type-Options", "nosniff")
		hs.Header().Set("X-Frame-Options", "DENY")
		hs.Header().Set("X-XSS-Protection", "1; mode=block")

		hs.headersWritten = true
	}
	hs.ResponseWriter.WriteHeader(code)
}

// isPortExternallyAccessible checks whether a specific TCP port
// is reachable via any non-loopback network interface.
func isPortExternallyAccessible(port int) (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, iface := range interfaces {
		// Skip loopback or down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Try connecting to the port from the interface IP
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip.String(), port), 500*time.Millisecond)
			if err == nil {
				conn.Close()
				return true, nil
			}
		}
	}

	return false, nil
}
