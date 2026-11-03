package jsonschema

import (
	"fmt"
	"regexp/syntax"
	"strings"
	"unicode"
)

// JSON Schema "pattern" uses the ECMA-262 dialect, while protovalidate and
// protoc-gen-validate patterns are RE2. ecmaPattern rewrites an RE2 pattern into
// an ECMA-262 pattern that matches exactly the same strings, built from the
// parsed RE2 syntax tree rather than the source text so that flags, Perl and
// POSIX classes and RE2-only escapes are expanded into plain constructs. It
// reports false when no exact equivalent exists; the caller then omits the
// keyword, since the server still enforces the rule.
//
// The output uses only escapes that mean the same with and without the ECMA "u"
// flag, and match semantics are those of the "u" flag (code points, as in RE2).
func ecmaPattern(re2 string) (string, bool) {
	re, err := syntax.Parse(re2, syntax.Perl)
	if err != nil {
		return "", false
	}
	var b strings.Builder
	if !writeECMA(&b, re) {
		return "", false
	}
	return b.String(), true
}

// ecmaQuote returns an ECMA-262 pattern matching s literally.
func ecmaQuote(s string) (string, bool) {
	var b strings.Builder
	for _, r := range s {
		if !writeECMARune(&b, r, false) {
			return "", false
		}
	}
	return b.String(), true
}

func writeECMA(b *strings.Builder, re *syntax.Regexp) bool {
	switch re.Op {
	case syntax.OpNoMatch:
		b.WriteString("[]")
	case syntax.OpEmptyMatch:
		b.WriteString("(?:)")
	case syntax.OpLiteral:
		for _, r := range re.Rune {
			if re.Flags&syntax.FoldCase != 0 && hasFold(r) {
				if !writeECMAClass(b, foldRanges(r)) {
					return false
				}
				continue
			}
			if !writeECMARune(b, r, false) {
				return false
			}
		}
	case syntax.OpCharClass:
		return writeECMAClass(b, re.Rune)
	case syntax.OpAnyCharNotNL:
		// RE2 "." excludes only \n; ECMA "." also excludes \r, U+2028 and U+2029.
		b.WriteString(`[^\n]`)
	case syntax.OpAnyChar:
		b.WriteString(`[\s\S]`)
	case syntax.OpBeginText:
		b.WriteString("^")
	case syntax.OpEndText:
		// Without the "m" flag, ECMA "$" matches only at the end of input, as RE2
		// "\z" and non-multiline "$" do.
		b.WriteString("$")
	case syntax.OpWordBoundary:
		b.WriteString(`\b`)
	case syntax.OpNoWordBoundary:
		b.WriteString(`\B`)
	case syntax.OpCapture:
		b.WriteString("(?:")
		if !writeECMA(b, re.Sub[0]) {
			return false
		}
		b.WriteString(")")
	case syntax.OpStar, syntax.OpPlus, syntax.OpQuest, syntax.OpRepeat:
		if !writeECMAAtom(b, re.Sub[0]) {
			return false
		}
		switch re.Op {
		case syntax.OpStar:
			b.WriteString("*")
		case syntax.OpPlus:
			b.WriteString("+")
		case syntax.OpQuest:
			b.WriteString("?")
		default:
			switch {
			case re.Max == -1:
				fmt.Fprintf(b, "{%d,}", re.Min)
			case re.Min == re.Max:
				fmt.Fprintf(b, "{%d}", re.Min)
			default:
				fmt.Fprintf(b, "{%d,%d}", re.Min, re.Max)
			}
		}
		if re.Flags&syntax.NonGreedy != 0 {
			b.WriteString("?")
		}
	case syntax.OpConcat:
		for _, sub := range re.Sub {
			if sub.Op == syntax.OpAlternate {
				if !writeECMAGroup(b, sub) {
					return false
				}
				continue
			}
			if !writeECMA(b, sub) {
				return false
			}
		}
	case syntax.OpAlternate:
		for i, sub := range re.Sub {
			if i > 0 {
				b.WriteString("|")
			}
			if !writeECMA(b, sub) {
				return false
			}
		}
	default:
		// OpBeginLine and OpEndLine come from the RE2 "m" flag, which cannot be
		// expressed inside an ECMA pattern.
		return false
	}
	return true
}

// writeECMAAtom writes re so that a following quantifier applies to all of it.
func writeECMAAtom(b *strings.Builder, re *syntax.Regexp) bool {
	switch re.Op {
	case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar, syntax.OpCapture:
		return writeECMA(b, re)
	case syntax.OpLiteral:
		if len(re.Rune) == 1 {
			return writeECMA(b, re)
		}
	}
	return writeECMAGroup(b, re)
}

func writeECMAGroup(b *strings.Builder, re *syntax.Regexp) bool {
	b.WriteString("(?:")
	if !writeECMA(b, re) {
		return false
	}
	b.WriteString(")")
	return true
}

// ecmaSyntaxChars are the characters ECMA-262 allows as identity escapes in
// "u" mode, plus "-", which is escapable inside a class in "u" mode.
const ecmaSyntaxChars = `^$\.*+?()[]{}|/-`

func writeECMARune(b *strings.Builder, r rune, inClass bool) bool {
	switch {
	case r > 0xFFFF:
		// \u{...} is only valid with the "u" flag and a surrogate pair only
		// without it, so there is no escape that means the same in both.
		return false
	case r == '-' && !inClass:
		b.WriteRune(r)
	case strings.ContainsRune(ecmaSyntaxChars, r):
		b.WriteByte('\\')
		b.WriteRune(r)
	case r >= 0x20 && r < 0x7F:
		b.WriteRune(r)
	default:
		fmt.Fprintf(b, `\u%04X`, r)
	}
	return true
}

// writeECMAClass writes a character class from RE2's sorted [lo, hi] range
// pairs. RE2 stores negated classes as their complement, so a class reaching
// past the BMP is written as a negated class of what it excludes.
func writeECMAClass(b *strings.Builder, ranges []rune) bool {
	if len(ranges) == 2 && ranges[0] == 0 && ranges[1] == unicode.MaxRune {
		b.WriteString(`[\s\S]`)
		return true
	}
	negate := false
	if ranges[len(ranges)-1] > 0xFFFF {
		ranges = complementRanges(ranges)
		negate = true
		if len(ranges) > 0 && ranges[len(ranges)-1] > 0xFFFF {
			return false
		}
	}
	b.WriteString("[")
	if negate {
		b.WriteString("^")
	}
	for i := 0; i < len(ranges); i += 2 {
		lo, hi := ranges[i], ranges[i+1]
		if !writeECMARune(b, lo, true) {
			return false
		}
		if hi != lo {
			if hi > lo+1 {
				b.WriteString("-")
			}
			if !writeECMARune(b, hi, true) {
				return false
			}
		}
	}
	b.WriteString("]")
	return true
}

func complementRanges(ranges []rune) []rune {
	var out []rune
	next := rune(0)
	for i := 0; i < len(ranges); i += 2 {
		if ranges[i] > next {
			out = append(out, next, ranges[i]-1)
		}
		next = ranges[i+1] + 1
	}
	if next <= unicode.MaxRune {
		out = append(out, next, unicode.MaxRune)
	}
	return out
}

func hasFold(r rune) bool {
	return unicode.SimpleFold(r) != r
}

// foldRanges returns the sorted single-rune ranges of r's case-folding orbit.
func foldRanges(r rune) []rune {
	orbit := []rune{r}
	for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
		orbit = append(orbit, f)
	}
	for i := 1; i < len(orbit); i++ {
		for j := i; j > 0 && orbit[j] < orbit[j-1]; j-- {
			orbit[j], orbit[j-1] = orbit[j-1], orbit[j]
		}
	}
	out := make([]rune, 0, 2*len(orbit))
	for _, f := range orbit {
		out = append(out, f, f)
	}
	return out
}
