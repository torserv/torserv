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

var aesgcm cipher.AEAD
var nonceSize int
var cssURLPattern = regexp.MustCompile(`url\(['"]?([^)'"]+)['"]?\)`)

func RewriteCSSLinks(css string) string {
	return cssURLPattern.ReplaceAllStringFunc(css, func(match string) string {
		matches := cssURLPattern.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		original := matches[1]
		if strings.HasPrefix(original, "data:") || strings.HasPrefix(original, "http") {
			return match
		}
		encrypted, err := EncryptPath(original)
		if err != nil {
			return match
		}
		replaced := "/" + encrypted
		return fmt.Sprintf("url('%s')", replaced)
	})
}

func Init() error {
	key := make([]byte, 32) // AES-256
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

func EncryptPath(path string) (string, error) {
	if aesgcm == nil {
		return "", errors.New("encryption not initialized")
	}
	plaintext := []byte(path)
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func DecryptPath(token string) (string, error) {
	if aesgcm == nil {
		return "", errors.New("decryption not initialized")
	}
	ciphertext, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// RewriteHTMLLinks finds href="..." and replaces internal paths with encrypted tokens.
func RewriteHTMLLinks(html string) string {
	const hrefTag = `href="`
	const srcTag = `src="`
	var result strings.Builder
	start := 0

	processTag := func(tag string) {
		for {
			tagStart := strings.Index(html[start:], tag)
			if tagStart == -1 {
				result.WriteString(html[start:])
				break
			}
			tagStart += start
			result.WriteString(html[start : tagStart+len(tag)])
			quoteEnd := strings.IndexByte(html[tagStart+len(tag):], '"')
			if quoteEnd == -1 {
				break
			}
			quoteEnd += tagStart + len(tag)
			original := html[tagStart+len(tag) : quoteEnd]
			replaced := original
			if !strings.Contains(original, "://") && !strings.HasPrefix(original, "data:") {
				if encrypted, err := EncryptPath(original); err == nil {
					replaced = "/" + encrypted
				}
			}
			result.WriteString(replaced)
			start = quoteEnd
		}
	}

	processTag(hrefTag)
	html = result.String()
	result.Reset()
	start = 0
	processTag(srcTag)
	html = result.String()

	// 👇 Add this line to process embedded <style> blocks too
	html = RewriteCSSLinks(html)

	return html
}
