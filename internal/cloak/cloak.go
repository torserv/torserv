package cloak

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Global AEAD cipher for encryption and its nonce size
var aesgcm cipher.AEAD
var nonceSize int

// Regular expression pattern to match CSS url() values
var cssURLPattern = regexp.MustCompile(`url\(['"]?([^)'"]+)['"]?\)`)

// RewriteCSSLinks rewrites relative CSS url() paths by encrypting them.
// It skips URLs that are data URIs or start with http/https.
func RewriteCSSLinks(css string) string {
	return cssURLPattern.ReplaceAllStringFunc(css, func(match string) string {
		// Extract the URL inside the url(...) function
		matches := cssURLPattern.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match // Return unchanged if the match is malformed
		}
		original := matches[1]

		// Skip rewriting if the URL is a data URI or absolute HTTP(S) link
		if strings.HasPrefix(original, "data:") || strings.HasPrefix(original, "http://") || strings.HasPrefix(original, "https://") {
			return match
		}

		// Encrypt the relative path
		encrypted, err := EncryptPath(original)
		if err != nil {
			return match // If encryption fails, return the original match
		}

		// Replace with encrypted path
		replaced := "/" + encrypted
		return fmt.Sprintf("url('%s')", replaced)
	})
}

// Init initializes the AES-GCM cipher with a randomly generated 256-bit key.
// It also stores the nonce size required for future encryption/decryption.
func Init() error {
	key := make([]byte, 32) // Generate a 256-bit key for AES-256
	if _, err := rand.Read(key); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize = aesgcm.NonceSize()
	return nil
}

// EncryptPath encrypts the provided path string using AES-GCM and returns
// the result encoded in base64 (URL-safe, no padding).
func EncryptPath(path string) (string, error) {
	if aesgcm == nil {
		return "", errors.New("encryption not initialized")
	}

	plaintext := []byte(path)

	// Generate a random nonce for this encryption
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal appends the ciphertext to the nonce and encrypts the data
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)

	// Encode to base64 (URL-safe, no padding)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// DecryptPath takes a base64-encoded encrypted string and decrypts it
// back into the original plaintext path.
func DecryptPath(token string) (string, error) {
	if aesgcm == nil {
		return "", errors.New("decryption not initialized")
	}

	// Decode from base64
	ciphertext, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}

	// Separate nonce and actual ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// RewriteHTMLLinks finds and rewrites href="..." and src="..." attributes in the HTML string.
// Internal paths are encrypted and replaced with tokenized URLs.
// It also processes embedded <style> blocks for CSS path rewriting.
func RewriteHTMLLinks(html string) string {
	const hrefTag = `href="`
	const srcTag = `src="`
	var result strings.Builder
	start := 0

	// processTag scans for a given HTML tag pattern (e.g., href=" or src="),
	// encrypts internal paths, and replaces them with encrypted tokens.
	processTag := func(tag string) {
		for {
			// Find the start of the next tag
			tagStart := strings.Index(html[start:], tag)
			if tagStart == -1 {
				result.WriteString(html[start:])
				break
			}
			tagStart += start

			// Write up to and including the tag
			result.WriteString(html[start : tagStart+len(tag)])

			// Find the end quote of the URL value
			quoteEnd := strings.IndexByte(html[tagStart+len(tag):], '"')
			if quoteEnd == -1 {
				break // malformed HTML
			}
			quoteEnd += tagStart + len(tag)

			// Extract the original path value
			original := html[tagStart+len(tag) : quoteEnd]
			replaced := original

			// Encrypt only internal paths (not full URLs or data URIs)
			if !strings.Contains(original, "://") && !strings.HasPrefix(original, "data:") {
				if encrypted, err := EncryptPath(original); err == nil {
					replaced = "/" + encrypted
				}
			}

			// Write the encrypted or original path
			result.WriteString(replaced)

			// Move the start pointer forward
			start = quoteEnd
		}
	}

	// Process href and src attributes
	processTag(hrefTag)
	html = result.String()
	result.Reset()
	start = 0
	processTag(srcTag)
	html = result.String()

	// Process embedded CSS in <style> blocks
	html = RewriteCSSLinks(html)

	return html
}
