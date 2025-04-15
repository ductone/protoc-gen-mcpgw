package schema

import (
	"sync"

	"github.com/ductone/protoc-gen-mcpgw/internal/jsonschema"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var cache = sync.Map{}

func GenerateSchema(md protoreflect.MessageDescriptor) (map[string]any, error) {
	cacheKey := md.FullName()
	if cached, ok := cache.Load(cacheKey); ok {
		return cached.(map[string]any), nil
	}
	schema, err := jsonschema.GenerateJSONSchema(md)
	if err != nil {
		return nil, err
	}
	cache.Store(cacheKey, schema)
	return schema, nil
}

func MustGenerateSchema(md protoreflect.MessageDescriptor) map[string]any {
	schema, err := GenerateSchema(md)
	if err != nil {
		panic(err)
	}
	return schema
}
