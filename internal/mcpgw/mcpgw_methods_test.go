package mcpgw

import (
	"net/http"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"", ""},
		{"Apps", "apps"},
		{"Get", "get"},
		{"AppEntitlements", "app_entitlements"},
		{"APIKey", "api_key"},
		{"MultiGet", "multi_get"},
		{"HTTPRequest", "http_request"},
		{"alreadySnake", "already_snake"},
	}
	for _, c := range cases {
		if got := toSnakeCase(c.in); got != c.out {
			t.Errorf("toSnakeCase(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

func TestVerbClassification(t *testing.T) {
	cases := []struct {
		verb        string
		readOnly    bool
		destructive bool
		idempotent  bool
	}{
		{http.MethodGet, true, false, true},
		{http.MethodHead, true, false, true},
		{http.MethodOptions, true, false, true},
		{http.MethodPost, false, false, false},
		{http.MethodPut, false, false, true},
		{http.MethodPatch, false, false, false},
		{http.MethodDelete, false, true, true},
		{"get", true, false, true},
	}
	for _, c := range cases {
		if got := isReadOnlyVerb(c.verb); got != c.readOnly {
			t.Errorf("isReadOnlyVerb(%q) = %v, want %v", c.verb, got, c.readOnly)
		}
		if got := isDestructiveVerb(c.verb); got != c.destructive {
			t.Errorf("isDestructiveVerb(%q) = %v, want %v", c.verb, got, c.destructive)
		}
		if got := isIdempotentVerb(c.verb); got != c.idempotent {
			t.Errorf("isIdempotentVerb(%q) = %v, want %v", c.verb, got, c.idempotent)
		}
	}
}

func TestNormalizeComment(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"", ""},
		{" Get an app by ID.", "Get an app by ID."},
		{" Get an app by ID.\n Returns the app.", "Get an app by ID. Returns the app."},
		{"\n Line one.\n\n Line two.\n", "Line one. Line two."},
	}
	for _, c := range cases {
		if got := normalizeComment(c.in); got != c.out {
			t.Errorf("normalizeComment(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}
