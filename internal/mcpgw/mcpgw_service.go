package mcpgw

import (
	"fmt"
	"io"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type serviceTemplateContext struct {
	ServerName         string
	FullyQualifiedName string
	Methods            []*methodTemplateContext
}

func (module *Module) renderService(ctx pgsgo.Context, w io.Writer, f pgs.File, in pgs.Service, ix *importTracker) error {
	ix.MCPGWV1 = true

	c := &serviceTemplateContext{
		ServerName: ctx.ServerName(in).String(),
		FullyQualifiedName: fmt.Sprintf("%s.%s",
			in.Package().ProtoName().String(),
			in.Name().String(),
		),
	}
	for _, method := range in.Methods() {
		methodCtx, err := module.methodContext(ctx, w, f, in, method, ix)
		if err != nil {
			return fmt.Errorf("method generation failed [%st: %w", method.FullyQualifiedName(), err)
		}
		if methodCtx == nil {
			continue
		}
		c.Methods = append(c.Methods, methodCtx)
	}

	return templates["service.tmpl"].Execute(w, c)
}
