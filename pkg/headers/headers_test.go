package headers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestSliceToMap(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "single header",
			input: []string{"X-Custom:value"},
			want:  map[string]string{"X-Custom": "value"},
		},
		{
			name:  "multiple headers",
			input: []string{"X-Custom:value", "Content-Type:application/json"},
			want: map[string]string{
				"X-Custom":     "value",
				"Content-Type": "application/json",
			},
		},
		{
			name:  "header with spaces",
			input: []string{" X-Custom : value "},
			want:  map[string]string{"X-Custom": "value"},
		},
		{
			name:  "value with colons",
			input: []string{"Authorization:Bearer token:secret"},
			want:  map[string]string{"Authorization": "Bearer token:secret"},
		},
		{
			name:    "missing colon",
			input:   []string{"InvalidHeader"},
			wantErr: true,
		},
		{
			name:    "empty key",
			input:   []string{":value"},
			wantErr: true,
			want:    map[string]string{"": "value"},
		},
		{
			// Empty value is allowed by RFC 7230 spec.
			name:  "empty value",
			input: []string{"X-Empty:"},
			want:  map[string]string{"X-Empty": ""},
		},
		{
			name:    "mixed valid and invalid",
			input:   []string{"Valid:header", "InvalidHeader"},
			wantErr: true,
		},
		{
			name:  "empty input slice",
			input: []string{},
			want:  map[string]string{},
		},
		{
			name:    "whitespace only key and value",
			input:   []string{"   :   "},
			wantErr: true,
		},
		{
			name:  "duplicate header keys",
			input: []string{"X-Custom:first", "X-Custom:second"},
			want:  map[string]string{"X-Custom": "second"},
		},
		{
			name:  "case sensitive keys",
			input: []string{"x-custom:value1", "X-Custom:value2"},
			want: map[string]string{
				"x-custom": "value1",
				"X-Custom": "value2",
			},
		},
		{
			name:  "special characters in value",
			input: []string{"X-Custom:!@#$%^&*()"},
			want:  map[string]string{"X-Custom": "!@#$%^&*()"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SliceToMap(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("SliceToMap() unexpected error = %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatalf("SliceToMap() expected error, got nil (%v)", tt.input)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SliceToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		dropHeaders []string
		privMode    bool
		want        map[string]string
	}{
		{
			name: "basic headers",
			headers: http.Header{
				"X-Custom": {"value"},
			},
			dropHeaders: []string{},
			privMode:    false,
			want: map[string]string{
				"X-Custom": "value",
			},
		},
		{
			name: "drop headers",
			headers: http.Header{
				"X-Custom": {"value"},
			},
			dropHeaders: []string{"x-custom"},
			privMode:    false,
			want:        map[string]string{},
		},
		{
			name: "privMode drops Cloudflare headers",
			headers: http.Header{
				"CF-Ray":           {"9cbdc3515d22baf3-MXP"},
				"Cf-Visitor":       {"{\"scheme\":\"https\"}"},
				"cf-Connecting-Ip": {"10.22.0.0"},
				"X-Custom":         {"value"},
			},
			dropHeaders: []string{},
			privMode:    true,
			want: map[string]string{
				"X-Custom": "value",
			},
		},
		{
			name: "privMode drops X-Forwarded headers",
			headers: http.Header{
				"X-Forwarded-For":   {"10.22.0.0"},
				"X-forwarded-Host":  {"example.com"},
				"x-forwarded-proto": {"https"},
				"X-Custom":          {"value"},
			},
			dropHeaders: []string{},
			privMode:    true,
			want: map[string]string{
				"X-Custom": "value",
			},
		},
		{
			name: "privMode drops X-Real-IP header",
			headers: http.Header{
				"X-Real-Ip": {"10.22.0.0"},
				"x-Real-Ip": {"10.22.0.0"},
				"X-real-ip": {"10.22.0.0"},
				"X-Custom":  {"value"},
			},
			dropHeaders: []string{},
			privMode:    true,
			want: map[string]string{
				"X-Custom": "value",
			},
		},
		{
			name: "privMode false keeps all headers",
			headers: http.Header{
				"CF-Ray":          {"9cbdc3515d22baf3-MXP"},
				"X-Forwarded-For": {"10.22.0.0"},
				"X-Real-Ip":       {"10.22.0.0"},
				"X-Custom":        {"value"},
			},
			dropHeaders: []string{},
			privMode:    false,
			want: map[string]string{
				"CF-Ray":          "9cbdc3515d22baf3-MXP",
				"X-Forwarded-For": "10.22.0.0",
				"X-Real-Ip":       "10.22.0.0",
				"X-Custom":        "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToMap(tt.headers, tt.dropHeaders, tt.privMode)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRemoteHostInfo(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		headers    http.Header
		method     string
		urlString  string
		expectedIP string
	}{
		{
			name:       "uses CF-Connecting-IP with Cf-Ipcountry",
			remoteAddr: "127.0.0.1:5000",
			headers: http.Header{
				"CF-Connecting-IP": {"1.2.3.4"},
				"Cf-Ipcountry":     {"US"},
				"X-Real-Ip":        {"5.6.7.8"},
				"X-Forwarded-For":  {"9.10.11.12"},
			},
			method:     "GET",
			urlString:  "http://example.com/",
			expectedIP: "1.2.3.4(US)",
		},
		{
			name:       "uses X-Real-IP when CF-Connecting-IP not available",
			remoteAddr: "127.0.0.1:5000",
			headers: http.Header{
				"X-Real-Ip":       {"5.6.7.8"},
				"X-Forwarded-For": {"9.10.11.12"},
			},
			method:     "GET",
			urlString:  "http://example.com/",
			expectedIP: "5.6.7.8",
		},
		{
			name:       "uses X-Forwarded-For when CF-Connecting-IP and X-Real-IP not available",
			remoteAddr: "127.0.0.1:5000",
			headers: http.Header{
				"X-Forwarded-For": {"9.10.11.12"},
			},
			method:     "GET",
			urlString:  "http://example.com/",
			expectedIP: "9.10.11.12",
		},
		{
			name:       "uses r.RemoteAddr when no proxy headers available",
			remoteAddr: "192.168.1.1:8080",
			headers:    http.Header{},
			method:     "GET",
			urlString:  "http://example.com/",
			expectedIP: "192.168.1.1:8080",
		},
		{
			name:       "prefers cf-Connecting-IP over X-Real-IP",
			remoteAddr: "127.0.0.1:5000",
			headers: http.Header{
				"cf-Connecting-IP": {"1.2.3.4"},
				"cf-Ipcountry":     {"IT"},
				"X-Real-Ip":        {"5.6.7.8"},
				"X-Forwarded-For":  {"9.10.11.12"},
			},
			method:     "GET",
			urlString:  "http://example.com/",
			expectedIP: "1.2.3.4(IT)",
		},
		{
			name:       "prefers x-Real-ip over X-Forwarded-For when CF-Connecting-IP missing",
			remoteAddr: "127.0.0.1:5000",
			headers: http.Header{
				"x-Real-ip":       {"5.6.7.8"},
				"X-Forwarded-For": {"9.10.11.12"},
			},
			method:     "GET",
			urlString:  "http://example.com/",
			expectedIP: "5.6.7.8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.urlString, nil)
			req.RemoteAddr = tt.remoteAddr
			for key, values := range tt.headers {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}

			got := GetRemoteHostInfo(req)

			// Verify that the expected IP is contained in the result
			if !strings.Contains(got, tt.expectedIP) {
				t.Fatalf("GetRemoteHostInfo() = %q, expected to contain IP %q", got, tt.expectedIP)
			}

			// Verify the format includes method and proto
			if !strings.Contains(got, tt.method) {
				t.Fatalf("GetRemoteHostInfo() = %q, expected to contain method %q", got, tt.method)
			}
		})
	}
}
