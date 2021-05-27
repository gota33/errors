package errors

import (
	"encoding/json"
	"errors"
)

var (
	ErrNoEncoder = errors.New("encoder: inner encoder is required")
	ErrNoDecoder = errors.New("decoder: inner decoder is required")
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

type DetailFilter func(a Any) bool

func HideDebugInfo(a Any) bool {
	return a.TypeUrl() != TypeUrlDebugInfo
}

type encoder interface {
	Encode(v interface{}) error
}

type Encoder struct {
	Filters []DetailFilter
	encoder
}

func NewEncoder(enc encoder) *Encoder {
	return &Encoder{encoder: enc}
}

func (e *Encoder) Encode(in error) error {
	var body messageBody
	if a, ok := Flatten(in).(*annotated); ok {
		body.Code = a.code.Http()
		body.Status = a.code.Name()
		body.Message = a.message
		body.Details = e.filter(a.details)
	}

	if e.encoder == nil {
		return ErrNoEncoder
	}

	return e.encoder.Encode(message{body})
}

func (e *Encoder) filter(details []Any) (out []Any) {
Loop:
	for _, detail := range details {
		for _, filter := range e.Filters {
			if !filter(detail) {
				continue Loop
			}
		}
		out = append(out, detail)
	}
	return
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

type decoder interface {
	Decode(v interface{}) error
}

type Decoder struct {
	dec decoder
}

func NewDecoder(dec decoder) *Decoder {
	return &Decoder{dec: dec}
}

func (d Decoder) Decode() (err error) {
	if d.dec == nil {
		return ErrNoDecoder
	}

	var msg encodedMessage
	if err = d.dec.Decode(&msg); err != nil {
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
