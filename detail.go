package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"strings"
)

const (
	typeUrlPrefix              = "type.googleapis.com/google.rpc."
	TypeUrlDebugInfo           = typeUrlPrefix + "DebugInfo"
	TypeUrlResourceInfo        = typeUrlPrefix + "ResourceInfo"
	TypeUrlBadRequest          = typeUrlPrefix + "BadRequest"
	TypeUrlPreconditionFailure = typeUrlPrefix + "PreconditionFailure"
	TypeUrlErrorInfo           = typeUrlPrefix + "ErrorInfo"
	TypeUrlQuotaFailure        = typeUrlPrefix + "QuotaFailure"
	TypeUrlRequestInfo         = typeUrlPrefix + "RequestInfo"
	TypeUrlHelp                = typeUrlPrefix + "Help"
	TypeUrlLocalizedMessage    = typeUrlPrefix + "LocalizedMessage"
)

type Any interface {
	TypeUrl() string
}

type AnyDetail map[string]interface{}

func (d AnyDetail) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(d))
}

func (d AnyDetail) TypeUrl() string {
	return fmt.Sprintf("%v", d["@type"])
}

type DebugInfo struct {
	StackEntries []string `json:"stackEntries,omitempty"`
	Detail       string   `json:"detail,omitempty"`
}

func (d DebugInfo) TypeUrl() string     { return TypeUrlDebugInfo }
func (d DebugInfo) Annotate(m Modifier) { m.AppendDetails(d) }

func (d DebugInfo) MarshalJSON() ([]byte, error) {
	type payload DebugInfo
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d DebugInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s: %q\n", "detail", d.Detail)
			_, _ = fmt.Fprintf(f, "%s:\n", "stack")
			for _, entry := range d.StackEntries {
				_, _ = fmt.Fprintf(f, "\t%s\n", entry)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(f, d.Detail)
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", d.Detail)
	}
}

type ResourceInfo struct {
	ResourceType string `json:"resourceType,omitempty"`
	ResourceName string `json:"resourceName,omitempty"`
	Owner        string `json:"owner,omitempty"`
	Description  string `json:"description,omitempty"`
}

func (d ResourceInfo) TypeUrl() string     { return TypeUrlResourceInfo }
func (d ResourceInfo) Annotate(m Modifier) { m.AppendDetails(d) }

func (d ResourceInfo) MarshalJSON() ([]byte, error) {
	type payload ResourceInfo
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d ResourceInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s: %q\n", "resource_type", d.ResourceType)
			_, _ = fmt.Fprintf(f, "%s: %q\n", "resource_name", d.ResourceName)
			_, _ = fmt.Fprintf(f, "%s: %q\n", "owner", d.Owner)
			_, _ = fmt.Fprintf(f, "%s: %q\n", "description", d.Description)
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f,
			"resource type: %s, name: %s, owner: %s, description: %s",
			d.ResourceType, d.ResourceName, d.Owner, d.Description)
	}
}

type FieldViolation struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
}

type BadRequest struct {
	FieldViolations []FieldViolation `json:"fieldViolations,omitempty"`
}

func (d BadRequest) TypeUrl() string     { return TypeUrlBadRequest }
func (d BadRequest) Annotate(m Modifier) { m.AppendDetails(d) }

func (d BadRequest) MarshalJSON() ([]byte, error) {
	type payload BadRequest
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d BadRequest) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s:\n", "field_violations")
			for _, v := range d.FieldViolations {
				_, _ = fmt.Fprintf(f, "\t%s: %q\n", v.Field, v.Description)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "field violations: %d", len(d.FieldViolations))
	}
}

type TypedViolation struct {
	Type        string `json:"type,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
}

type PreconditionFailure struct {
	Violations []TypedViolation `json:"violations,omitempty"`
}

func (d PreconditionFailure) TypeUrl() string     { return TypeUrlPreconditionFailure }
func (d PreconditionFailure) Annotate(m Modifier) { m.AppendDetails(d) }

func (d PreconditionFailure) MarshalJSON() ([]byte, error) {
	type payload PreconditionFailure
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d PreconditionFailure) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s:\n", "violations")
			for _, v := range d.Violations {
				_, _ = fmt.Fprintf(f, "\t[%s] %s: %q\n", v.Type, v.Subject, v.Description)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "violations: %d", len(d.Violations))
	}
}

type ErrorInfo struct {
	Reason   string            `json:"reason,omitempty"`
	Domain   string            `json:"domain,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (d ErrorInfo) TypeUrl() string     { return TypeUrlErrorInfo }
func (d ErrorInfo) Annotate(m Modifier) { m.AppendDetails(d) }

func (d ErrorInfo) MarshalJSON() ([]byte, error) {
	type payload ErrorInfo
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d ErrorInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s: %q\n", "reason", d.Reason)
			_, _ = fmt.Fprintf(f, "%s: %q\n", "domain", d.Domain)
			_, _ = fmt.Fprintf(f, "%s:\n", "metadata")
			for k, v := range d.Metadata {
				_, _ = fmt.Fprintf(f, "\t%s: %q\n", k, v)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "[%s] %s", d.Domain, d.Reason)
	}
}

type Violation struct {
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
}

type QuotaFailure struct {
	Violations []Violation `json:"violations,omitempty"`
}

func (d QuotaFailure) TypeUrl() string     { return TypeUrlQuotaFailure }
func (d QuotaFailure) Annotate(m Modifier) { m.AppendDetails(d) }

func (d QuotaFailure) MarshalJSON() ([]byte, error) {
	type payload QuotaFailure
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d QuotaFailure) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s:\n", "violations")
			for _, v := range d.Violations {
				_, _ = fmt.Fprintf(f, "\t%s: %q\n", v.Subject, v.Description)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "violations: %d", len(d.Violations))
	}
}

type RequestInfo struct {
	RequestId   string `json:"requestId,omitempty"`
	ServingData string `json:"servingData,omitempty"`
}

func (d RequestInfo) TypeUrl() string     { return TypeUrlRequestInfo }
func (d RequestInfo) Annotate(m Modifier) { m.AppendDetails(d) }

func (d RequestInfo) MarshalJSON() ([]byte, error) {
	type payload RequestInfo
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d RequestInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s: %q\n", "request_id", d.RequestId)
			_, _ = fmt.Fprintf(f, "%s: %q\n", "serving_data", d.ServingData)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(f, d.RequestId)
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", d.RequestId)
	}
}

type Link struct {
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
}

type Help struct {
	Links []Link `json:"links,omitempty"`
}

func (d Help) TypeUrl() string     { return TypeUrlHelp }
func (d Help) Annotate(m Modifier) { m.AppendDetails(d) }

func (d Help) MarshalJSON() ([]byte, error) {
	type payload Help
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d Help) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s:\n", "links")
			for _, v := range d.Links {
				_, _ = fmt.Fprintf(f, "\t%s: %q\n", v.Description, v.Url)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "help(%d)", len(d.Links))
	}
}

type LocalizedMessage struct {
	Local   string `json:"local,omitempty"`
	Message string `json:"message,omitempty"`
}

func (d LocalizedMessage) TypeUrl() string     { return TypeUrlLocalizedMessage }
func (d LocalizedMessage) Annotate(m Modifier) { m.AppendDetails(d) }

func (d LocalizedMessage) MarshalJSON() ([]byte, error) {
	type payload LocalizedMessage
	return json.Marshal(struct {
		Type string `json:"@type"`
		payload
	}{d.TypeUrl(), payload(d)})
}

func (d LocalizedMessage) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "%s: %q\n", "local", d.Local)
			_, _ = fmt.Fprintf(f, "%s: %q\n", "message", d.Message)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(f, d.Message)
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", d.Message)
	}
}

type StackTrace string

func (s StackTrace) Annotate(m Modifier) {
	entries := strings.Split(string(debug.Stack()), "\n")
	m.AppendDetails(DebugInfo{entries, string(s)})
}
