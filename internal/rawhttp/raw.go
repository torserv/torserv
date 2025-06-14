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

// Start runs the raw HTTP server on port 8080
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

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(50+rng.Intn(150)) * time.Millisecond)

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR] Failed to read request line:", err)
		return
	}

	method, path, _, ok := parseRequestLine(line)

	if !ok || (method != "GET" && method != "HEAD") {
		fmt.Println("[ERROR] Invalid method or request line")
		writeNotFound(conn)
		return
	}

	var decrypted string
	if path == "/" {
		decrypted = "/index.html"
	} else if path == "/favicon.ico" {
		// Ignore the favicon.ico request by returning a 404 without decryption
		writeNotFound(conn)
		return
	} else {
		decrypted, err = cloak.DecryptPath(strings.TrimPrefix(path, "/"))
		if err != nil {
			fmt.Println("[ERROR] Failed to decrypt path:", err)
			writeNotFound(conn)
			return
		}
	}

	cleanPath := filepath.Clean(decrypted)

	if strings.Contains(cleanPath, "..") {
		fmt.Println("[ERROR] Path traversal attempt blocked")
		writeNotFound(conn)
		return
	}
	if cleanPath == "/" {
		cleanPath = "/index.html"
	}

	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Println("[ERROR] Failed to get working directory:", err)
		writeServerError(conn)
		return
	}
	fullPath := filepath.Join(baseDir, "public", cleanPath)

	if !strings.HasPrefix(fullPath, filepath.Join(baseDir, "public")) {
		fmt.Println("[ERROR] File outside allowed directory")
		writeNotFound(conn)
		return
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Println("[ERROR] Failed to read file:", err)
		writeNotFound(conn)
		return
	}

	ext := filepath.Ext(fullPath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	paddedLen := ((len(data) + 512*1024 - 1) / (512 * 1024)) * (512 * 1024)

	if method == "GET" {
		if strings.HasSuffix(fullPath, ".html") {
			rewritten := cloak.RewriteHTMLLinks(string(data))
			rewrittenBytes := []byte(rewritten)
			paddedLen = ((len(rewrittenBytes) + 512*1024 - 1) / (512 * 1024)) * (512 * 1024)
			writeHeaders(conn, paddedLen, mimeType)
			conn.Write(rewrittenBytes)
			conn.Write(make([]byte, paddedLen-len(rewrittenBytes)))
		} else if strings.HasSuffix(fullPath, ".css") {
			rewritten := cloak.RewriteCSSLinks(string(data))
			rewrittenBytes := []byte(rewritten)
			paddedLen = ((len(rewrittenBytes) + 512*1024 - 1) / (512 * 1024)) * (512 * 1024)
			writeHeaders(conn, paddedLen, mimeType)
			conn.Write(rewrittenBytes)
			conn.Write(make([]byte, paddedLen-len(rewrittenBytes)))
		} else {
			writeHeaders(conn, paddedLen, mimeType)
			conn.Write(data)
			conn.Write(make([]byte, paddedLen-len(data)))
		}
	}
}

func parseRequestLine(line string) (method, path, version string, ok bool) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) < 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

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
			"\r\n", contentType, contentLength)
	conn.Write([]byte(headers))
}

func writeNotFound(conn net.Conn) {
	body := "<h1>404 Not Found</h1>"
	headers := fmt.Sprintf(
		"HTTP/1.1 404 Not Found\r\n"+
			"Content-Type: text/html\r\n"+
			"Content-Length: %d\r\n\r\n%s", len(body), body)
	conn.Write([]byte(headers))
}

func writeServerError(conn net.Conn) {
	body := "<h1>500 Internal Server Error</h1>"
	headers := fmt.Sprintf(
		"HTTP/1.1 500 Internal Server Error\r\n"+
			"Content-Type: text/html\r\n"+
			"Content-Length: %d\r\n\r\n%s", len(body), body)
	conn.Write([]byte(headers))
}
