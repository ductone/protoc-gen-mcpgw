package v1

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func UnmarshalFromMap(args map[string]any, out proto.Message) error {
	// TODO(pquerna): future optimziation: avoid json marshalling
	//
	// We use this function from generated code, so we can later
	// optimize this to avoid json marshalling.
	jsonArgs, err := json.Marshal(args)
	if err != nil {
		return err
	}
	return protojson.Unmarshal(jsonArgs, out)
}
