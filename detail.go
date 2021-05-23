package errors

import (
	"encoding/json"
	"fmt"
	"io"
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

var detailProvider = map[string]func() Detail{
	TypeUrlDebugInfo: func() Detail { return new(DebugInfo) },
}

func RegisterDetailProvider(typeUrl string, provider func() Detail) {
	if provider == nil {
		delete(detailProvider, typeUrl)
	} else {
		detailProvider[typeUrl] = provider
	}
}

type Detail interface {
	TypeUrl() string
}

func decodeDetails(raw []byte) (details []Detail, err error) {
	var wrappers []json.RawMessage
	if err = json.Unmarshal(raw, &wrappers); err != nil {
		return
	}

	details = make([]Detail, len(wrappers))
	for i, wrapper := range wrappers {
		if details[i], err = decodeDetail(wrapper); err != nil {
			return
		}
	}
	return
}

func decodeDetail(raw []byte) (detail Detail, err error) {
	var w struct {
		Type string `json:"@type"`
	}
	if err = json.Unmarshal(raw, &w); err != nil {
		return
	}

	if provide, ok := detailProvider[w.Type]; ok {
		detail = provide()
	} else {
		detail = new(AnyDetail)
	}
	err = json.Unmarshal(raw, detail)
	return
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

func (d DebugInfo) TypeUrl() string {
	return TypeUrlDebugInfo
}

type debugInfo DebugInfo

func (d DebugInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		debugInfo
	}{d.TypeUrl(), debugInfo(d)})
}

func (d DebugInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "detail", d.Detail)
			_, _ = fmt.Fprintf(f, "\t%-6s:\n", "stack")
			for _, entry := range d.StackEntries {
				_, _ = fmt.Fprintf(f, "\t\t%s\n", entry)
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

func (d ResourceInfo) TypeUrl() string {
	return TypeUrlResourceInfo
}

type resourceInfo ResourceInfo

func (d ResourceInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		resourceInfo
	}{d.TypeUrl(), resourceInfo(d)})
}

func (d ResourceInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-13s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-13s: %q\n", "resource_type", d.ResourceType)
			_, _ = fmt.Fprintf(f, "\t%-13s: %q\n", "resource_name", d.ResourceName)
			_, _ = fmt.Fprintf(f, "\t%-13s: %q\n", "owner", d.Owner)
			_, _ = fmt.Fprintf(f, "\t%-13s: %q\n", "description", d.Description)
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
	Field       string `json:"field"`
	Description string `json:"description"`
}

type BadRequest struct {
	FieldViolations []FieldViolation `json:"fieldViolations"`
}

func (d BadRequest) TypeUrl() string {
	return TypeUrlBadRequest
}

type badRequest BadRequest

func (d BadRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		badRequest
	}{d.TypeUrl(), badRequest(d)})
}

func (d BadRequest) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-6s:\n", "field_violations")
			for _, v := range d.FieldViolations {
				_, _ = fmt.Fprintf(f, "\t\t%-8s: %s\n", v.Field, v.Description)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "field violations: %d", len(d.FieldViolations))
	}
}

type TypedViolation struct {
	Type        string `json:"type"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type PreconditionFailure struct {
	Violations []TypedViolation `json:"violations"`
}

func (d PreconditionFailure) TypeUrl() string {
	return TypeUrlPreconditionFailure
}

type preconditionFailure PreconditionFailure

func (d PreconditionFailure) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		preconditionFailure
	}{d.TypeUrl(), preconditionFailure(d)})
}

func (d PreconditionFailure) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-6s:\n", "violations")
			for _, v := range d.Violations {
				_, _ = fmt.Fprintf(f, "\t\t[%s] %-8s: %s\n", v.Type, v.Subject, v.Description)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "violations: %d", len(d.Violations))
	}
}

type ErrorInfo struct {
	Reason   string            `json:"reason"`
	Domain   string            `json:"domain"`
	Metadata map[string]string `json:"metadata"`
}

func (d ErrorInfo) TypeUrl() string {
	return TypeUrlErrorInfo
}

type errorInfo ErrorInfo

func (d ErrorInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		errorInfo
	}{d.TypeUrl(), errorInfo(d)})
}

func (d ErrorInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "reason", d.Reason)
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "domain", d.Domain)
			_, _ = fmt.Fprintf(f, "\t%-6s:\n", "metadata")
			for k, v := range d.Metadata {
				_, _ = fmt.Fprintf(f, "\t\t%s: %s\n", k, v)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "[%s] %s", d.Domain, d.Reason)
	}
}

type Violation struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type QuotaFailure struct {
	Violations []Violation `json:"violations"`
}

func (d QuotaFailure) TypeUrl() string {
	return TypeUrlQuotaFailure
}

type quotaFailure QuotaFailure

func (d QuotaFailure) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		quotaFailure
	}{d.TypeUrl(), quotaFailure(d)})
}

func (d QuotaFailure) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-6s:\n", "violations")
			for _, v := range d.Violations {
				_, _ = fmt.Fprintf(f, "\t\t%-8s: %s\n", v.Subject, v.Description)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "violations: %d", len(d.Violations))
	}
}

type RequestInfo struct {
	RequestId   string `json:"requestId"`
	ServingData string `json:"servingData"`
}

func (d RequestInfo) TypeUrl() string {
	return TypeUrlRequestInfo
}

type requestInfo RequestInfo

func (d RequestInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		requestInfo
	}{d.TypeUrl(), requestInfo(d)})
}

func (d RequestInfo) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "request_id", d.RequestId)
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "serving_data", d.ServingData)
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
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Help struct {
	Links []Link `json:"links"`
}

func (d Help) TypeUrl() string {
	return TypeUrlHelp
}

type help Help

func (d Help) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		help
	}{d.TypeUrl(), help(d)})
}

func (d Help) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%s:\n", "links")
			for _, v := range d.Links {
				_, _ = fmt.Fprintf(f, "\t\t%s: %s\n", v.Description, v.Url)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(f, "help(%d)", len(d.Links))
	}
}

type LocalizedMessage struct {
	Local   string `json:"local"`
	Message string `json:"message"`
}

func (d LocalizedMessage) TypeUrl() string {
	return TypeUrlLocalizedMessage
}

type localizedMessage LocalizedMessage

func (d LocalizedMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"@type"`
		localizedMessage
	}{d.TypeUrl(), localizedMessage(d)})
}

func (d LocalizedMessage) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "type", d.TypeUrl())
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "local", d.Local)
			_, _ = fmt.Fprintf(f, "\t%-6s: %q\n", "message", d.Message)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(f, d.Message)
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", d.Message)
	}
}
