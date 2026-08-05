package jsonschema

// TODO: Implement JSON Schema Generation from Protobuf Messages
//
// PLAN:
// 1. Define Function Signature: Create a function, likely `GenerateJSONSchema(md protoreflect.MessageDescriptor) (map[string]any, error)`,
//    that takes a message descriptor and returns a map representing its JSON Schema or an error.
// 2. Initialize Schema: Start with a basic JSON Schema structure (map[string]any):
//    - "$schema": "https://json-schema.org/draft/2020-12/schema"
//    - "title": Message name (from md.Name())
//    - "type": "object"
//    - "properties": Empty map[string]any{}
//    - "required": Empty []string{}
//    - "additionalProperties": false
// 3. Iterate Fields: Loop through `md.Fields()`.
// 4. Process Each Field (fd protoreflect.FieldDescriptor):
//    a. Determine JSON Schema Type: Map `fd.Kind()` and `fd.IsList()`, `fd.IsMap()` to JSON Schema types ("string", "integer", "number", "boolean", "array", "object").
//    b. Handle Nested Types: For `MessageKind`, recursively call `GenerateJSONSchema` for the field's message type (`fd.Message()`). Consider using "$ref" for complex/recursive structures if needed, otherwise embed. Manage recursion depth. For `EnumKind`, generate a schema with type and `enum` values (preferring string names).
//    c. Extract Description: Use `proto.GetExtension` to retrieve `mcpgw.v1.FieldOptions` from `fd.Options()` and set the "description" field if the extension is present. Requires importing the generated mcpgw proto Go package.
//    d. Extract Validation Rules: Use `proto.GetExtension` to retrieve `buf.validate.FieldConstraints` from `fd.Options()`. Requires importing the generated buf validate proto Go package.
//    e. Map Validation Rules: Convert buf validate rules (e.g., `string.min_len`, `int64.gt`, `repeated.min_items`) to corresponding JSON Schema keywords (e.g., "minLength", "exclusiveMinimum", "minItems").
//    f. Determine Required: Add the field's JSON name (`fd.JSONName()`) to the root schema's "required" list if it's considered mandatory. This depends on proto3 presence rules (non-optional scalars are implicitly required) and potentially buf validate rules like (`min_len >= 1`, `required: true`).
//    g. Handle Well-Known Types (WKTs): Implement specific mappings for WKTs found via `fd.Message().FullName()` (e.g., "google.protobuf.Timestamp" -> `{type: "string", format: "date-time"}`).
//    h. Add to Properties: Add the generated schema for the field to the root schema's "properties" map, using `fd.JSONName()` as the key.
// 5. Return Result: Return the completed root schema map.
//
// CHECKLIST:
//   [x] Define `GenerateJSONSchema` function signature.
//   [x] Implement basic schema initialization ($schema, title, type, etc.).
//   [x] Implement field iteration using `md.Fields()`.
//   [x] Implement Protobuf Kind -> JSON Schema type mapping:
//       [x] StringKind -> "string"
//       [x] BoolKind -> "boolean"
//       [x] Integer Kinds -> "integer"
//       [x] Float/Double Kinds -> "number"
//       [x] BytesKind -> "string" (consider "contentEncoding": "base64")
//       [x] EnumKind -> "string" / "integer" + "enum" values
//       [x] MessageKind -> "object" (recursive call)
//   [x] Handle `fd.IsList()` -> "array" with "items".
//   [x] Handle `fd.IsMap()` -> "object" with "additionalProperties".
//   [x] Implement WKT handling (Timestamp, Duration, etc.).
//   [x] Implement recursive call for nested messages (handle cycles/depth).
//   [x] Import mcpgw generated code and implement `proto.GetExtension` for description.
//   [x] Import buf validate generated code and implement `proto.GetExtension` for constraints.
//   [x] Implement mapping for buf validate rules to JSON Schema keywords:
//       [x] String rules (min/max_len, pattern, etc.)
//       [x] Numeric rules (min/max, gt/gte, lt/lte, etc.)
//       [x] Repeated rules (min/max_items, unique)
//       [x] Map rules (min/max_pairs)
//       [x] Required rule (`required: true`)
//       [x] WKT rules (duration/timestamp constraints)
//   [x] Implement logic to determine required fields (proto3 presence + validation rules).
//   [x] Use `fd.JSONName()` for property keys.
//   [x] Add error handling (e.g., for extension retrieval, recursion).
//   [ ] Write unit tests covering various types, validations, WKTs, nesting, and edge cases.

import (
	"fmt"
	"regexp"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
)

// Well-known type full names
const (
	// Wrapper types
	stringValueFullName = "google.protobuf.StringValue"
	boolValueFullName   = "google.protobuf.BoolValue"
	int32ValueFullName  = "google.protobuf.Int32Value"
	uint32ValueFullName = "google.protobuf.UInt32Value"
	int64ValueFullName  = "google.protobuf.Int64Value"
	uint64ValueFullName = "google.protobuf.UInt64Value"
	floatValueFullName  = "google.protobuf.FloatValue"
	doubleValueFullName = "google.protobuf.DoubleValue"
	bytesValueFullName  = "google.protobuf.BytesValue"

	// Other well-known types
	timestampFullName = "google.protobuf.Timestamp"
	durationFullName  = "google.protobuf.Duration"
	emptyFullName     = "google.protobuf.Empty"
	fieldMaskFullName = "google.protobuf.FieldMask"
	structFullName    = "google.protobuf.Struct"
	valueFullName     = "google.protobuf.Value"
	listValueFullName = "google.protobuf.ListValue"
	anyFullName       = "google.protobuf.Any"
)

// GenerateJSONSchema generates a JSON Schema (draft-2020-12) for a Protobuf message.
// The schema is returned as a map[string]any that can be marshaled to JSON.
// All schemas are fully inlined (no $ref or $defs).
func GenerateJSONSchema(md protoreflect.MessageDescriptor) (map[string]any, error) {
	if md == nil {
		return nil, fmt.Errorf("message descriptor cannot be nil")
	}

	// Track visited message types to prevent infinite recursion
	visited := make(map[protoreflect.FullName]bool)

	// Determine message-level options
	allowNulls := messageAllowsNulls(md)

	return schemaForMessage(md, visited, allowNulls)
}

// messageAllowsNulls reads mcpgw.v1.message options and returns true if allow_nulls is set.
func messageAllowsNulls(md protoreflect.MessageDescriptor) bool {
	if md == nil {
		return false
	}
	opts := md.Options()
	if opts == nil {
		return false
	}
	if !proto.HasExtension(opts, mcpgw_v1.E_Message) {
		return false
	}
	if extVal := proto.GetExtension(opts, mcpgw_v1.E_Message); extVal != nil {
		if extMsg, ok := extVal.(*mcpgw_v1.MessageOptions); ok && extMsg != nil {
			// Reflectively check presence to respect explicit false
			m := extMsg.ProtoReflect()
			if f := m.Descriptor().Fields().ByName("allow_nulls"); f != nil {
				if m.Has(f) {
					return m.Get(f).Bool()
				}
			}
			// If not explicitly set, default is false
		}
	}
	return false
}

// schemaForMessage generates a JSON Schema for a message type, tracking visited messages to prevent recursion
func schemaForMessage(md protoreflect.MessageDescriptor, visited map[protoreflect.FullName]bool, allowNulls bool) (map[string]any, error) {
	// Check for recursion
	if visited[md.FullName()] {
		// For recursive types, return a generic object schema
		return map[string]any{
			"type": "object",
			// No further details to prevent infinite recursion
		}, nil
	}

	// Mark this message as visited
	visited[md.FullName()] = true
	defer func() {
		// Remove from visited when done with this branch
		delete(visited, md.FullName())
	}()

	// Initialize the schema
	schema := map[string]any{
		"$schema":              "https://json-schema.org/draft/2020-12/schema",
		"title":                string(md.Name()),
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
	}

	// Fields that are required
	var requiredFields []string

	// Process all fields
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)

		// Generate schema for this field
		fieldSchema, err := schemaForField(fd, visited, allowNulls)
		if err != nil {
			return nil, fmt.Errorf("error processing field %s: %w", fd.Name(), err)
		}

		// Add field schema to properties
		properties := schema["properties"].(map[string]any)
		properties[string(fd.JSONName())] = fieldSchema

		// Check if field is required
		if isFieldRequired(fd) {
			requiredFields = append(requiredFields, string(fd.JSONName()))
		}
	}

	// Add required fields if any
	if len(requiredFields) > 0 {
		schema["required"] = requiredFields
	}

	return schema, nil
}

// schemaForField generates a JSON Schema for a single field
func schemaForField(fd protoreflect.FieldDescriptor, visited map[protoreflect.FullName]bool, allowNulls bool) (map[string]any, error) {
	// Handle repeated fields (non-map)
	if fd.IsList() && !fd.IsMap() {
		return schemaForRepeatedField(fd, visited, allowNulls)
	}

	// Handle map fields
	if fd.IsMap() {
		return schemaForMapField(fd, visited, allowNulls)
	}

	// Handle well-known types
	if fd.Kind() == protoreflect.MessageKind {
		wktSchema, isWKT := schemaForWellKnownType(fd.Message(), allowNulls)
		if isWKT {
			applyCustomFieldOptions(fd, wktSchema)
			applyValidationRules(fd, wktSchema)
			return wktSchema, nil
		}
	}

	// Handle regular fields based on kind
	fieldSchema, err := schemaForKind(fd, visited, allowNulls)
	if err != nil {
		return nil, err
	}

	// Apply custom field options (if any)
	applyCustomFieldOptions(fd, fieldSchema)

	// Apply validation rules from buf.validate (if any)
	applyValidationRules(fd, fieldSchema)

	// Handle nullability for proto3 optional fields or proto2 optional fields
	if allowNulls {
		applyNullability(fd, fieldSchema)
	}

	return fieldSchema, nil
}

// schemaForKind generates a schema based on the field's kind
func schemaForKind(fd protoreflect.FieldDescriptor, visited map[protoreflect.FullName]bool, allowNulls bool) (map[string]any, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return map[string]any{"type": "boolean"}, nil

	case protoreflect.StringKind:
		return map[string]any{"type": "string"}, nil

	case protoreflect.BytesKind:
		return map[string]any{
			"type":            "string",
			"contentEncoding": "base64",
		}, nil

	// Integer types
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return map[string]any{"type": "integer"}, nil

	// 64-bit integers - treat as strings in JSON schema per Proto JSON spec
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		// In Proto JSON, 64-bit integers are represented as strings by default
		return map[string]any{
			"type": "string",
			// Allow integers as well since parsers accept both
			// Use pattern to match valid integer format
			"pattern": "^-?[0-9]+$",
		}, nil

	// Floating point types
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		schema := map[string]any{"type": "number"}
		// Optionally allow "NaN", "Infinity", "-Infinity" as strings
		// which are valid in Proto JSON
		schema["oneOf"] = []map[string]any{
			{"type": "number"},
			{"type": "string", "enum": []string{"NaN", "Infinity", "-Infinity"}},
		}
		return schema, nil

	case protoreflect.EnumKind:
		return schemaForEnum(fd.Enum()), nil

	case protoreflect.MessageKind:
		// Generate schema for nested message
		msgSchema, err := schemaForMessage(fd.Message(), visited, allowNulls)
		if err != nil {
			return nil, fmt.Errorf("error processing nested message %s: %w", fd.Message().Name(), err)
		}
		return msgSchema, nil

	default:
		return nil, fmt.Errorf("unsupported field kind: %s", fd.Kind())
	}
}

// schemaForEnum generates a schema for an enum field
func schemaForEnum(ed protoreflect.EnumDescriptor) map[string]any {
	schema := map[string]any{
		"type": "string",
	}

	// Collect enum values
	values := ed.Values()
	enumValues := make([]string, 0, values.Len())

	for i := 0; i < values.Len(); i++ {
		enumValue := values.Get(i)
		enumValues = append(enumValues, string(enumValue.Name()))
	}

	if len(enumValues) > 0 {
		schema["enum"] = enumValues
	}

	return schema
}

// schemaForRepeatedField handles repeated fields (creates an array schema)
func schemaForRepeatedField(fd protoreflect.FieldDescriptor, visited map[protoreflect.FullName]bool, allowNulls bool) (map[string]any, error) {
	// Create item schema (schema for a single element of the array)
	itemSchema, err := schemaForKind(fd, visited, allowNulls)
	if err != nil {
		return nil, err
	}

	// Create array schema
	schema := map[string]any{
		"type":  "array",
		"items": itemSchema,
	}

	if allowNulls {
		// In Proto JSON, null can be accepted as an alias for empty array
		schema["type"] = []string{"array", "null"}
	}

	// The description documents the field, so it belongs on the array rather
	// than on a single element. Nullability is already handled above; repeated
	// fields have no presence, so applyNullability would be a no-op here.
	applyCustomFieldOptions(fd, schema)
	applyRepeatedFieldRules(fd, schema, itemSchema)

	return schema, nil
}

// schemaForMapField handles map fields
func schemaForMapField(fd protoreflect.FieldDescriptor, visited map[protoreflect.FullName]bool, allowNulls bool) (map[string]any, error) {
	// Get the value descriptor and generate its schema
	valueDesc := fd.MapValue()
	valueSchema, err := schemaForField(valueDesc, visited, allowNulls)
	if err != nil {
		return nil, err
	}

	// Create map schema (object with additionalProperties)
	schema := map[string]any{
		"type":                 "object",
		"additionalProperties": valueSchema,
	}

	applyCustomFieldOptions(fd, schema)
	applyMapFieldRules(fd, schema, valueSchema)

	return schema, nil
}

// schemaForWellKnownType returns specialized schema for well-known Protobuf types
func schemaForWellKnownType(md protoreflect.MessageDescriptor, allowNulls bool) (map[string]any, bool) {
	fullName := string(md.FullName())

	switch fullName {
	// Wrapper types
	case stringValueFullName:
		if allowNulls {
			return map[string]any{"type": []string{"string", "null"}}, true
		}
		return map[string]any{"type": "string"}, true
	case boolValueFullName:
		if allowNulls {
			return map[string]any{"type": []string{"boolean", "null"}}, true
		}
		return map[string]any{"type": "boolean"}, true
	case int32ValueFullName, uint32ValueFullName:
		if allowNulls {
			return map[string]any{"type": []string{"integer", "null"}}, true
		}
		return map[string]any{"type": "integer"}, true
	case int64ValueFullName, uint64ValueFullName:
		if allowNulls {
			return map[string]any{
				"type":    []string{"string", "null"},
				"pattern": "^-?[0-9]+$",
			}, true
		}
		return map[string]any{
			"type":    "string",
			"pattern": "^-?[0-9]+$",
		}, true
	case floatValueFullName, doubleValueFullName:
		if allowNulls {
			return map[string]any{
				"oneOf": []map[string]any{
					{"type": "number"},
					{"type": "string", "enum": []string{"NaN", "Infinity", "-Infinity"}},
					{"type": "null"},
				},
			}, true
		}
		return map[string]any{
			"oneOf": []map[string]any{
				{"type": "number"},
				{"type": "string", "enum": []string{"NaN", "Infinity", "-Infinity"}},
			},
		}, true
	case bytesValueFullName:
		if allowNulls {
			return map[string]any{
				"type":            []string{"string", "null"},
				"contentEncoding": "base64",
			}, true
		}
		return map[string]any{
			"type":            "string",
			"contentEncoding": "base64",
		}, true

	// Timestamp
	case timestampFullName:
		return map[string]any{
			"type":   "string",
			"format": "date-time",
		}, true

	// Duration
	case durationFullName:
		return map[string]any{
			"type":    "string",
			"pattern": "^-?[0-9]+(\\.[0-9]+)?s$",
		}, true

	// Empty
	case emptyFullName:
		return map[string]any{
			"type":                 "object",
			"additionalProperties": false,
		}, true

	// FieldMask
	case fieldMaskFullName:
		return map[string]any{
			"type": "string",
			// Optionally enforce format with a pattern for comma-separated paths
			"pattern": "^([a-zA-Z0-9_.]+)(,[a-zA-Z0-9_.]+)*$",
		}, true

	// Struct
	case structFullName:
		return map[string]any{
			"type": "object",
		}, true

	// Value (can be any JSON value)
	case valueFullName:
		// Allow any JSON value
		return map[string]any{
			"type": []string{"object", "array", "string", "number", "boolean", "null"},
		}, true

	// ListValue
	case listValueFullName:
		return map[string]any{
			"type": "array",
			"items": map[string]any{
				// Allow any JSON value for items
				"type": []string{"object", "array", "string", "number", "boolean", "null"},
			},
		}, true

	// Any
	case anyFullName:
		return map[string]any{
			"type": "object",
			// Optionally require @type field
			"properties": map[string]any{
				"@type": map[string]any{
					"type": "string",
				},
			},
		}, true

	default:
		return nil, false
	}
}

// isFieldRequired determines if a field should be marked as required in the schema
func isFieldRequired(fd protoreflect.FieldDescriptor) bool {
	// Proto2 required fields
	if fd.Cardinality() == protoreflect.Required {
		return true
	}

	opts := fd.Options()
	if opts == nil {
		return false
	}

	// If the mcpgw field option explicitly sets required, honor it (true or
	// false) — an explicit author choice overrides anything inferred below.
	if proto.HasExtension(opts, mcpgw_v1.E_Field) {
		if extVal := proto.GetExtension(opts, mcpgw_v1.E_Field); extVal != nil {
			if extMsg, ok := extVal.(proto.Message); ok && extMsg != nil {
				m := extMsg.ProtoReflect()
				if reqField := m.Descriptor().Fields().ByName("required"); reqField != nil {
					if m.Has(reqField) {
						return m.Get(reqField).Bool()
					}
				}
			}
		}
	}

	// Infer requiredness from validation rules (modern buf.validate and legacy
	// protoc-gen-validate). A field is required only when its rules provably
	// reject the zero value.
	return requiredByValidate(fd)
}

// applyCustomFieldOptions applies mcpgw.v1.field options if present
func applyCustomFieldOptions(fd protoreflect.FieldDescriptor, schema map[string]any) {
	opts := fd.Options()
	if opts == nil {
		return
	}

	// Check if our extension is registered
	if !proto.HasExtension(opts, mcpgw_v1.E_Field) {
		return
	}

	// Extract our field options
	fieldOpts, ok := proto.GetExtension(opts, mcpgw_v1.E_Field).(*mcpgw_v1.FieldOptions)
	if !ok || fieldOpts == nil {
		return
	}

	// Add description from our field options
	if fieldOpts.GetDescription() != "" {
		schema["description"] = fieldOpts.GetDescription()
	}
}

// fieldRules returns the validation rules set on a field, or nil. Both
// validation dialects are read: modern buf.validate (buf.validate.field) and
// legacy protoc-gen-validate (validate.rules), the latter translated to the
// modern shape so callers map one type. A schema generated for a PGV-annotated
// proto would otherwise carry no bounds at all, leaving an MCP client to
// discover a limit by tripping a runtime InvalidArgument.
func fieldRules(fd protoreflect.FieldDescriptor) *validate.FieldRules {
	opts := fd.Options()
	if opts == nil {
		return nil
	}

	// Check if our extension is registered
	if proto.HasExtension(opts, validate.E_Field) {
		if rules, ok := proto.GetExtension(opts, validate.E_Field).(*validate.FieldRules); ok {
			return rules
		}
		return nil
	}

	return legacyPGVFieldRules(fd)
}

// applyValidationRules applies buf.validate rules for a singular field. The
// required rule is not handled here; isFieldRequired drives the parent's
// "required" array.
func applyValidationRules(fd protoreflect.FieldDescriptor, schema map[string]any) {
	rules := fieldRules(fd)
	if rules == nil {
		return
	}

	applyKindValidationRules(fd.Kind(), rules, schema)
}

// applyKindValidationRules maps the type-specific buf.validate rules for a
// single value of the given kind onto schema.
func applyKindValidationRules(kind protoreflect.Kind, rules *validate.FieldRules, schema map[string]any) {
	switch kind {
	case protoreflect.StringKind:
		applyStringValidationRules(rules, schema)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		applyInt32ValidationRules(rules, schema)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		applyInt64ValidationRules(rules, schema)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		applyUint32ValidationRules(rules, schema)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		applyUint64ValidationRules(rules, schema)
	case protoreflect.FloatKind:
		applyFloatValidationRules(rules, schema)
	case protoreflect.DoubleKind:
		applyDoubleValidationRules(rules, schema)
	case protoreflect.BytesKind:
		applyBytesValidationRules(rules, schema)
	case protoreflect.BoolKind:
		applyBoolValidationRules(rules, schema)
	case protoreflect.EnumKind:
		applyEnumValidationRules(rules, schema)
	}
}

// applyRepeatedFieldRules splits a repeated field's buf.validate rules by the
// level each one constrains: min_items, max_items and unique describe the list
// as a whole, while repeated.items describes a single element and so must land
// on the item schema. Applying an element rule such as min_len to the array
// would produce a constraint that is silently never satisfiable.
func applyRepeatedFieldRules(fd protoreflect.FieldDescriptor, arraySchema map[string]any, itemSchema map[string]any) {
	rules := fieldRules(fd)
	if rules == nil {
		return
	}

	repeatedRules := rules.GetRepeated()
	if repeatedRules == nil {
		return
	}

	applyRepeatedValidationRules(rules, arraySchema)

	if items := repeatedRules.GetItems(); items != nil {
		applyKindValidationRules(fd.Kind(), items, itemSchema)
	}
}

// applyMapFieldRules splits a map field's buf.validate rules by level, the same
// way applyRepeatedFieldRules does: min_pairs and max_pairs describe the object,
// while map.keys and map.values describe a single key or value.
func applyMapFieldRules(fd protoreflect.FieldDescriptor, objectSchema map[string]any, valueSchema map[string]any) {
	rules := fieldRules(fd)
	if rules == nil {
		return
	}

	mapRules := rules.GetMap()
	if mapRules == nil {
		return
	}

	applyMapValidationRules(rules, objectSchema)

	// Proto JSON renders every map key as a JSON property name, so key rules are
	// only faithfully expressible for string keys; numeric bounds on a stringified
	// key would not mean what the rule says.
	if keys := mapRules.GetKeys(); keys != nil && fd.MapKey().Kind() == protoreflect.StringKind {
		keySchema := map[string]any{"type": "string"}
		applyKindValidationRules(protoreflect.StringKind, keys, keySchema)
		if len(keySchema) > 1 {
			objectSchema["propertyNames"] = keySchema
		}
	}

	if values := mapRules.GetValues(); values != nil {
		applyKindValidationRules(fd.MapValue().Kind(), values, valueSchema)
	}
}

// applyStringValidationRules applies string validation rules to a schema
func applyStringValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	stringRules := fieldOpts.GetString()
	if stringRules == nil {
		return
	}

	// Handle min_len
	if stringRules.MinLen != nil {
		schema["minLength"] = stringRules.GetMinLen()
	}

	// Handle max_len
	if stringRules.MaxLen != nil {
		schema["maxLength"] = stringRules.GetMaxLen()
	}

	// Handle pattern
	if stringRules.Pattern != nil {
		schema["pattern"] = stringRules.GetPattern()
	}

	// Handle email format
	if stringRules.GetEmail() {
		schema["format"] = "email"
	}

	// Handle hostname format
	if stringRules.GetHostname() {
		schema["format"] = "hostname"
	}

	// Handle uri format
	if stringRules.GetUri() {
		schema["format"] = "uri"
	}

	// Handle const
	if stringRules.Const != nil {
		schema["const"] = stringRules.GetConst()
	}

	// Handle in (enum)
	if len(stringRules.GetIn()) > 0 {
		schema["enum"] = stringRules.GetIn()
	}

	// Handle prefix
	if stringRules.Prefix != nil {
		// JSON Schema doesn't have prefix directly, we could use pattern
		prefix := stringRules.GetPrefix()
		schema["pattern"] = fmt.Sprintf("^%s.*", regexp.QuoteMeta(prefix))
	}

	// Handle suffix
	if stringRules.Suffix != nil {
		// JSON Schema doesn't have suffix directly, we could use pattern
		suffix := stringRules.GetSuffix()
		schema["pattern"] = fmt.Sprintf(".*%s$", regexp.QuoteMeta(suffix))
	}

	// Handle contains
	if stringRules.Contains != nil {
		// JSON Schema doesn't have contains for strings directly
		// We might need a more complex schema for this
	}
}

// applyInt32ValidationRules applies int32 validation rules to a schema
func applyInt32ValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	intRules := fieldOpts.GetInt32()
	if intRules == nil {
		return
	}

	// Handle const
	if intRules.Const != nil {
		schema["const"] = intRules.GetConst()
	}

	// The bound rules are oneofs, so presence is the set case and not a
	// non-zero value; testing the value would drop an explicit gte:0 floor.
	switch bound := intRules.GetGreaterThan().(type) {
	case *validate.Int32Rules_Gte:
		schema["minimum"] = bound.Gte
	case *validate.Int32Rules_Gt:
		schema["exclusiveMinimum"] = bound.Gt
	}

	switch bound := intRules.GetLessThan().(type) {
	case *validate.Int32Rules_Lte:
		schema["maximum"] = bound.Lte
	case *validate.Int32Rules_Lt:
		schema["exclusiveMaximum"] = bound.Lt
	}

	// Handle in (enum)
	if len(intRules.GetIn()) > 0 {
		schema["enum"] = intRules.GetIn()
	}
}

// applyInt64ValidationRules applies int64 validation rules to a schema
func applyInt64ValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	intRules := fieldOpts.GetInt64()
	if intRules == nil {
		return
	}

	// Handle const
	if intRules.Const != nil {
		schema["const"] = intRules.GetConst()
	}

	// The bound rules are oneofs, so presence is the set case and not a
	// non-zero value; testing the value would drop an explicit gte:0 floor.
	switch bound := intRules.GetGreaterThan().(type) {
	case *validate.Int64Rules_Gte:
		schema["minimum"] = bound.Gte
	case *validate.Int64Rules_Gt:
		schema["exclusiveMinimum"] = bound.Gt
	}

	switch bound := intRules.GetLessThan().(type) {
	case *validate.Int64Rules_Lte:
		schema["maximum"] = bound.Lte
	case *validate.Int64Rules_Lt:
		schema["exclusiveMaximum"] = bound.Lt
	}

	// Handle in (enum)
	if len(intRules.GetIn()) > 0 {
		schema["enum"] = intRules.GetIn()
	}
}

// applyUint32ValidationRules applies uint32 validation rules to a schema
func applyUint32ValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	uintRules := fieldOpts.GetUint32()
	if uintRules == nil {
		return
	}

	// Handle const
	if uintRules.Const != nil {
		schema["const"] = uintRules.GetConst()
	}

	// The bound rules are oneofs, so presence is the set case and not a
	// non-zero value; testing the value would drop an explicit gte:0 floor.
	switch bound := uintRules.GetGreaterThan().(type) {
	case *validate.UInt32Rules_Gte:
		schema["minimum"] = bound.Gte
	case *validate.UInt32Rules_Gt:
		schema["exclusiveMinimum"] = bound.Gt
	}

	switch bound := uintRules.GetLessThan().(type) {
	case *validate.UInt32Rules_Lte:
		schema["maximum"] = bound.Lte
	case *validate.UInt32Rules_Lt:
		schema["exclusiveMaximum"] = bound.Lt
	}

	// Handle in (enum)
	if len(uintRules.GetIn()) > 0 {
		schema["enum"] = uintRules.GetIn()
	}
}

// applyUint64ValidationRules applies uint64 validation rules to a schema
func applyUint64ValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	uintRules := fieldOpts.GetUint64()
	if uintRules == nil {
		return
	}

	// Handle const
	if uintRules.Const != nil {
		schema["const"] = uintRules.GetConst()
	}

	// The bound rules are oneofs, so presence is the set case and not a
	// non-zero value; testing the value would drop an explicit gte:0 floor.
	switch bound := uintRules.GetGreaterThan().(type) {
	case *validate.UInt64Rules_Gte:
		schema["minimum"] = bound.Gte
	case *validate.UInt64Rules_Gt:
		schema["exclusiveMinimum"] = bound.Gt
	}

	switch bound := uintRules.GetLessThan().(type) {
	case *validate.UInt64Rules_Lte:
		schema["maximum"] = bound.Lte
	case *validate.UInt64Rules_Lt:
		schema["exclusiveMaximum"] = bound.Lt
	}

	// Handle in (enum)
	if len(uintRules.GetIn()) > 0 {
		schema["enum"] = uintRules.GetIn()
	}
}

// applyFloatValidationRules applies float validation rules to a schema
func applyFloatValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	floatRules := fieldOpts.GetFloat()
	if floatRules == nil {
		return
	}

	// Handle const
	if floatRules.Const != nil {
		schema["const"] = floatRules.GetConst()
	}

	// The bound rules are oneofs, so presence is the set case and not a
	// non-zero value; testing the value would drop an explicit gte:0 floor.
	switch bound := floatRules.GetGreaterThan().(type) {
	case *validate.FloatRules_Gte:
		schema["minimum"] = bound.Gte
	case *validate.FloatRules_Gt:
		schema["exclusiveMinimum"] = bound.Gt
	}

	switch bound := floatRules.GetLessThan().(type) {
	case *validate.FloatRules_Lte:
		schema["maximum"] = bound.Lte
	case *validate.FloatRules_Lt:
		schema["exclusiveMaximum"] = bound.Lt
	}

	// Handle in (enum)
	if len(floatRules.GetIn()) > 0 {
		schema["enum"] = floatRules.GetIn()
	}
}

// applyDoubleValidationRules applies double validation rules to a schema
func applyDoubleValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	doubleRules := fieldOpts.GetDouble()
	if doubleRules == nil {
		return
	}

	// Handle const
	if doubleRules.Const != nil {
		schema["const"] = doubleRules.GetConst()
	}

	// The bound rules are oneofs, so presence is the set case and not a
	// non-zero value; testing the value would drop an explicit gte:0 floor.
	switch bound := doubleRules.GetGreaterThan().(type) {
	case *validate.DoubleRules_Gte:
		schema["minimum"] = bound.Gte
	case *validate.DoubleRules_Gt:
		schema["exclusiveMinimum"] = bound.Gt
	}

	switch bound := doubleRules.GetLessThan().(type) {
	case *validate.DoubleRules_Lte:
		schema["maximum"] = bound.Lte
	case *validate.DoubleRules_Lt:
		schema["exclusiveMaximum"] = bound.Lt
	}

	// Handle in (enum)
	if len(doubleRules.GetIn()) > 0 {
		schema["enum"] = doubleRules.GetIn()
	}
}

// applyBytesValidationRules applies bytes validation rules to a schema
func applyBytesValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	bytesRules := fieldOpts.GetBytes()
	if bytesRules == nil {
		return
	}

	// Handle const
	if bytesRules.Const != nil {
		schema["const"] = bytesRules.GetConst()
	}

	// Handle min_len
	if bytesRules.MinLen != nil {
		schema["minByteLength"] = bytesRules.GetMinLen()
	}

	// Handle max_len
	if bytesRules.MaxLen != nil {
		schema["maxByteLength"] = bytesRules.GetMaxLen()
	}

	// Handle pattern
	if bytesRules.Pattern != nil {
		// Not directly applicable to bytes in JSON Schema
	}
}

// applyBoolValidationRules applies bool validation rules to a schema
func applyBoolValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	boolRules := fieldOpts.GetBool()
	if boolRules == nil {
		return
	}

	// Handle const
	if boolRules.Const != nil {
		schema["const"] = boolRules.GetConst()
	}
}

// applyEnumValidationRules applies enum validation rules to a schema
func applyEnumValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	enumRules := fieldOpts.GetEnum()
	if enumRules == nil {
		return
	}

	// Handle const
	if enumRules.Const != nil {
		schema["const"] = enumRules.GetConst()
	}

	// Handle defined_only
	if enumRules.DefinedOnly != nil && enumRules.GetDefinedOnly() {
		// Already handled by the enum schema generator
	}

	// Handle in
	if len(enumRules.GetIn()) > 0 {
		schema["enum"] = enumRules.GetIn()
	}

	// Handle not_in
	if len(enumRules.GetNotIn()) > 0 {
		// Not directly mappable to JSON Schema without more complex schema
	}
}

// applyRepeatedValidationRules applies validation rules for repeated fields
func applyRepeatedValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	repeatedRules := fieldOpts.GetRepeated()
	if repeatedRules == nil {
		return
	}

	// Handle min_items
	if repeatedRules.MinItems != nil {
		schema["minItems"] = repeatedRules.GetMinItems()
	}

	// Handle max_items
	if repeatedRules.MaxItems != nil {
		schema["maxItems"] = repeatedRules.GetMaxItems()
	}

	// Handle unique
	if repeatedRules.Unique != nil && repeatedRules.GetUnique() {
		schema["uniqueItems"] = true
	}
}

// applyMapValidationRules applies validation rules for map fields
func applyMapValidationRules(fieldOpts *validate.FieldRules, schema map[string]any) {
	mapRules := fieldOpts.GetMap()
	if mapRules == nil {
		return
	}

	// Handle min_pairs
	if mapRules.MinPairs != nil {
		schema["minProperties"] = mapRules.GetMinPairs()
	}

	// Handle max_pairs
	if mapRules.MaxPairs != nil {
		schema["maxProperties"] = mapRules.GetMaxPairs()
	}
}

// applyNullability handles nullability for proto3 optional fields
func applyNullability(fd protoreflect.FieldDescriptor, schema map[string]any) {
	// Skip for required fields
	if fd.Cardinality() == protoreflect.Required {
		return
	}

	// For proto3 optional fields or proto2 optional fields, allow null
	if fd.HasPresence() {
		// If schema already has a type field
		if typeVal, ok := schema["type"]; ok {
			// If type is already an array of types
			if types, ok := typeVal.([]string); ok {
				// Check if "null" is already in the list
				hasNull := false
				for _, t := range types {
					if t == "null" {
						hasNull = true
						break
					}
				}

				if !hasNull {
					schema["type"] = append(types, "null")
				}
			} else if typeStr, ok := typeVal.(string); ok {
				// Single type string, convert to array with null
				schema["type"] = []string{typeStr, "null"}
			}
		}

		// If schema uses oneOf
		if oneOf, ok := schema["oneOf"].([]map[string]any); ok {
			// Add null as an option if not already present
			hasNull := false
			for _, option := range oneOf {
				if typeVal, ok := option["type"]; ok && typeVal == "null" {
					hasNull = true
					break
				}
			}

			if !hasNull {
				schema["oneOf"] = append(oneOf, map[string]any{"type": "null"})
			}
		}
	}
}
