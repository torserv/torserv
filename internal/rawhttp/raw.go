package rawhttp

import (
	"bufio"
	"fmt"
	"math/rand"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"torserve/internal/cloak"
)

// Padding const
const paddingSize = 512 * 1024

// init registers custom MIME types for specific file extensions.
// This ensures correct Content-Type headers are set when serving these files.
func init() {
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".woff2", "font/woff2")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".ttf", "font/ttf")
	mime.AddExtensionType(".otf", "font/otf")
	mime.AddExtensionType(".svg", "image/svg+xml")
}

// Start runs the raw HTTP server on port 8080 and handles incoming connections.
func Start() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("failed to start raw HTTP server: %w", err)
	}
	fmt.Println("[*] HTTP server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// handleConnection processes a single HTTP request over a raw TCP connection.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Add jitter to response timing to avoid fingerprinting
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(50+rng.Intn(150)) * time.Millisecond)

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR] Failed to read request line:", err)
		return
	}

	// Parse request method and path
	method, path, _, ok := parseRequestLine(line)
	if !ok || (method != "GET" && method != "HEAD") {
		fmt.Println("[ERROR] Invalid method or request line")
		writeNotFound(conn)
		return
	}

	// Determine and decrypt requested path
	var decrypted string
	switch path {
	case "/":
		decrypted = "/index.html"
	case "/favicon.ico":
		writeNotFound(conn) // Skip favicon
		return
	default:
		decrypted, err = cloak.DecryptPath(strings.TrimPrefix(path, "/"))
		if err != nil {
			fmt.Println("[ERROR] Failed to decrypt path:", err)
			writeNotFound(conn)
			return
		}
	}

	// Normalize and validate the decrypted path
	cleanPath := filepath.Clean(decrypted)
	if strings.Contains(cleanPath, "..") {
		fmt.Println("[ERROR] Path traversal attempt blocked")
		writeNotFound(conn)
		return
	}
	if cleanPath == "/" {
		cleanPath = "/index.html"
	}

	// Build the full file path to serve
	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Println("[ERROR] Failed to get working directory:", err)
		writeServerError(conn)
		return
	}
	fullPath := filepath.Join(baseDir, "public", cleanPath)

	// Prevent access outside the public directory
	if !strings.HasPrefix(fullPath, filepath.Join(baseDir, "public")) {
		fmt.Println("[ERROR] File outside allowed directory")
		writeNotFound(conn)
		return
	}

	// Read the requested file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Println("[ERROR] Failed to read file:", err)
		writeNotFound(conn)
		return
	}

	// Determine MIME type based on file extension
	ext := filepath.Ext(fullPath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Pad file content to 512KB alignment
	paddedLen := ((len(data) + paddingSize - 1) / paddingSize) * paddingSize

	// Serve the file with optional rewriting for HTML or CSS
	if method == "GET" {
		switch {
		case strings.HasSuffix(fullPath, ".html"):
			rewritten := cloak.RewriteHTMLLinks(string(data))
			rewrittenBytes := []byte(rewritten)
			paddedLen = ((len(rewrittenBytes) + 512*1024 - 1) / (512 * 1024)) * (512 * 1024)
			writeHeaders(conn, paddedLen, mimeType)
			conn.Write(rewrittenBytes)
			conn.Write(make([]byte, paddedLen-len(rewrittenBytes)))

		case strings.HasSuffix(fullPath, ".css"):
			rewritten := cloak.RewriteCSSLinks(string(data))
			rewrittenBytes := []byte(rewritten)
			paddedLen = ((len(rewrittenBytes) + 512*1024 - 1) / (512 * 1024)) * (512 * 1024)
			writeHeaders(conn, paddedLen, mimeType)
			conn.Write(rewrittenBytes)
			conn.Write(make([]byte, paddedLen-len(rewrittenBytes)))

		default:
			writeHeaders(conn, paddedLen, mimeType)
			conn.Write(data)
			conn.Write(make([]byte, paddedLen-len(data)))
		}
	}
}

// parseRequestLine splits an HTTP request line into method, path, and version.
// Returns ok=false if the line is malformed.
func parseRequestLine(line string) (method, path, version string, ok bool) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) < 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// writeHeaders sends HTTP 200 OK headers to the client with the specified content length and MIME type.
// Includes standard security headers and disables caching.
func writeHeaders(conn net.Conn, contentLength int, contentType string) {
	headers := fmt.Sprintf(
		"HTTP/1.0 200 OK\r\n"+
			"Content-Type: %s\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n"+
			"Cache-Control: no-store\r\n"+
			"Pragma: no-cache\r\n"+
			"Expires: 0\r\n"+
			"X-Content-Type-Options: nosniff\r\n"+
			"X-Frame-Options: DENY\r\n"+
			"X-XSS-Protection: 1; mode=block\r\n"+
			"Referrer-Policy: no-referrer\r\n"+
			"\r\n", contentType, contentLength)
	conn.Write([]byte(headers))
}

// writeNotFound sends a fake 200 OK with a chunked transfer response that pretends to serve content
// but actually trickles garbage data to trap, confuse and waste resources of bots, scanners and probes.
// This costs our server almost nothing to do but has a sizeable impact on hostile bots.
func writeNotFound(conn net.Conn) {
	headers := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: application/octet-stream\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Connection: keep-alive\r\n" +
		"X-Trap: You fell for it\r\n\r\n"
	conn.Write([]byte(headers))

	for i := 0; i < 256; i++ {
		time.Sleep(2 * time.Second)
		chunk := fmt.Sprintf("%x\r\n%s\r\n", 5, "trash")
		_, err := conn.Write([]byte(chunk))
		if err != nil {
			break
		}
	}

	// Indicate end of chunked transfer
	conn.Write([]byte("0\r\n\r\n"))
	conn.Close()
}

// writeServerError sends a simple 500 Internal Server Error page with HTML content.
func writeServerError(conn net.Conn) {
	body := "<h1>500 Internal Server Error</h1>"
	headers := fmt.Sprintf(
		"HTTP/1.1 500 Internal Server Error\r\n"+
			"Content-Type: text/html\r\n"+
			"Content-Length: %d\r\n\r\n%s", len(body), body)
	conn.Write([]byte(headers))
}
