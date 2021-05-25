package errors

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	return Annotate(err, WithCode(Internal), WithMessage(req.URL.String()))
}

func (e *RoundTripper) onBusinessError(req *http.Request, resp *http.Response, err error) error {
	defer func() { _ = resp.Body.Close() }()

	if e.isJson(resp) {
		return Decode(resp.Body)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return e.onNetworkError(req, err)
	}

	return Annotate(
		fmt.Errorf("%s", data),
		WithHttpCode(resp.StatusCode))
}

func (e *RoundTripper) isSuccess(resp *http.Response) bool {
	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest
}

func (e *RoundTripper) isJson(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/json")
}
