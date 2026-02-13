package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fgiudici/headertrace/api"
	hdrs "github.com/fgiudici/headertrace/pkg/headers"
	"github.com/fgiudici/headertrace/pkg/logging"
	"github.com/spf13/pflag"
)

var (
	port         string
	host         string
	headers      []string
	dropHeaders  []string
	sentHeaders  bool
	privMode     bool
	printVersion bool
)

func init() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "HeaderTrace %s - A simple HTTP server that echoes back received HTTP headers\n\n", getVersion())
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n\n", filepath.Base(os.Args[0]))
		pflag.PrintDefaults()
	}

	pflag.StringVarP(&host, "address", "a", "0.0.0.0", "IP address (or domain) to bind to")
	pflag.StringVarP(&port, "port", "p", "8080", "TCP port to bind to")
	pflag.StringSliceVarP(&headers, "header", "H", []string{}, "Custom HTTP headers to add to responses (key1:value1,key2:value2)")
	pflag.StringSliceVarP(&dropHeaders, "drop-header", "D", []string{}, "HTTP headers to redact from request headers echoed in the response body (key1,key2)")
	pflag.BoolVarP(&privMode, "privacy", "P", false, "Drop X-Forwarded and Cloudflare headers from request headers echoed in the response body")
	pflag.BoolVarP(&sentHeaders, "sent", "s", false, "Dump the HTTP headers added in the response in the response body")
	pflag.BoolVarP(&printVersion, "version", "v", false, "Print version and exit")
}

type server struct {
	headers     map[string]string
	dropHeaders []string
	privMode    bool
	sentHeaders bool
}

// Get implements api.ServerInterface
func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	logging.Infof("Received request: %s", hdrs.RemoteHostInfo(r))

	// Convert headers to map
	headers := hdrs.ToMap(r.Header, s.dropHeaders, s.privMode)
	var xHeadersPtr *map[string]string

	protocol := r.Proto
	if protocol == "" {
		protocol = "HTTP/1.1"
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	for key, value := range s.headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(http.StatusOK)

	if s.sentHeaders {
		xHeaders := hdrs.ToMap(w.Header(), nil, false)
		xHeadersPtr = &xHeaders
	}

	// Create the response
	response := api.HeaderResponse{
		Headers:  headers,
		Host:     r.Host,
		Method:   r.Method,
		Path:     r.RequestURI,
		Protocol: protocol,
		Sent:     xHeadersPtr,
	}

	// Encode and send the response
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(response); err != nil {
		logging.Errorf("Error encoding response: %v", err)
	}
}

func (s *server) GetMatchall(w http.ResponseWriter, r *http.Request, matchall string) {
	s.Get(w, r) // Reuse the same logic for all paths
}

// Execute starts the HTTP server
func Execute() error {
	pflag.Parse()
	if printVersion {
		fmt.Println(getVersion())
		return nil
	}
	// Parse custom headers
	customHeaders, err := hdrs.SliceToMap(headers)
	if err != nil {
		logging.Fatalf("%v", err)
	}

	// Create server instance
	srv := &server{headers: customHeaders,
		dropHeaders: dropHeaders,
		privMode:    privMode,
		sentHeaders: sentHeaders}

	// Create handler from the generated code
	handler := api.Handler(srv)

	// Start listening
	addr := fmt.Sprintf("%s:%s", host, port)
	logging.Infof("Starting server on %s", addr)
	return http.ListenAndServe(addr, handler)
}
