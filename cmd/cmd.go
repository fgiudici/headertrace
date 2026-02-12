package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fgiudici/headertrace/api"
	hdrs "github.com/fgiudici/headertrace/pkg/headers"
	"github.com/fgiudici/headertrace/pkg/logging"
	"github.com/spf13/pflag"
)

var (
	port         string
	host         string
	headers      []string
	sentHeaders  bool
	printVersion bool
)

func init() {
	pflag.StringVarP(&host, "address", "a", "0.0.0.0", "IP address (or domain) to bind to")
	pflag.StringVarP(&port, "port", "p", "8080", "TCP port to bind to")
	pflag.StringSliceVarP(&headers, "header", "H", []string{}, "Custom HTTP headers to add to the HTTP responses (key:value format)")
	pflag.BoolVarP(&sentHeaders, "sent", "s", false, "Include the original HTTP headers added to the response in the body")
	pflag.BoolVarP(&printVersion, "version", "v", false, "Print version and exit")
}

type server struct {
	headers     map[string]string
	sentHeaders bool
}

// Get implements api.ServerInterface
func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	logging.Infof("Received request: %s", hdrs.RemoteHostInfo(r))

	// Convert headers to map
	headers := hdrs.ToMap(r.Header)
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
		xHeaders := hdrs.ToMap(w.Header())
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
	srv := &server{headers: customHeaders, sentHeaders: sentHeaders}

	// Create handler from the generated code
	handler := api.Handler(srv)

	// Start listening
	addr := fmt.Sprintf("%s:%s", host, port)
	logging.Infof("Starting server on %s", addr)
	return http.ListenAndServe(addr, handler)
}
