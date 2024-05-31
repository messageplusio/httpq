package types

import (
	"database/sql/driver"
	"encoding/json"
)

type HTTPRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (r *HTTPRequest) Scan(src interface{}) error {
	// handle nil value and case when src is not []byte
	if src == nil {
		return nil
	}
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, r)
	default:
		return nil
	}
}

func (r *HTTPRequest) Value() (driver.Value, error) {
	// handle nil value
	if r == nil {
		return "{}", nil
	}
	return json.Marshal(r)
}
