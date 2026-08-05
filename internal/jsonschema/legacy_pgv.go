package jsonschema

import (
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	pgv "github.com/envoyproxy/protoc-gen-validate/validate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// legacyPGVFieldRules reads the legacy protoc-gen-validate (validate.rules)
// extension and re-expresses the constraints that have JSON Schema equivalents
// as modern buf.validate rules, so the whole apply*ValidationRules mapping is
// reused for both extensions rather than duplicated per dialect.
//
// Only the bound and length rules are translated. Rules whose JSON Schema
// keyword does not depend on the dialect (const, in, formats, well-known-type
// rules) are deliberately left out: they carry over no numeric bound, and
// converting them would mean mirroring most of the PGV surface here.
func legacyPGVFieldRules(fd protoreflect.FieldDescriptor) *validate.FieldRules {
	opts := fd.Options()
	if opts == nil || !proto.HasExtension(opts, pgv.E_Rules) {
		return nil
	}
	rules, ok := proto.GetExtension(opts, pgv.E_Rules).(*pgv.FieldRules)
	if !ok || rules == nil {
		return nil
	}

	if fd.IsList() || fd.IsMap() {
		return legacyPGVCollectionRules(fd, rules)
	}

	return legacyPGVKindRules(fd.Kind(), rules)
}

// legacyPGVKindRules translates the rules for a single value of the given kind.
func legacyPGVKindRules(kind protoreflect.Kind, rules *pgv.FieldRules) *validate.FieldRules {
	switch kind {
	case protoreflect.StringKind:
		return legacyPGVStringRules(rules.GetString_())
	case protoreflect.BytesKind:
		return legacyPGVBytesRules(rules.GetBytes())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return legacyPGVInt32Rules(rules.GetInt32())
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return legacyPGVInt64Rules(rules.GetInt64())
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return legacyPGVUint32Rules(rules.GetUint32())
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return legacyPGVUint64Rules(rules.GetUint64())
	case protoreflect.FloatKind:
		return legacyPGVFloatRules(rules.GetFloat())
	case protoreflect.DoubleKind:
		return legacyPGVDoubleRules(rules.GetDouble())
	default:
		return nil
	}
}

// legacyPGVCollectionRules translates the list/map size rules and recurses into
// the per-element rules, mirroring how applyRepeatedFieldRules and
// applyMapFieldRules split rules by the level they constrain.
func legacyPGVCollectionRules(fd protoreflect.FieldDescriptor, rules *pgv.FieldRules) *validate.FieldRules {
	if fd.IsMap() {
		m := rules.GetMap()
		if m == nil {
			return nil
		}
		out := &validate.MapRules{}
		if m.MinPairs != nil {
			out.MinPairs = proto.Uint64(m.GetMinPairs())
		}
		if m.MaxPairs != nil {
			out.MaxPairs = proto.Uint64(m.GetMaxPairs())
		}
		if keys := m.GetKeys(); keys != nil {
			out.Keys = legacyPGVKindRules(fd.MapKey().Kind(), keys)
		}
		if values := m.GetValues(); values != nil {
			out.Values = legacyPGVKindRules(fd.MapValue().Kind(), values)
		}
		return &validate.FieldRules{Type: &validate.FieldRules_Map{Map: out}}
	}

	r := rules.GetRepeated()
	if r == nil {
		return nil
	}
	out := &validate.RepeatedRules{}
	if r.MinItems != nil {
		out.MinItems = proto.Uint64(r.GetMinItems())
	}
	if r.MaxItems != nil {
		out.MaxItems = proto.Uint64(r.GetMaxItems())
	}
	if r.Unique != nil {
		out.Unique = proto.Bool(r.GetUnique())
	}
	if items := r.GetItems(); items != nil {
		out.Items = legacyPGVKindRules(fd.Kind(), items)
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Repeated{Repeated: out}}
}

func legacyPGVStringRules(r *pgv.StringRules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.StringRules{}
	if r.MinLen != nil {
		out.MinLen = proto.Uint64(r.GetMinLen())
	}
	if r.MaxLen != nil {
		out.MaxLen = proto.Uint64(r.GetMaxLen())
	}
	// PGV len is an exact character count, which JSON Schema expresses as an
	// equal minLength/maxLength pair.
	if r.Len != nil {
		out.MinLen = proto.Uint64(r.GetLen())
		out.MaxLen = proto.Uint64(r.GetLen())
	}
	if r.Pattern != nil {
		out.Pattern = proto.String(r.GetPattern())
	}
	return &validate.FieldRules{Type: &validate.FieldRules_String_{String_: out}}
}

func legacyPGVBytesRules(r *pgv.BytesRules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.BytesRules{}
	if r.MinLen != nil {
		out.MinLen = proto.Uint64(r.GetMinLen())
	}
	if r.MaxLen != nil {
		out.MaxLen = proto.Uint64(r.GetMaxLen())
	}
	if r.Len != nil {
		out.MinLen = proto.Uint64(r.GetLen())
		out.MaxLen = proto.Uint64(r.GetLen())
	}
	if r.Pattern != nil {
		out.Pattern = proto.String(r.GetPattern())
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Bytes{Bytes: out}}
}

func legacyPGVInt32Rules(r *pgv.Int32Rules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.Int32Rules{}
	if r.Gte != nil {
		out.GreaterThan = &validate.Int32Rules_Gte{Gte: r.GetGte()}
	} else if r.Gt != nil {
		out.GreaterThan = &validate.Int32Rules_Gt{Gt: r.GetGt()}
	}
	if r.Lte != nil {
		out.LessThan = &validate.Int32Rules_Lte{Lte: r.GetLte()}
	} else if r.Lt != nil {
		out.LessThan = &validate.Int32Rules_Lt{Lt: r.GetLt()}
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Int32{Int32: out}}
}

func legacyPGVInt64Rules(r *pgv.Int64Rules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.Int64Rules{}
	if r.Gte != nil {
		out.GreaterThan = &validate.Int64Rules_Gte{Gte: r.GetGte()}
	} else if r.Gt != nil {
		out.GreaterThan = &validate.Int64Rules_Gt{Gt: r.GetGt()}
	}
	if r.Lte != nil {
		out.LessThan = &validate.Int64Rules_Lte{Lte: r.GetLte()}
	} else if r.Lt != nil {
		out.LessThan = &validate.Int64Rules_Lt{Lt: r.GetLt()}
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Int64{Int64: out}}
}

func legacyPGVUint32Rules(r *pgv.UInt32Rules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.UInt32Rules{}
	if r.Gte != nil {
		out.GreaterThan = &validate.UInt32Rules_Gte{Gte: r.GetGte()}
	} else if r.Gt != nil {
		out.GreaterThan = &validate.UInt32Rules_Gt{Gt: r.GetGt()}
	}
	if r.Lte != nil {
		out.LessThan = &validate.UInt32Rules_Lte{Lte: r.GetLte()}
	} else if r.Lt != nil {
		out.LessThan = &validate.UInt32Rules_Lt{Lt: r.GetLt()}
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Uint32{Uint32: out}}
}

func legacyPGVUint64Rules(r *pgv.UInt64Rules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.UInt64Rules{}
	if r.Gte != nil {
		out.GreaterThan = &validate.UInt64Rules_Gte{Gte: r.GetGte()}
	} else if r.Gt != nil {
		out.GreaterThan = &validate.UInt64Rules_Gt{Gt: r.GetGt()}
	}
	if r.Lte != nil {
		out.LessThan = &validate.UInt64Rules_Lte{Lte: r.GetLte()}
	} else if r.Lt != nil {
		out.LessThan = &validate.UInt64Rules_Lt{Lt: r.GetLt()}
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Uint64{Uint64: out}}
}

func legacyPGVFloatRules(r *pgv.FloatRules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.FloatRules{}
	if r.Gte != nil {
		out.GreaterThan = &validate.FloatRules_Gte{Gte: r.GetGte()}
	} else if r.Gt != nil {
		out.GreaterThan = &validate.FloatRules_Gt{Gt: r.GetGt()}
	}
	if r.Lte != nil {
		out.LessThan = &validate.FloatRules_Lte{Lte: r.GetLte()}
	} else if r.Lt != nil {
		out.LessThan = &validate.FloatRules_Lt{Lt: r.GetLt()}
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Float{Float: out}}
}

func legacyPGVDoubleRules(r *pgv.DoubleRules) *validate.FieldRules {
	if r == nil {
		return nil
	}
	out := &validate.DoubleRules{}
	if r.Gte != nil {
		out.GreaterThan = &validate.DoubleRules_Gte{Gte: r.GetGte()}
	} else if r.Gt != nil {
		out.GreaterThan = &validate.DoubleRules_Gt{Gt: r.GetGt()}
	}
	if r.Lte != nil {
		out.LessThan = &validate.DoubleRules_Lte{Lte: r.GetLte()}
	} else if r.Lt != nil {
		out.LessThan = &validate.DoubleRules_Lt{Lt: r.GetLt()}
	}
	return &validate.FieldRules{Type: &validate.FieldRules_Double{Double: out}}
}
