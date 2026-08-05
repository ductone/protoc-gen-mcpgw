package jsonschema_test

import (
	"testing"

	validationv1 "github.com/ductone/protoc-gen-mcpgw/example/validation/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The constraint mapping read only buf.validate.field, so a proto annotated with
// legacy protoc-gen-validate generated a schema carrying no bounds at all. The
// limit survived only as prose in the description, which an MCP client cannot
// act on -- it learns the ceiling by tripping a runtime InvalidArgument.
func TestLegacyPGVBoundsReachSchema(t *testing.T) {
	msg := (*validationv1.LegacyPGVBoundsRequest)(nil)

	pageSize := propertyOf(t, msg, "pageSize")
	assert.EqualValues(t, 0, pageSize["minimum"], "int32 gte must map to minimum")
	assert.EqualValues(t, 100, pageSize["maximum"], "int32 lte must map to maximum")
	assert.NotEmpty(t, pageSize["description"], "the mcpgw description must survive alongside the bounds")

	query := propertyOf(t, msg, "query")
	assert.EqualValues(t, 1024, query["maxLength"], "string max_len must map to maxLength")

	exclusive := propertyOf(t, msg, "exclusiveBound")
	assert.EqualValues(t, 5, exclusive["exclusiveMinimum"], "int32 gt must map to exclusiveMinimum")
	assert.EqualValues(t, 50, exclusive["exclusiveMaximum"], "int32 lt must map to exclusiveMaximum")

	name := propertyOf(t, msg, "name")
	assert.EqualValues(t, 3, name["minLength"], "string min_len must map to minLength")
	assert.Equal(t, "^[a-z]+$", name["pattern"], "string pattern must carry over")

	// PGV len is an exact count; JSON Schema expresses it as an equal pair.
	code := propertyOf(t, msg, "code")
	assert.EqualValues(t, 8, code["minLength"])
	assert.EqualValues(t, 8, code["maxLength"])

	tags := propertyOf(t, msg, "tags")
	assert.EqualValues(t, 10, tags["maxItems"], "repeated max_items must map to maxItems")
	items, ok := tags["items"].(map[string]any)
	require.True(t, ok, "array schema has no items")
	assert.EqualValues(t, 32, items["maxLength"], "repeated.items rules must land on the item schema")
}

// A field with no rules must stay unconstrained.
func TestLegacyPGVLeavesUnannotatedFieldsAlone(t *testing.T) {
	plain := propertyOf(t, (*validationv1.LegacyPGVBoundsRequest)(nil), "plain")

	assert.NotContains(t, plain, "minimum")
	assert.NotContains(t, plain, "maximum")
	assert.NotContains(t, plain, "minLength")
	assert.NotContains(t, plain, "maxLength")
	assert.NotContains(t, plain, "pattern")
}

// Modern buf.validate must keep working unchanged now that a second dialect is
// read; the legacy path is a fallback, not a replacement.
func TestModernBufValidateStillMapsBounds(t *testing.T) {
	tags := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "tags")

	assert.EqualValues(t, 1, tags["minItems"])
	assert.EqualValues(t, 10, tags["maxItems"])
	items, ok := tags["items"].(map[string]any)
	require.True(t, ok, "array schema has no items")
	assert.EqualValues(t, 2, items["minLength"])
	assert.EqualValues(t, 32, items["maxLength"])
}
