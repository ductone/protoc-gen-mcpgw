package mcpgw

import (
	"net/http"
	"testing"

	apigw_v1 "github.com/ductone/protoc-gen-apigw/apigw/v1"
	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
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
		{"Http2Stream", "http2_stream"},
		{"V1Get", "v1_get"},
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

func TestVerbRouteFallback(t *testing.T) {
	cases := []struct {
		name string
		op   *apigw_v1.Operation
		want string
	}{
		{"nil operation", nil, ""},
		{"empty fields", &apigw_v1.Operation{}, ""},
		{"missing route", &apigw_v1.Operation{Method: "GET"}, ""},
		{"missing verb", &apigw_v1.Operation{Route: "/v1/apps"}, ""},
		{"verb uppercased", &apigw_v1.Operation{Method: "get", Route: "/v1/apps"}, "GET /v1/apps"},
		{"already upper", &apigw_v1.Operation{Method: "DELETE", Route: "/v1/apps/{id}"}, "DELETE /v1/apps/{id}"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := verbRouteFallback(c.op); got != c.want {
				t.Errorf("verbRouteFallback(%+v) = %q, want %q", c.op, got, c.want)
			}
		})
	}
}

func TestMergeMethodMetadata(t *testing.T) {
	const derivedTitle = "apps_get"
	const derivedDescription = "GET /v1/apps"

	cases := []struct {
		name           string
		mext           *mcpgw_v1.MethodOptions
		aop            *apigw_v1.Operation
		wantTitle      string
		wantDesc       string
		wantReadOnly   bool
		wantDestruct   bool
		wantIdempotent bool
		wantOpenWorld  bool
	}{
		{
			name:           "no annotations falls back to derived",
			mext:           nil,
			aop:            nil,
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantReadOnly:   false,
			wantDestruct:   false,
			wantIdempotent: false,
			wantOpenWorld:  false,
		},
		{
			name:           "apigw GET derives read-only and idempotent",
			mext:           nil,
			aop:            &apigw_v1.Operation{Method: "GET", Route: "/v1/apps"},
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantReadOnly:   true,
			wantDestruct:   false,
			wantIdempotent: true,
			wantOpenWorld:  false,
		},
		{
			name:           "apigw DELETE derives destructive and idempotent",
			mext:           nil,
			aop:            &apigw_v1.Operation{Method: "DELETE", Route: "/v1/apps/{id}"},
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantReadOnly:   false,
			wantDestruct:   true,
			wantIdempotent: true,
			wantOpenWorld:  false,
		},
		{
			name:           "apigw POST derives nothing",
			mext:           nil,
			aop:            &apigw_v1.Operation{Method: "POST", Route: "/v1/apps"},
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantReadOnly:   false,
			wantDestruct:   false,
			wantIdempotent: false,
			wantOpenWorld:  false,
		},
		{
			name: "explicit title and description override derived",
			mext: mcpgw_v1.MethodOptions_builder{
				Title:       ptr("Custom Title"),
				Description: ptr("Custom desc"),
			}.Build(),
			aop:            &apigw_v1.Operation{Method: "GET", Route: "/v1/apps"},
			wantTitle:      "Custom Title",
			wantDesc:       "Custom desc",
			wantReadOnly:   true,
			wantIdempotent: true,
		},
		{
			name: "explicit empty title and description override derived",
			mext: mcpgw_v1.MethodOptions_builder{
				Title:       ptr(""),
				Description: ptr(""),
			}.Build(),
			aop:            &apigw_v1.Operation{Method: "GET", Route: "/v1/apps"},
			wantTitle:      "",
			wantDesc:       "",
			wantReadOnly:   true,
			wantIdempotent: true,
		},
		{
			name: "explicit read_only_hint=false suppresses GET-derived true",
			mext: mcpgw_v1.MethodOptions_builder{
				ReadOnlyHint: ptr(false),
			}.Build(),
			aop:            &apigw_v1.Operation{Method: "GET", Route: "/v1/apps"},
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantReadOnly:   false,
			wantIdempotent: true,
		},
		{
			name: "explicit destructive_hint=false suppresses DELETE-derived true",
			mext: mcpgw_v1.MethodOptions_builder{
				DestructiveHint: ptr(false),
			}.Build(),
			aop:            &apigw_v1.Operation{Method: "DELETE", Route: "/v1/apps/{id}"},
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantDestruct:   false,
			wantIdempotent: true,
		},
		{
			name: "explicit idempotent_hint=false suppresses PUT-derived true",
			mext: mcpgw_v1.MethodOptions_builder{
				IdempotentHint: ptr(false),
			}.Build(),
			aop:            &apigw_v1.Operation{Method: "PUT", Route: "/v1/apps/{id}"},
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantIdempotent: false,
		},
		{
			name: "explicit hints=true on POST add hints apigw wouldn't derive",
			mext: mcpgw_v1.MethodOptions_builder{
				ReadOnlyHint:   ptr(true),
				DestructiveHint: ptr(true),
				IdempotentHint:  ptr(true),
				OpenWorldHint:   ptr(true),
			}.Build(),
			aop:            &apigw_v1.Operation{Method: "POST", Route: "/v1/apps"},
			wantTitle:      derivedTitle,
			wantDesc:       derivedDescription,
			wantReadOnly:   true,
			wantDestruct:   true,
			wantIdempotent: true,
			wantOpenWorld:  true,
		},
		{
			name: "mcpgw-only (no apigw) takes explicit fields verbatim",
			mext: mcpgw_v1.MethodOptions_builder{
				Title:        ptr("Manual"),
				ReadOnlyHint: ptr(true),
			}.Build(),
			aop:            nil,
			wantTitle:      "Manual",
			wantDesc:       derivedDescription,
			wantReadOnly:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			title, desc, ro, destr, idem, ow := mergeMethodMetadata(c.mext, c.aop, derivedTitle, derivedDescription)
			if title != c.wantTitle {
				t.Errorf("title = %q, want %q", title, c.wantTitle)
			}
			if desc != c.wantDesc {
				t.Errorf("description = %q, want %q", desc, c.wantDesc)
			}
			if ro != c.wantReadOnly {
				t.Errorf("readOnly = %v, want %v", ro, c.wantReadOnly)
			}
			if destr != c.wantDestruct {
				t.Errorf("destructive = %v, want %v", destr, c.wantDestruct)
			}
			if idem != c.wantIdempotent {
				t.Errorf("idempotent = %v, want %v", idem, c.wantIdempotent)
			}
			if ow != c.wantOpenWorld {
				t.Errorf("openWorld = %v, want %v", ow, c.wantOpenWorld)
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
