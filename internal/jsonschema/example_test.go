package jsonschema_test

import (
	"encoding/json"
	"testing"

	v1 "github.com/ductone/protoc-gen-mcpgw/example/bookstore/v1"
	jsonschema "github.com/ductone/protoc-gen-mcpgw/internal/jsonschema"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenreRequestSchema(t *testing.T) {
	// Get the message descriptor for the CreateGenreRequest type
	md := (&v1.CreateGenreRequest{}).ProtoReflect().Descriptor()

	// Generate JSON schema for the message
	schema, err := jsonschema.GenerateJSONSchema(md)
	assert.NoError(t, err)

	// Basic schema checks
	assert.Equal(t, "CreateGenreRequest", schema["title"])
	assert.Equal(t, "object", schema["type"])
	assert.False(t, schema["additionalProperties"].(bool))

	// Check properties
	properties, ok := schema["properties"].(map[string]any)
	assert.True(t, ok, "Properties should be a map")

	// Check the name field
	nameProperty, ok := properties["name"].(map[string]any)
	assert.True(t, ok, "Name field should be present")

	// In the actual schema, type is an array of strings
	typeArray, ok := nameProperty["type"].([]string)
	assert.True(t, ok, "Type should be a string array")
	assert.Contains(t, typeArray, "string")
	assert.Contains(t, typeArray, "null")

	// Check validation constraints - in the implementation they're uint64
	assert.Equal(t, uint64(1), nameProperty["minLength"])
	assert.Equal(t, uint64(50), nameProperty["maxLength"])
	assert.Equal(t, "The name of the genre", nameProperty["description"])

	// Pretty print the schema for debugging
	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	t.Logf("Generated schema: %s", jsonBytes)
}

// TestGetBookRequestSchema tests schema generation for a message with multiple field types
func TestGetBookRequestSchema(t *testing.T) {
	md := (&v1.GetBookRequest{}).ProtoReflect().Descriptor()
	schema, err := jsonschema.GenerateJSONSchema(md)
	assert.NoError(t, err)

	// Basic schema checks
	assert.Equal(t, "GetBookRequest", schema["title"])
	assert.Equal(t, "object", schema["type"])

	// Check properties exists and has the expected fields
	properties, ok := schema["properties"].(map[string]any)
	assert.True(t, ok, "Properties should be a map")

	// Verify all field types are present
	assert.Contains(t, properties, "shelf")         // string
	assert.Contains(t, properties, "book")          // int64
	assert.Contains(t, properties, "includeAuthor") // bool
	assert.Contains(t, properties, "pageSize")      // int32
	assert.Contains(t, properties, "pageToken")     // string

	// Check field types
	shelfProp := properties["shelf"].(map[string]any)
	typeArray, ok := shelfProp["type"].([]string)
	assert.True(t, ok)
	assert.Contains(t, typeArray, "string")

	bookProp := properties["book"].(map[string]any)
	// Int64 type should include a pattern for string representation
	assert.Contains(t, bookProp, "pattern")
	typeArray, ok = bookProp["type"].([]string)
	assert.True(t, ok)
	assert.Contains(t, typeArray, "string") // int64 is represented as a string in JSON

	includeProp := properties["includeAuthor"].(map[string]any)
	typeArray, ok = includeProp["type"].([]string)
	assert.True(t, ok)
	assert.Contains(t, typeArray, "boolean")

	// Pretty print the schema for debugging
	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	t.Logf("Generated schema: %s", jsonBytes)
}

// TestCreateBookRequestSchema tests schema generation for a message with a nested message
func TestCreateBookRequestSchema(t *testing.T) {
	md := (&v1.CreateBookRequest{}).ProtoReflect().Descriptor()
	schema, err := jsonschema.GenerateJSONSchema(md)
	assert.NoError(t, err)

	// Basic schema checks
	assert.Equal(t, "CreateBookRequest", schema["title"])

	// Check properties exists
	properties, ok := schema["properties"].(map[string]any)
	assert.True(t, ok, "Properties should be a map")

	// Verify shelf field
	assert.Contains(t, properties, "shelf")

	// Verify book field (nested message)
	assert.Contains(t, properties, "book")
	bookProp := properties["book"].(map[string]any)

	// In the actual schema, type is an array of strings
	typeArray, ok := bookProp["type"].([]string)
	assert.True(t, ok, "Type should be a string array")
	assert.Contains(t, typeArray, "object")
	assert.Contains(t, typeArray, "null")

	// Check that the nested book has its own properties
	bookProperties, ok := bookProp["properties"].(map[string]any)
	assert.True(t, ok, "Book properties should be a map")
	assert.Contains(t, bookProperties, "id")
	assert.Contains(t, bookProperties, "author")
	assert.Contains(t, bookProperties, "title")
	assert.Contains(t, bookProperties, "quotes")

	// Pretty print the schema for debugging
	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	t.Logf("Generated schema: %s", jsonBytes)
}

// TestListGenresRequestSchema tests schema generation for an empty message
func TestListGenresRequestSchema(t *testing.T) {
	md := (&v1.ListGenresRequest{}).ProtoReflect().Descriptor()
	schema, err := jsonschema.GenerateJSONSchema(md)
	assert.NoError(t, err)

	// Basic schema checks
	assert.Equal(t, "ListGenresRequest", schema["title"])
	assert.Equal(t, "object", schema["type"])

	// Check properties exists and should be empty
	properties, ok := schema["properties"].(map[string]any)
	assert.True(t, ok, "Properties should be a map")
	assert.Empty(t, properties, "Properties should be empty for a message with no fields")

	// Make sure additionalProperties is false
	assert.False(t, schema["additionalProperties"].(bool))

	// Pretty print the schema for debugging
	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	t.Logf("Generated schema: %s", jsonBytes)
}

// TestRecursiveBookRequestSchema tests schema generation with potential recursive references
func TestRecursiveBookRequestSchema(t *testing.T) {
	md := (&v1.RecursiveBookRequest{}).ProtoReflect().Descriptor()
	schema, err := jsonschema.GenerateJSONSchema(md)
	assert.NoError(t, err)

	// Basic schema checks
	assert.Equal(t, "RecursiveBookRequest", schema["title"])

	// Check properties
	properties, ok := schema["properties"].(map[string]any)
	assert.True(t, ok, "Properties should be a map")

	// Verify bookId field
	assert.Contains(t, properties, "bookId")

	// Get one level of recursive response just to make sure it doesn't loop infinitely
	recursiveResponesMd := (&v1.RecursiveBookResponse{}).ProtoReflect().Descriptor()
	recursiveSchema, err := jsonschema.GenerateJSONSchema(recursiveResponesMd)
	assert.NoError(t, err)

	// Verify pages field which contains a recursive reference
	properties, ok = recursiveSchema["properties"].(map[string]any)
	assert.True(t, ok, "Properties should be a map")
	assert.Contains(t, properties, "page")
	assert.Contains(t, properties, "anotherProp")

	// Pretty print the schema for debugging
	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	t.Logf("Generated schema: %s", jsonBytes)
}

// TestGetAuthorResponseSchema tests schema generation for a message with a oneof field
func TestGetAuthorResponseSchema(t *testing.T) {
	md := (&v1.GetAuthorResponse{}).ProtoReflect().Descriptor()
	schema, err := jsonschema.GenerateJSONSchema(md)
	assert.NoError(t, err)

	// Basic schema checks
	assert.Equal(t, "GetAuthorResponse", schema["title"])

	// Check properties exist
	properties, ok := schema["properties"].(map[string]any)
	assert.True(t, ok, "Properties should be a map")

	// Verify author field
	assert.Contains(t, properties, "author")

	// Verify oneof fields (fiction and nonfiction) are present
	assert.Contains(t, properties, "fiction")
	assert.Contains(t, properties, "nonfiction")

	// Check the type of the fiction field
	fictionProp := properties["fiction"].(map[string]any)
	typeArray, ok := fictionProp["type"].([]string)
	assert.True(t, ok, "Type should be a string array")
	assert.Contains(t, typeArray, "boolean")
	assert.Contains(t, typeArray, "null") // Nullable since it's part of a oneof

	// Pretty print the schema for debugging
	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	t.Logf("Generated schema: %s", jsonBytes)
}
