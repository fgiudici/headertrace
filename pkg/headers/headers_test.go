package headers

import (
	"reflect"
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
