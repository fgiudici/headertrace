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
	for key, values := range headers {
		if slices.Contains(dropHeaders, key) {
			logging.Debugf("Dropping header '%s':'%s'", key, strings.Join(values, ","))
			continue
		}
		if privMode {
			if isCloudflareHeader(key) || isXForwardedHeader(key) {
				logging.Debugf("Dropping header '%s':'%s' (privacy mode)", key, strings.Join(values, ","))
				continue
			}
		}
		headerMap[key] = strings.Join(values, ",")
	}
	return headerMap
}

func isCloudflareHeader(header string) bool {
	return strings.HasPrefix(header, "CF-") || strings.HasPrefix(header, "Cf-")
}

func isXForwardedHeader(header string) bool {
	return strings.HasPrefix(header, "X-Forwarded-") || header == "X-Real-Ip"
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
