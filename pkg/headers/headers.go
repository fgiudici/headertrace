package headers

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/fgiudici/headertrace/pkg/logging"
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
// It takes a list of headers to drop and a privacy mode flag to exclude headers that may reveal
// sensitive information of the internal network. Note that enabling debug logging will log all dropped headers.
func ToMap(headers http.Header, dropHeaders []string, privMode bool) map[string]string {
	headerMap := make(map[string]string)
	normalizedDropHeaders := sliceToLower(dropHeaders)

	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if slices.Contains(normalizedDropHeaders, lowerKey) {
			logging.Debugf("Dropping header '%s':'%s'", key, strings.Join(values, ","))
			continue
		}
		if privMode {
			if isCloudflareHeader(lowerKey) || isXForwardedHeader(lowerKey) {
				logging.Debugf("Dropping header '%s':'%s' (privacy mode)", key, strings.Join(values, ","))
				continue
			}
		}
		headerMap[key] = strings.Join(values, ",")
	}
	return headerMap
}

func sliceToLower(headers []string) []string {
	lower := make([]string, len(headers))
	for i, h := range headers {
		lower[i] = strings.ToLower(h)
	}
	return lower
}

// isCloudflareHeader checks if a header is a Cloudflare-specific header that should be dropped in privacy mode.
// NOTE: it expects headers to be already normalized to lowercase.
func isCloudflareHeader(header string) bool {
	return strings.HasPrefix(header, "cf-")
}

// isXForwardedHeader checks if a header is an X-Forwarded or X-Real-IP header that should be dropped in privacy mode.
// NOTE: it expects headers to be already normalized to lowercase.
func isXForwardedHeader(header string) bool {
	return strings.HasPrefix(header, "x-forwarded-") || header == "x-real-ip"
}

// GetRemoteHostInfo extracts the remote host information from the request, inspecting common proxy headers like CF-Connecting-IP, X-Real-IP, and X-Forwarded-For.
// It returns a formatted string with the remote address and user agent.
func GetRemoteHostInfo(r *http.Request) string {
	// Example of received headers:
	// "Accept": "*/*",
	// "Accept-Encoding": "gzip",
	// "Cdn-Loop": "cloudflare; loops=1",
	// "Cf-Connecting-Ip": "1.2.3.4",
	// "Cf-Ipcountry": "IT",
	// "Cf-Ray": "9cbdc3515d22baf3-MXP",
	// "Cf-Visitor": "{\"scheme\":\"http\"}",
	// "User-Agent": "curl/7.88.1",
	// "X-Forwarded-For": "10.22.0.0",
	// "X-Forwarded-Host": "headers.example.com",
	// "X-Forwarded-Port": "80",
	// "X-Forwarded-Proto": "http",
	// "X-Forwarded-Server": "traefik-73f98ac65-z1drx",
	// "X-Real-Ip": "10.22.0.0"

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
