package utilities

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

// GetIPAddress returns the real IP address of the HTTP request. It parses the
// X-Forwarded-For header.
func GetIPAddress(r *http.Request) string {
	if r.Header != nil {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		if xForwardedFor != "" {
			ips := strings.Split(xForwardedFor, ",")
			for i := range ips {
				ips[i] = strings.TrimSpace(ips[i])
			}

			for _, ip := range ips {
				if ip != "" {
					return ip
				}
			}
		}
	}

	ipPort := r.RemoteAddr
	ip, _, err := net.SplitHostPort(ipPort)
	if err != nil {
		return ipPort
	}

	return ip
}

// GetBodyBytes reads the whole request body properly into a byte array.
func GetBodyBytes(req *http.Request) ([]byte, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return nil, nil
	}

	originalBody := req.Body
	defer SafeClose(originalBody)

	buf, err := io.ReadAll(originalBody)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(buf))

	return buf, nil
}

func matchPattern(str string, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

func FindDuplicate(data string) []string {
	words := strings.Fields(data)
	pattern := `^".+":$`
	// Create a map to store word frequencies
	wordFrequency := make(map[string]int)

	// Iterate over the words and count their occurrences
	for _, word := range words {
		if matchPattern(word, pattern) {
			wordFrequency[word]++
		}
	}

	// Find the duplicate words
	var duplicates []string
	for word, count := range wordFrequency {
		if count > 1 {
			duplicates = append(duplicates, word)
		}
	}

	return duplicates
}
