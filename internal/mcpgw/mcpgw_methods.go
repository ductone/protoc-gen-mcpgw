package mcpgw

import (
	"fmt"
	"io"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type methodTemplateContext struct {
	Method                 string
	Title                  string
	Description            string
	ReadOnlyHint           bool
	Destructive            bool
	Idempotent             bool
	OpenWorldHint          bool
	MethodHandlerName      string
	DecoderHandlerName     string
	InputSchemaHandlerName string
	RequestType            string
	ResponseType           string
	ServerName             string
	MethodName             string
	FullMethodName         string
}

func (module *Module) methodContext(ctx pgsgo.Context, w io.Writer, f pgs.File, service pgs.Service, method pgs.Method, ix *importTracker) (*methodTemplateContext, error) {
	mext := getMethodOptions(method)
	if mext == nil {
		return nil, fmt.Errorf("apigw: methodContext: failed to extract Method extension from '%s' (on enabled service '%s')", method.FullyQualifiedName(), service.FullyQualifiedName())
	}

	ix.Protojson = true
	ix.MCPGWV1 = true
	ix.MCPGWV1Schema = true

	serviceShortName := strings.TrimSuffix(ctx.Name(service).String(), "Server")
	methodFullName := fmt.Sprintf("%s_%s_FullMethodName", serviceShortName, ctx.Name(method).String())
	ix.Context = true
	ix.Proto = true
	ix.Protojson = true
	ix.GRPC = true

	rv := &methodTemplateContext{
		Method:         methodFullName,
		Title:          mext.GetTitle(),
		Description:    mext.GetDescription(),
		ReadOnlyHint:   mext.GetReadOnlyHint(),
		Destructive:    mext.GetDestructiveHint(),
		Idempotent:     mext.GetIdempotentHint(),
		OpenWorldHint:  mext.GetOpenWorldHint(),
		ServerName:     ctx.ServerName(service).String(),
		MethodName:     ctx.Name(method).String(),
		FullMethodName: methodFullName,
		MethodHandlerName: fmt.Sprintf("_%s_%s_MCPGW_Handler",
			serviceShortName,
			ctx.Name(method).String(),
		),
		DecoderHandlerName: fmt.Sprintf("_%s_%s_MCPGW_Decoder",
			serviceShortName,
			ctx.Name(method).String(),
		),
		InputSchemaHandlerName: fmt.Sprintf("_%s_%s_MCPGW_InputSchema",
			serviceShortName,
			ctx.Name(method).String(),
		),
		RequestType:  ctx.Name(method.Input()).String(),
		ResponseType: ctx.Name(method.Output()).String(),
	}
	return rv, nil
}
