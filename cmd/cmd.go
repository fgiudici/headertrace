package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/fgiudici/headertrace/api"
)

var (
	port = flag.String("port", "8080", "Port to listen on")
	host = flag.String("host", "0.0.0.0", "Host to listen on")
)

type server struct{}

// Get implements api.ServerInterface
func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	// Convert headers to map
	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ",")
	}

	// Determine protocol version
	protocol := r.Proto
	if protocol == "" {
		protocol = "HTTP/1.1"
	}

	// Create the response
	response := api.HeaderResponse{
		Headers:  headers,
		Host:     r.Host,
		Method:   r.Method,
		Path:     r.RequestURI,
		Protocol: protocol,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
	}
}

// Execute starts the HTTP server
func Execute() error {
	flag.Parse()

	// Create server instance
	srv := &server{}

	// Create handler from the generated code
	handler := api.Handler(srv)

	// Start listening
	addr := fmt.Sprintf("%s:%s", *host, *port)
	fmt.Printf("Starting server on %s\n", addr)
	return http.ListenAndServe(addr, handler)
}
