package headers

import (
	"fmt"
	"net/http"
	"strings"
)

// Slice2Map takes a slice of header strings in "key:value" format and returns a map.
// Returns an error if any header has an invalid format.
func SliceToMap(headerStrings []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, h := range headerStrings {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format '%s', expected 'key:value'", h)
		}
		if parts[0] == "" {
			return nil, fmt.Errorf("header key cannot be empty in '%s'", h)
		}
		headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return headers, nil
}

// ToMap converts an http.Header to a "key:value" map.
func ToMap(headers http.Header) map[string]string {
	headerMap := make(map[string]string)
	for key, values := range headers {
		headerMap[key] = strings.Join(values, ",")
	}
	return headerMap
}

func RemoteHostInfo(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")

	// Proxied through Cloudflare?
	if remote := r.Header.Get("CF-Connecting-IP"); remote != "" {
		remoteAddr = fmt.Sprintf("%s (%s)", remote, r.Header.Get("Cf-Ipcountry"))
	} else if remote := r.Header.Get("X-Real-Ip"); remote != "" {
		remoteAddr = remote
	} else if remote := r.Header.Get("X-Forwarded-For"); remote != "" {
		remoteAddr = remote
	}

	return fmt.Sprintf("%s %q - %s %s %q", remoteAddr, userAgent, r.Method, r.Proto, r.URL.String())
}
