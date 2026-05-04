package mcpgw

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	apigw_v1 "github.com/ductone/protoc-gen-apigw/apigw/v1"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	mcpgw_v1 "github.com/ductone/protoc-gen-mcpgw/mcpgw/v1"
)

type methodTemplateContext struct {
	mcpgw_v1.MethodDesc
	MethodHandlerName      string
	DecoderHandlerName     string
	InputSchemaHandlerName string
	RequestType            string
	ServerName             string
	MethodName             string
	FullMethodName         string
}

// methodContext builds the template context for a single RPC method.
//
// A method is renderable as an MCP tool when it has either:
//  1. An explicit (mcpgw.v1.method) annotation — fields taken verbatim.
//  2. An (apigw.v1.method) annotation with at least one Operation — tool
//     metadata is auto-derived from the HTTP verb, route, and proto
//     leading comment so any REST endpoint automatically becomes an MCP
//     tool without additional annotation effort.
//
// If the method has neither annotation the method is skipped by returning
// (nil, nil); callers filter out nil entries.
func (module *Module) methodContext(ctx pgsgo.Context, w io.Writer, f pgs.File, service pgs.Service, method pgs.Method, ix *importTracker) (*methodTemplateContext, error) {
	mext := getMethodOptions(method)
	aop := getApigwOperation(method)

	if mext == nil && aop == nil {
		return nil, nil
	}

	ix.Protojson = true
	ix.MCPGWV1 = true
	ix.MCPGWV1Schema = true

	serviceShortName := strings.TrimSuffix(ctx.Name(service).String(), "Server")
	methodShortName := ctx.Name(method).String()
	methodFullName := fmt.Sprintf("%s_%s_FullMethodName", serviceShortName, methodShortName)
	ix.Context = true
	ix.Proto = true
	ix.Protojson = true
	ix.GRPC = true

	title, description, readOnly, destructive, idempotent, openWorld := mergeMethodMetadata(mext, aop, deriveToolName(service, method), deriveDescription(method, aop))

	rv := &methodTemplateContext{
		MethodDesc: mcpgw_v1.MethodDesc{
			Method:        methodFullName,
			Title:         title,
			Description:   description,
			ReadOnlyHint:  readOnly,
			Destructive:   destructive,
			Idempotent:    idempotent,
			OpenWorldHint: openWorld,
		},
		ServerName:     ctx.ServerName(service).String(),
		MethodName:     methodShortName,
		FullMethodName: methodFullName,
		MethodHandlerName: fmt.Sprintf("_%s_%s_MCPGW_Handler",
			serviceShortName,
			methodShortName,
		),
		DecoderHandlerName: fmt.Sprintf("_%s_%s_MCPGW_Decoder",
			serviceShortName,
			methodShortName,
		),
		InputSchemaHandlerName: fmt.Sprintf("_%s_%s_MCPGW_InputSchema",
			serviceShortName,
			methodShortName,
		),
		RequestType: ctx.Name(method.Input()).String(),
	}
	return rv, nil
}

// mergeMethodMetadata combines explicit (mcpgw.v1.method) options with
// fallbacks derived from (apigw.v1.method). Each field falls back to the
// apigw-derived default when the mcpgw option does not set it explicitly.
//
// Explicit presence wins: a hint set to false in the mcpgw annotation
// suppresses the verb-derived default (e.g. read_only_hint = false on a
// GET endpoint). Only unset fields fall through to the apigw verb.
//
// derivedTitle and derivedDescription are precomputed so the merge logic
// stays a pure function of the proto options — see deriveToolName /
// deriveDescription for the pgs-dependent halves.
func mergeMethodMetadata(mext *mcpgw_v1.MethodOptions, aop *apigw_v1.Operation, derivedTitle, derivedDescription string) (title, description string, readOnly, destructive, idempotent, openWorld bool) {
	if mext != nil && mext.HasTitle() {
		title = mext.GetTitle()
	} else {
		title = derivedTitle
	}

	if mext != nil && mext.HasDescription() {
		description = mext.GetDescription()
	} else {
		description = derivedDescription
	}

	verb := ""
	if aop != nil {
		verb = aop.GetMethod()
	}
	readOnly = resolveHint(mext, (*mcpgw_v1.MethodOptions).HasReadOnlyHint, (*mcpgw_v1.MethodOptions).GetReadOnlyHint, verb, isReadOnlyVerb)
	destructive = resolveHint(mext, (*mcpgw_v1.MethodOptions).HasDestructiveHint, (*mcpgw_v1.MethodOptions).GetDestructiveHint, verb, isDestructiveVerb)
	idempotent = resolveHint(mext, (*mcpgw_v1.MethodOptions).HasIdempotentHint, (*mcpgw_v1.MethodOptions).GetIdempotentHint, verb, isIdempotentVerb)
	if mext != nil && mext.HasOpenWorldHint() {
		openWorld = mext.GetOpenWorldHint()
	}
	return
}

// resolveHint picks an explicit mcpgw bool when present, otherwise falls
// back to deriving from the HTTP verb. An empty verb (no apigw annotation)
// yields false.
func resolveHint(
	mext *mcpgw_v1.MethodOptions,
	has func(*mcpgw_v1.MethodOptions) bool,
	get func(*mcpgw_v1.MethodOptions) bool,
	verb string,
	derive func(string) bool,
) bool {
	if mext != nil && has(mext) {
		return get(mext)
	}
	if verb == "" {
		return false
	}
	return derive(verb)
}

// deriveToolName returns "{service}_{method}" in snake_case.
// Example: service "Apps", method "Get" -> "apps_get".
func deriveToolName(service pgs.Service, method pgs.Method) string {
	return toSnakeCase(service.Name().String()) + "_" + toSnakeCase(method.Name().String())
}

// deriveDescription uses the proto leading comment on the RPC as the MCP
// tool description. If no comment is present, falls back to the HTTP verb
// and route so the LLM at least has some usable metadata.
func deriveDescription(method pgs.Method, aop *apigw_v1.Operation) string {
	if sci := method.SourceCodeInfo(); sci != nil {
		if c := strings.TrimSpace(sci.LeadingComments()); c != "" {
			return normalizeComment(c)
		}
	}
	return verbRouteFallback(aop)
}

// verbRouteFallback formats an apigw operation as "VERB /route", or
// returns "" if either field is missing.
func verbRouteFallback(aop *apigw_v1.Operation) string {
	if aop == nil {
		return ""
	}
	verb := strings.ToUpper(aop.GetMethod())
	route := aop.GetRoute()
	if verb == "" || route == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", verb, route)
}

// normalizeComment strips protoc's per-line leading space and joins the
// comment into a single space-separated string so it's suitable as a
// JSON-schema description field.
func normalizeComment(c string) string {
	lines := strings.Split(c, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimPrefix(ln, " ")
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		out = append(out, ln)
	}
	return strings.Join(out, " ")
}

func isReadOnlyVerb(verb string) bool {
	switch strings.ToUpper(verb) {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}

func isDestructiveVerb(verb string) bool {
	return strings.ToUpper(verb) == http.MethodDelete
}

// isIdempotentVerb follows RFC 9110 §9.2.2: GET, HEAD, OPTIONS, PUT, DELETE
// are idempotent. POST and PATCH are not.
func isIdempotentVerb(verb string) bool {
	switch strings.ToUpper(verb) {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodDelete:
		return true
	}
	return false
}

// toSnakeCase converts an UpperCamelCase or camelCase identifier to
// snake_case. It handles consecutive uppercase runs (e.g. "APIKey" ->
// "api_key") by treating a lowercase letter following an uppercase run as
// the start of a new word.
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 4)
	runes := []rune(s)
	for i, r := range runes {
		isUpper := r >= 'A' && r <= 'Z'
		if isUpper && i > 0 {
			prev := runes[i-1]
			prevUpper := prev >= 'A' && prev <= 'Z'
			var nextLower bool
			if i+1 < len(runes) {
				nr := runes[i+1]
				nextLower = nr >= 'a' && nr <= 'z'
			}
			if !prevUpper || nextLower {
				b.WriteByte('_')
			}
		}
		if isUpper {
			b.WriteRune(r - 'A' + 'a')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
