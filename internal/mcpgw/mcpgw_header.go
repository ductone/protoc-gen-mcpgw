package mcpgw

import (
	"io"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type headerTemplateContext struct {
	Version     string
	PackageName string
	SourceFile  string
	Imports     *importTracker
}

func (module *Module) renderHeader(ctx pgsgo.Context, w io.Writer, in pgs.File, ix *importTracker) error {
	c := &headerTemplateContext{
		Version:     version,
		SourceFile:  in.Name().String(),
		PackageName: ctx.PackageName(in).String(),
		Imports:     ix,
	}

	return templates["header.tmpl"].Execute(w, c)
}
