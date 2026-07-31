package jsonschema_test

import (
	"testing"

	validationv1 "github.com/ductone/protoc-gen-mcpgw/example/validation/v1"
	jsonschema "github.com/ductone/protoc-gen-mcpgw/internal/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func propertiesOf(t *testing.T, m protoreflect.ProtoMessage) map[string]any {
	t.Helper()
	schema, err := jsonschema.GenerateJSONSchema(m.ProtoReflect().Descriptor())
	require.NoError(t, err)
	props, ok := schema["properties"].(map[string]any)
	require.True(t, ok, "properties should be a map, got %T", schema["properties"])
	return props
}

func propertyOf(t *testing.T, m protoreflect.ProtoMessage, name string) map[string]any {
	t.Helper()
	props := propertiesOf(t, m)
	prop, ok := props[name].(map[string]any)
	require.True(t, ok, "property %q should be a map, got %T", name, props[name])
	return prop
}

// A repeated scalar field carries its mcpgw description on the array schema.
func TestCollection_RepeatedScalarDescription(t *testing.T) {
	tags := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "tags")

	assert.Equal(t, "array", tags["type"])
	assert.Equal(t, "Tags applied to the resource", tags["description"])

	// The description documents the field, never a single element.
	items, ok := tags["items"].(map[string]any)
	require.True(t, ok, "items should be a map, got %T", tags["items"])
	assert.NotContains(t, items, "description")
}

// A repeated message field carries its description on the array schema while the
// item schema is left exactly as the nested message generator produced it.
func TestCollection_RepeatedMessageDescription(t *testing.T) {
	field := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "items")

	assert.Equal(t, "array", field["type"])
	assert.Equal(t, "Nested items", field["description"])

	items, ok := field["items"].(map[string]any)
	require.True(t, ok, "items should be a map, got %T", field["items"])
	assert.NotContains(t, items, "description")
	assert.Equal(t, "CollectionItem", items["title"])
	assert.Equal(t, "object", items["type"])

	// Fields nested inside the item type keep their own descriptions.
	itemProps, ok := items["properties"].(map[string]any)
	require.True(t, ok, "item properties should be a map, got %T", items["properties"])
	id, ok := itemProps["id"].(map[string]any)
	require.True(t, ok, "item id should be a map, got %T", itemProps["id"])
	assert.Equal(t, "The item id", id["description"])
}

// A map field carries its description on the object schema.
func TestCollection_MapDescription(t *testing.T) {
	labels := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "labels")

	assert.Equal(t, "object", labels["type"])
	assert.Equal(t, "Arbitrary labels", labels["description"])
}

// repeated.min_items/max_items/unique constrain the list, so they land on the
// array schema and not on items.
func TestCollection_ArrayLevelRules(t *testing.T) {
	tags := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "tags")

	assert.Equal(t, uint64(1), tags["minItems"])
	assert.Equal(t, uint64(10), tags["maxItems"])
	assert.Equal(t, true, tags["uniqueItems"])

	items, ok := tags["items"].(map[string]any)
	require.True(t, ok, "items should be a map, got %T", tags["items"])
	assert.NotContains(t, items, "minItems")
	assert.NotContains(t, items, "maxItems")
	assert.NotContains(t, items, "uniqueItems")
}

// repeated.items constrains each element, so it lands on items and never on the
// array — minLength on an array would be unsatisfiable.
func TestCollection_ItemLevelRules(t *testing.T) {
	tags := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "tags")

	items, ok := tags["items"].(map[string]any)
	require.True(t, ok, "items should be a map, got %T", tags["items"])
	assert.Equal(t, uint64(2), items["minLength"])
	assert.Equal(t, uint64(32), items["maxLength"])

	assert.NotContains(t, tags, "minLength")
	assert.NotContains(t, tags, "maxLength")
}

// map.min_pairs/max_pairs constrain the object; map.keys and map.values
// constrain a single key or value.
func TestCollection_MapLevelRules(t *testing.T) {
	labels := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "labels")

	assert.Equal(t, uint64(1), labels["minProperties"])
	assert.Equal(t, uint64(5), labels["maxProperties"])
	assert.NotContains(t, labels, "minLength")
	assert.NotContains(t, labels, "maxLength")

	keys, ok := labels["propertyNames"].(map[string]any)
	require.True(t, ok, "propertyNames should be a map, got %T", labels["propertyNames"])
	assert.Equal(t, "string", keys["type"])
	assert.Equal(t, uint64(1), keys["minLength"])

	values, ok := labels["additionalProperties"].(map[string]any)
	require.True(t, ok, "additionalProperties should be a map, got %T", labels["additionalProperties"])
	assert.Equal(t, uint64(64), values["maxLength"])
	assert.NotContains(t, values, "minLength")
}

// Regression guard: a singular scalar still takes the normal path unchanged.
func TestCollection_SingularScalarUnchanged(t *testing.T) {
	name := propertyOf(t, (*validationv1.CollectionOptionsRequest)(nil), "name")

	assert.Equal(t, "string", name["type"])
	assert.Equal(t, "The name of the resource", name["description"])
	assert.Equal(t, uint64(1), name["minLength"])
}

// allow_nulls widens a repeated field to ["array", "null"] exactly once —
// applyNullability must not append a second "null".
func TestCollection_NullabilityNotDoubleApplied(t *testing.T) {
	props := propertiesOf(t, (*validationv1.NullableCollectionRequest)(nil))

	tags, ok := props["tags"].(map[string]any)
	require.True(t, ok, "tags should be a map, got %T", props["tags"])
	assert.Equal(t, []string{"array", "null"}, tags["type"])
	assert.Equal(t, "Tags applied to the resource", tags["description"])

	// Maps are not widened, matching the generator's existing behaviour.
	labels, ok := props["labels"].(map[string]any)
	require.True(t, ok, "labels should be a map, got %T", props["labels"])
	assert.Equal(t, "object", labels["type"])
	assert.Equal(t, "Arbitrary labels", labels["description"])
}
