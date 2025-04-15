package v1

import "encoding/json"

type DecoderInput interface {
	Method() string
	Arguments() map[string]any
	// If available, the raw arguments as a json.RawMessage
	// Otherwise, the arguments are already unmarshaled into the Arguments map
	RawArguments() json.RawMessage
}
