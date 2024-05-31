package types

import (
	"database/sql/driver"
	"reflect"
	"testing"
)

func TestHTTPRequest_Scan(t *testing.T) {
	tests := []struct {
		name string
		src  interface{}
		want *HTTPRequest
	}{
		{
			name: "Valid JSON",
			src:  []byte(`{"url":"https://example.com","method":"GET","headers":{"Content-Type":"application/json"},"body":"{\"key\":\"value\"}"}`),
			want: &HTTPRequest{
				URL:     "https://example.com",
				Method:  "GET",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    `{"key":"value"}`,
			},
		},
		{
			name: "Empty JSON",
			src:  []byte(`{}`),
			want: &HTTPRequest{
				URL:     "",
				Method:  "",
				Headers: nil,
				Body:    "",
			},
		},
		{
			name: "Nil value",
			src:  nil,
			want: &HTTPRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HTTPRequest{}
			err := r.Scan(tt.src)
			if err != nil {
				t.Errorf("Scan() error = %v, want nil", err)
			}
			if !reflect.DeepEqual(r, tt.want) {
				t.Errorf("Scan() got = %v, want %v", r, tt.want)
			}
		})
	}
}

func TestHTTPRequest_Value(t *testing.T) {
	tests := []struct {
		name    string
		request *HTTPRequest
		want    driver.Value
		wantErr bool
	}{
		{
			name: "Non-nil request",
			request: &HTTPRequest{
				URL:     "https://example.com",
				Method:  "GET",
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    `{"key":"value"}`,
			},
			want:    []byte(`{"url":"https://example.com","method":"GET","headers":{"Content-Type":"application/json"},"body":"{\"key\":\"value\"}"}`),
			wantErr: false,
		},
		{
			name:    "Nil request",
			request: nil,
			want:    "{}",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.request
			got, err := r.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}
