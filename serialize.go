package errors

import (
	"encoding/json"
	"io"
)

var typeProvider = map[string]func() Any{
	TypeUrlDebugInfo:           func() Any { return new(DebugInfo) },
	TypeUrlResourceInfo:        func() Any { return new(ResourceInfo) },
	TypeUrlBadRequest:          func() Any { return new(BadRequest) },
	TypeUrlPreconditionFailure: func() Any { return new(PreconditionFailure) },
	TypeUrlErrorInfo:           func() Any { return new(ErrorInfo) },
	TypeUrlQuotaFailure:        func() Any { return new(QuotaFailure) },
	TypeUrlRequestInfo:         func() Any { return new(RequestInfo) },
	TypeUrlHelp:                func() Any { return new(Help) },
	TypeUrlLocalizedMessage:    func() Any { return new(LocalizedMessage) },
}

func RegisterType(typeUrl string, provider func() Any) {
	if provider == nil {
		delete(typeProvider, typeUrl)
	} else {
		typeProvider[typeUrl] = provider
	}
}

type encoded struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Status  string          `json:"status,omitempty"`
	Details json.RawMessage `json:"details,omitempty"`
}

func Encode(err error) (data []byte, _err error) {
	var o struct {
		Code    int    `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
		Status  string `json:"status,omitempty"`
		Details []Any  `json:"details,omitempty"`
	}
	if a, ok := Flatten(err).(*annotated); ok {
		o.Code = a.code.HttpCode()
		o.Status = a.code.String()
		o.Message = a.message
		o.Details = a.details
	}
	return json.Marshal(o)
}

func Decode(body io.Reader) (err error) {
	var msg encoded
	if err = json.NewDecoder(body).Decode(&msg); err != nil {
		return
	}

	details, err := decodeDetails(msg.Details)
	if err != nil {
		return
	}

	code := StrToCode(msg.Status)

	return &annotated{
		cause:   code,
		code:    code,
		message: msg.Message,
		details: details,
	}
}

func decodeDetails(raw []byte) (details []Any, err error) {
	var wrappers []json.RawMessage
	if err = json.Unmarshal(raw, &wrappers); err != nil {
		return
	}

	details = make([]Any, len(wrappers))
	for i, wrapper := range wrappers {
		if details[i], err = decodeDetail(wrapper); err != nil {
			return
		}
	}
	return
}

func decodeDetail(raw []byte) (detail Any, err error) {
	var w struct {
		Type string `json:"@type"`
	}
	if err = json.Unmarshal(raw, &w); err != nil {
		return
	}

	if provide, ok := typeProvider[w.Type]; ok {
		detail = provide()
	} else {
		detail = new(AnyDetail)
	}
	err = json.Unmarshal(raw, detail)
	return
}
