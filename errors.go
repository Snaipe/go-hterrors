// Copyright 2020, Franklin "Snaipe" Mathieu <me@snai.pe>
//
// Use of this source-code is govered by the MIT license, which
// can be found in the LICENSE file.

package hterrors

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/k3a/html2text"
)

// StatusError represents a non-2xx HTTP status code, and the associated message
// returned by the server, if any.
type StatusError struct {
	StatusCode int
	Message    string
}

func (err *StatusError) Error() string {
	code := fmt.Sprintf("%d", err.StatusCode)
	text := http.StatusText(err.StatusCode)

	if strings.Contains(err.Message, code) && strings.Contains(err.Message, text) {
		return err.Message
	}
	return fmt.Sprintf("%s %s: %s", code, text, err.Message)
}

var (
	nlre  = regexp.MustCompile(`(\r?\n)+`)
	space = regexp.MustCompile(`\s\s+`)
)

func extractMessage(resp *http.Response) string {
	mtype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		// assume text
		mtype = "text/plain"
	}

	// The MIME type might be a vendor type, which looks like application/vnd.*+type;
	// in which case we try to change it back to the appropriate application/type
	// MIME.
	// This isn't always correct, but is a good enough heuristic for most API
	// bodies.
	if strings.HasPrefix(mtype, "application/vnd.") {
		if i := strings.IndexRune(mtype, '+'); i != -1 {
			mtype = "application/" + mtype[i+1:]
		}
	}

	switch mtype {
	case "text/plain":
		var out strings.Builder
		io.Copy(&out, resp.Body)
		return out.String()

	case "text/html":
		var out strings.Builder
		io.Copy(&out, resp.Body)
		body := strings.TrimSpace(html2text.HTML2Text(out.String()))
		return space.ReplaceAllString(nlre.ReplaceAllString(body, ": "), " ")

	case "application/json":
		var doc map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
			return fmt.Sprintf("<invalid json in response body: %v>", err)
		}

		keys := make([]string, 0, len(doc))
		for k := range doc {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fields := make([]string, 0, len(doc))
		for _, k := range keys {
			fields = append(fields, fmt.Sprintf("%s: %v", k, doc[k]))
		}
		return strings.Join(fields, ", ")

	default:
		return ""
	}
}

// CheckStatus returns an error if the status code of the specified response
// is not in the 2xx family. The returned error contains a digested version
// of the response body, and the reponse body is consumed.
func CheckStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	msg := extractMessage(resp)
	//if prefix := fmt.Sprintf("%d %s", err.StatusCode, http.StatusText(err.StatusCode)); strings.HasPrefix(msg, prefix) {
	//} else if prefix := fmt.Sprintf("%d %s", err.StatusCode, http.StatusText(err.StatusCode)); strings.HasPrefix(msg, prefix) {
	//}
	//	fallthrough
	//case msg == fmt.Sprintf("%s", http.StatusText(err.StatusCode))
	//	msg = 
	//}
	//if msg == fmt.Sprintf("%d %s")

	return &url.Error{
		Op:  resp.Request.Method,
		URL: resp.Request.URL.String(),
		Err: &StatusError{StatusCode: resp.StatusCode, Message: msg},
	}
}

// Check is a convenience wrapper over CheckStatus -- if the passed error is
// non-nil, it is returned; otherwise, CheckStatus(resp) is returned.
//
// This function exists to make it easier to write error handling code, by
// directly taking the expression that makes the request. For instance:
//
//     Check(http.Get("http://example.com"))
//
//     Check(http.Do(request))
//
func Check(resp *http.Response, err error) (*http.Response, error) {
	if err == nil {
		err = CheckStatus(resp)
	}
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil, err
	}
	return resp, nil
}
