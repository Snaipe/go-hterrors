package hterrors

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestJSONResponse(t *testing.T) {
	err := CheckStatus(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(strings.NewReader(`{"foo": "bar", "baz": "quux"}`)),
		Request: &http.Request{
			Method: "GET",
			URL:    &url.URL{},
		},
	})
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
	if !strings.Contains(err.Error(), "foo: bar") {
		t.Errorf("Expected error message to contain fields from JSON response; instead got %q", err.Error())
	}
}
