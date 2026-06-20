package jsonschema_test

import (
	"testing"

	validationv1 "github.com/ductone/protoc-gen-mcpgw/example/validation/v1"
	jsonschema "github.com/ductone/protoc-gen-mcpgw/internal/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func requiredOf(t *testing.T, m protoreflect.ProtoMessage) []string {
	t.Helper()
	schema, err := jsonschema.GenerateJSONSchema(m.ProtoReflect().Descriptor())
	require.NoError(t, err)
	if schema["required"] == nil {
		return nil
	}
	got, ok := schema["required"].([]string)
	require.True(t, ok, "required should be []string, got %T", schema["required"])
	return got
}

// Modern buf.validate: a field is required purely because its string.min_len
// rule rejects the empty value — no mcpgw field option involved. This proves the
// buf.validate required-derivation path actually reads the extension set on the
// field.
func TestRequired_ModernBufValidate(t *testing.T) {
	got := requiredOf(t, (*validationv1.ModernValidateRequest)(nil))
	assert.Equal(t, []string{"name"}, got)
}

// Legacy protoc-gen-validate (validate.rules): exercises the rules that real
// request protos rely on. Only the genuinely-mandatory fields appear.
func TestRequired_LegacyPGV(t *testing.T) {
	got := requiredOf(t, (*validationv1.LegacyPGVRequest)(nil))

	// id: pattern rejects empty; tags: repeated min_items>=1.
	assert.ElementsMatch(t, []string{"id", "tags"}, got)

	// Optional fields must never be required:
	//   query     -> ignore_empty
	//   pageSize  -> int32 gte:0 (zero is valid)
	//   pageToken -> bytes min_len:0
	assert.NotContains(t, got, "query")
	assert.NotContains(t, got, "pageSize")
	assert.NotContains(t, got, "pageToken")
}
