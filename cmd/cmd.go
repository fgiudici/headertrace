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
	logLevel     string
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
	pflag.StringVarP(&logLevel, "log-level", "l", "", "Logging level: TRACE, DEBUG, INFO, WARN, ERROR (overrides the LOG_LEVEL env variable)")
}

type server struct {
	headers     map[string]string
	dropHeaders []string
	privMode    bool
	sentHeaders bool
}

// Get implements api.ServerInterface
func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	logging.Infof("Received request: %s", hdrs.GetRemoteHostInfo(r))

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
		logging.Tracef("Dumping sent headers to response body")
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

	// Init logging level
	if logLevel != "" {
		os.Setenv("LOG_LEVEL", logLevel)
	}
	if err := logging.Init(logLevel); err != nil {
		logging.Fatalf("Failed to initialize logging: %v", err)
	}

	logging.Infof("Starting HeaderTrace version %s", getVersion())

	// Parse custom headers
	customHeaders, err := hdrs.SliceToMap(headers)
	if err != nil {
		logging.Fatalf("Custom headers: %v", err)
	}
	if len(customHeaders) > 0 {
		logging.Debugf("Custom headers to add in responses: %v", customHeaders)
	}

	if len(dropHeaders) > 0 {
		logging.Debugf("Headers to drop from echoed request headers: %v", dropHeaders)
	}

	logging.Debugf("Privacy mode: %v", privMode)
	logging.Debugf("Dump sent headers: %v", sentHeaders)

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
