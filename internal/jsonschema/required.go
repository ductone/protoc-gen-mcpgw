package jsonschema

import (
	"regexp"
	"slices"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	pgv "github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// requiredByValidate reports whether a field's validation rules provably reject
// the zero value, which is the only thing that makes the field genuinely
// mandatory in the generated JSON Schema.
//
// It honors both validation dialects:
//   - modern buf.validate (protovalidate), extension buf.validate.field
//   - legacy protoc-gen-validate (PGV), extension validate.rules
//
// Many request protos still carry only the legacy PGV rules. A generator that
// reads only buf.validate therefore emits an empty "required" array for them,
// even for fields whose rules clearly reject an absent/empty value.
//
// The check is deliberately conservative: a field becomes required only when a
// rule makes the zero/empty value invalid (e.g. string min_len>=1, a pattern
// that rejects "", repeated min_items>=1, message.required). Optional filters,
// pagination tokens (page_size gte:0, page_token min_bytes:0), and any field
// flagged ignore-empty are never marked required.
func requiredByValidate(fd protoreflect.FieldDescriptor) bool {
	return requiredByBufValidate(fd) || requiredByLegacyPGV(fd)
}

// requiredByBufValidate handles the modern buf.validate (protovalidate)
// extension. It reads the extension actually set on the field — the previous
// approach iterated extensions declared inside FieldOptions, which never matches
// a buf.validate rule and so detected nothing.
func requiredByBufValidate(fd protoreflect.FieldDescriptor) bool {
	opts := fd.Options()
	if opts == nil || !proto.HasExtension(opts, validate.E_Field) {
		return false
	}
	rules, ok := proto.GetExtension(opts, validate.E_Field).(*validate.FieldRules)
	if !ok || rules == nil {
		return false
	}

	// Any ignore mode tolerates the zero value, so the field is never required.
	if rules.GetIgnore() != validate.Ignore_IGNORE_UNSPECIFIED {
		return false
	}
	// Explicit required rule.
	if rules.GetRequired() {
		return true
	}

	if fd.IsList() && !fd.IsMap() {
		r := rules.GetRepeated()
		return r != nil && r.GetMinItems() >= 1
	}
	if fd.IsMap() {
		return false
	}

	switch fd.Kind() {
	case protoreflect.StringKind:
		return bufStringRequired(rules.GetString())
	case protoreflect.BytesKind:
		b := rules.GetBytes()
		return b != nil && (b.GetMinLen() >= 1 || b.GetLen() >= 1)
	case protoreflect.EnumKind:
		e := rules.GetEnum()
		if e == nil {
			return false
		}
		return e.GetConst() != 0 || inExcludesZeroInt32(e.GetIn())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		r := rules.GetInt32()
		return r != nil && (r.GetConst() != 0 || inExcludesZeroInt32(r.GetIn()))
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		r := rules.GetInt64()
		return r != nil && (r.GetConst() != 0 || inExcludesZeroInt64(r.GetIn()))
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		r := rules.GetUint32()
		return r != nil && r.GetConst() != 0
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		r := rules.GetUint64()
		return r != nil && r.GetConst() != 0
	default:
		return false
	}
}

// bufStringRequired reports whether modern string rules reject the empty value.
func bufStringRequired(s *validate.StringRules) bool {
	if s == nil {
		return false
	}
	if s.GetMinLen() >= 1 || s.GetMinBytes() >= 1 || s.GetLen() >= 1 || s.GetLenBytes() >= 1 {
		return true
	}
	if s.GetConst() != "" {
		return true
	}
	if in := s.GetIn(); len(in) > 0 && !slices.Contains(in, "") {
		return true
	}
	// Well-known string formats all reject the empty string.
	if s.GetEmail() || s.GetHostname() || s.GetUri() || s.GetUriRef() ||
		s.GetAddress() || s.GetUuid() || s.GetIp() || s.GetIpv4() || s.GetIpv6() {
		return true
	}
	return patternRejectsEmpty(s.GetPattern())
}

// requiredByLegacyPGV handles the legacy protoc-gen-validate (validate.rules)
// extension.
func requiredByLegacyPGV(fd protoreflect.FieldDescriptor) bool {
	opts := fd.Options()
	if opts == nil || !proto.HasExtension(opts, pgv.E_Rules) {
		return false
	}
	rules, ok := proto.GetExtension(opts, pgv.E_Rules).(*pgv.FieldRules)
	if !ok || rules == nil {
		return false
	}

	// Repeated (non-map): a min_items floor rejects the empty list.
	if fd.IsList() && !fd.IsMap() {
		r := rules.GetRepeated()
		return r != nil && !r.GetIgnoreEmpty() && r.GetMinItems() >= 1
	}
	if fd.IsMap() {
		return false
	}

	switch fd.Kind() {
	case protoreflect.StringKind:
		return pgvStringRequired(rules.GetString_())
	case protoreflect.BytesKind:
		b := rules.GetBytes()
		return b != nil && !b.GetIgnoreEmpty() && (b.GetMinLen() >= 1 || b.GetLen() >= 1)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		m := rules.GetMessage()
		return m != nil && m.GetRequired()
	case protoreflect.BoolKind:
		b := rules.GetBool()
		return b != nil && b.Const != nil && b.GetConst()
	case protoreflect.EnumKind:
		e := rules.GetEnum()
		if e == nil {
			return false
		}
		if e.Const != nil && e.GetConst() != 0 {
			return true
		}
		return inExcludesZeroInt32(e.GetIn())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		r := rules.GetInt32()
		if r == nil || r.GetIgnoreEmpty() {
			return false
		}
		if r.Const != nil && r.GetConst() != 0 {
			return true
		}
		return inExcludesZeroInt32(r.GetIn())
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		r := rules.GetInt64()
		if r == nil || r.GetIgnoreEmpty() {
			return false
		}
		if r.Const != nil && r.GetConst() != 0 {
			return true
		}
		return inExcludesZeroInt64(r.GetIn())
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		r := rules.GetUint32()
		if r == nil || r.GetIgnoreEmpty() {
			return false
		}
		return r.Const != nil && r.GetConst() != 0
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		r := rules.GetUint64()
		if r == nil || r.GetIgnoreEmpty() {
			return false
		}
		return r.Const != nil && r.GetConst() != 0
	default:
		return false
	}
}

// pgvStringRequired reports whether legacy string rules reject the empty value.
func pgvStringRequired(s *pgv.StringRules) bool {
	if s == nil || s.GetIgnoreEmpty() {
		return false
	}
	if s.GetMinLen() >= 1 || s.GetMinBytes() >= 1 || s.GetLen() >= 1 || s.GetLenBytes() >= 1 {
		return true
	}
	if s.Const != nil && s.GetConst() != "" {
		return true
	}
	if in := s.GetIn(); len(in) > 0 && !slices.Contains(in, "") {
		return true
	}
	// Well-known string formats all reject the empty string.
	if s.GetEmail() || s.GetHostname() || s.GetUri() || s.GetUriRef() ||
		s.GetAddress() || s.GetUuid() || s.GetIp() || s.GetIpv4() || s.GetIpv6() {
		return true
	}
	return patternRejectsEmpty(s.GetPattern())
}

// patternRejectsEmpty reports whether a regex pattern is set and provably fails
// to match the empty string. A pattern that can't be compiled stays
// conservative (treated as not rejecting empty).
func patternRejectsEmpty(pattern string) bool {
	if pattern == "" {
		return false
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return !re.MatchString("")
}

func inExcludesZeroInt32(in []int32) bool {
	return len(in) > 0 && !slices.Contains(in, 0)
}

func inExcludesZeroInt64(in []int64) bool {
	return len(in) > 0 && !slices.Contains(in, 0)
}
