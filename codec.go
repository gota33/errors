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

func Register(typeUrl string, provider func() Any) {
	if provider == nil {
		delete(typeProvider, typeUrl)
	} else {
		typeProvider[typeUrl] = provider
	}
}

type message struct {
	Error messageBody `json:"error,omitempty"`
}

type messageBody struct {
	Code    int        `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
	Status  StatusName `json:"status,omitempty"`
	Details []Any      `json:"details,omitempty"`
}

func Encode(w io.Writer, in error) error {
	var body messageBody
	if a, ok := Flatten(in).(*annotated); ok {
		body.Code = a.code.Http()
		body.Status = a.code.Name()
		body.Message = a.message
		body.Details = a.details
	}
	return json.NewEncoder(w).Encode(message{body})
}

type encodedMessage struct {
	Error encodedBody `json:"error,omitempty"`
}

type encodedBody struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Status  StatusName      `json:"status,omitempty"`
	Details json.RawMessage `json:"details,omitempty"`
}

func Decode(r io.Reader) (err error) {
	var msg encodedMessage
	if err = json.NewDecoder(r).Decode(&msg); err != nil {
		return
	}

	details, err := decodeDetails(msg.Error.Details)
	if err != nil {
		return
	}

	cause := msg.Error
	code := cause.Status.StatusCode()

	return &annotated{
		cause:   code,
		code:    code,
		message: cause.Message,
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
