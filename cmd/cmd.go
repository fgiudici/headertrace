package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fgiudici/headertrace/api"
	hdrs "github.com/fgiudici/headertrace/pkg/headers"
	"github.com/spf13/pflag"
)

var (
	port         string
	host         string
	headers      []string
	printVersion bool
)

func init() {
	pflag.StringVarP(&port, "port", "p", "8080", "Port to listen on")
	pflag.StringVarP(&host, "host", "", "0.0.0.0", "Host to listen on")
	pflag.StringSliceVarP(&headers, "header", "H", []string{}, "Custom response header (key:value format)")
	pflag.BoolVarP(&printVersion, "version", "v", false, "Print version and exit")
}

type server struct {
	headers map[string]string
}

// Get implements api.ServerInterface
func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	// Convert headers to map
	headers := hdrs.ToMap(r.Header)

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
	for key, value := range s.headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(http.StatusOK)

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
	}
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
		log.Fatal(err)
	}

	// Create server instance
	srv := &server{headers: customHeaders}

	// Create handler from the generated code
	handler := api.Handler(srv)

	// Start listening
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Starting server on %s\n", addr)
	return http.ListenAndServe(addr, handler)
}
