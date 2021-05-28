package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type RoundTripper struct {
	Parent http.RoundTripper
}

func (e *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if resp, err = e.next(req); err != nil {
		return
	}

	if e.isSuccess(resp) {
		return
	}

	err = e.onError(resp, err)
	resp = nil
	return
}

func (e *RoundTripper) next(req *http.Request) (resp *http.Response, err error) {
	rt := http.DefaultTransport
	if e.Parent != nil {
		rt = e.Parent
	}
	return rt.RoundTrip(req)
}

func (e *RoundTripper) onError(resp *http.Response, err error) error {
	defer func() { _ = resp.Body.Close() }()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	if e.isJson(resp) {
		dec := NewDecoder(json.NewDecoder(&buf))
		return dec.Decode()
	}

	msg := strings.TrimSpace(buf.String())
	return fmt.Errorf("%s: %w", msg, Unknown)
}

func (e *RoundTripper) isSuccess(r *http.Response) bool {
	return r.StatusCode > 199 && r.StatusCode < 300
}

func (e *RoundTripper) isJson(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/json")
}
