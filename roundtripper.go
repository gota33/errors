package errors

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type RoundTripper struct {
	Parent http.RoundTripper
}

func (e *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if resp, err = e.next(req); err != nil {
		err = e.onNetworkError(req, err)
		return
	}

	if e.isSuccess(resp) {
		return
	}

	err = e.onBusinessError(req, resp, err)
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

func (e *RoundTripper) onNetworkError(req *http.Request, err error) error {
	var (
		annotations = []Annotation{WithCode(Internal)}
		uErr        *url.Error
	)
	if !errors.As(err, &uErr) {
		msg := fmt.Sprintf("%s %s", req.Method, req.URL)
		annotations = append(annotations, WithMessage(msg))
	}
	return Annotate(err, WithCode(Internal))
}

func (e *RoundTripper) onBusinessError(req *http.Request, resp *http.Response, err error) error {
	defer func() { _ = resp.Body.Close() }()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return e.onNetworkError(req, err)
	}

	if e.isJson(resp) {
		if err = Decode(&buf); err != nil {
			return e.onNetworkError(req, err)
		}
	}

	return Annotate(
		errors.New(buf.String()),
		WithHttpCode(resp.StatusCode))
}

func (e *RoundTripper) isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest
}

func (e *RoundTripper) isJson(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/json")
}
