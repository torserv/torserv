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

	"torserve/internal/cloak"
	"torserve/internal/rawhttp"
	"torserve/internal/scrub"

	"github.com/fsnotify/fsnotify"
)

type headerFilter struct {
	http.ResponseWriter
	headersWritten bool
}

func (h *headerFilter) WriteHeader(code int) {
	if !h.headersWritten {
		// Remove Go-added headers right before writing
		h.Header().Del("Date")
		h.Header().Del("Last-Modified")
		h.Header().Del("ETag")
		h.Header().Del("Accept-Ranges")
		h.headersWritten = true
	}
	h.ResponseWriter.WriteHeader(code)
}

// --- First: serveStaticFile function ---
func serveStaticFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		path = "/index.html"
	}

	filePath := fmt.Sprintf("public%s", path)

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("[ERROR] Failed to read file %s: %v", filePath, err)
		http.NotFound(w, r)
		return
	}

	contentType := http.DetectContentType(data)
	w.Header().Set("Content-Type", contentType)

	// Security headers
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Server", "")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	w.WriteHeader(http.StatusOK)

	if strings.HasSuffix(filePath, ".html") {
		rewritten := cloak.RewriteHTMLLinks(string(data))
		w.Write([]byte(rewritten))
	} else {
		w.Write(data)
	}
}

func WatchLive(dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	seen := make(map[string]time.Time)

	go func() {
		defer watcher.Close()
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Create != 0 {
					if t, ok := seen[ev.Name]; ok && time.Since(t) < time.Second {
						continue
					}
					seen[ev.Name] = time.Now()

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

// --- Then: Start function ---
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

// Wrapper to strip Go-injected headers after writing
type headerSanitizer struct {
	http.ResponseWriter
	headersWritten bool
}

func (hs *headerSanitizer) WriteHeader(code int) {
	if !hs.headersWritten {
		hs.Header().Del("Date")
		hs.Header().Del("Last-Modified")
		hs.Header().Del("ETag")
		hs.Header().Del("Accept-Ranges")

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

func isPortExternallyAccessible(port int) (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, iface := range interfaces {
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

			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip.String(), port), 500*time.Millisecond)
			if err == nil {
				conn.Close()
				return true, nil
			}
		}
	}

	return false, nil
}
