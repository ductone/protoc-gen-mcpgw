package jsonschema_test

import (
	"testing"

	validationv1 "github.com/ductone/protoc-gen-mcpgw/example/validation/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Enum rules hold numbers; the schema must list value names, or it describes
// a string that can never match.
func TestEnumRules_InMapsNumbersToNames(t *testing.T) {
	prop := propertyOf(t, (*validationv1.EnumRulesRequest)(nil), "allowed")
	assert.Equal(t, "string", prop["type"])
	assert.Equal(t, []string{"SHADE_LIGHT", "SHADE_DARK"}, prop["enum"])
}

func TestEnumRules_NotInFiltersNames(t *testing.T) {
	prop := propertyOf(t, (*validationv1.EnumRulesRequest)(nil), "excluded")
	assert.Equal(t, []string{"SHADE_LIGHT", "SHADE_MEDIUM", "SHADE_DARK"}, prop["enum"])
}

func TestEnumRules_ConstMapsToName(t *testing.T) {
	prop := propertyOf(t, (*validationv1.EnumRulesRequest)(nil), "fixed")
	assert.Equal(t, "SHADE_MEDIUM", prop["const"])
}

func TestEnumRules_RepeatedItems(t *testing.T) {
	prop := propertyOf(t, (*validationv1.EnumRulesRequest)(nil), "allowedList")
	items, ok := prop["items"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, []string{"SHADE_LIGHT", "SHADE_MEDIUM"}, items["enum"])
}

func TestPatternRules_TranslatedToECMA(t *testing.T) {
	assert.Equal(t, `^a[^\n]b$`, propertyOf(t, (*validationv1.PatternRulesRequest)(nil), "dotted")["pattern"])
	assert.Equal(t, `^[Aa][Bb]$`, propertyOf(t, (*validationv1.PatternRulesRequest)(nil), "folded")["pattern"])
}

// A pattern with no exact ECMA-262 equivalent is omitted; the server still
// enforces it.
func TestPatternRules_UntranslatableOmitted(t *testing.T) {
	assert.NotContains(t, propertyOf(t, (*validationv1.PatternRulesRequest)(nil), "multiline"), "pattern")
}

// prefix and suffix are ANDed with pattern instead of replacing it.
func TestPatternRules_PrefixSuffixCombine(t *testing.T) {
	prop := propertyOf(t, (*validationv1.PatternRulesRequest)(nil), "prefixed")
	assert.Equal(t, "^[a-z]+$", prop["pattern"])
	assert.Equal(t, []any{
		map[string]any{"pattern": `^a\.b`},
		map[string]any{"pattern": "z$"},
	}, prop["allOf"])
}
