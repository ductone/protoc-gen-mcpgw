package mcpgw

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

func New() pgs.Module {
	return &Module{ModuleBase: &pgs.ModuleBase{}}
}

const moduleName = "mcpgw"
const version = "0.1.0"

type Module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

var _ pgs.Module = (*Module)(nil)

func (m *Module) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Module) Name() string {
	return moduleName
}

func (m *Module) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		m.processFile(m.ctx, f)
	}
	return m.Artifacts()
}

func (m *Module) processFile(ctx pgsgo.Context, f pgs.File) {
	out := bytes.Buffer{}
	rendered, err := m.applyTemplate(ctx, &out, f)
	if err != nil {
		m.Logf("couldn't apply template: %s", err)
		m.Fail("code generation failed")
		return
	}
	// We didn't find anything to render so skip writing the file
	if !rendered {
		return
	}
	generatedFileName := m.ctx.OutputPath(f).SetExt(fmt.Sprintf(".%s.go", moduleName)).String()
	if ok, _ := strconv.ParseBool(os.Getenv("MCPGW_DEBUG_FILE_RAW")); ok {
		spew.Fdump(os.Stderr, out.String())
		_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", out.String())
	}
	m.AddGeneratorFile(generatedFileName, out.String())
}

func (m *Module) applyTemplate(ctx pgsgo.Context, w *bytes.Buffer, f pgs.File) (bool, error) {
	ix := &importTracker{
		ctx:        ctx,
		input:      f,
		typeMapper: make(map[pgs.Name]pgs.FilePath),
	}
	headerBuf := &bytes.Buffer{}
	bodyBuf := &bytes.Buffer{}

	services := f.Services()
	for _, service := range services {
		sopt := getServiceOptions(service)
		if sopt != nil && !sopt.GetEnabled() {
			continue
		}
		err := m.renderService(ctx, bodyBuf, f, service, ix)
		if err != nil {
			return false, err
		}
	}

	if len(services) == 0 {
		return false, nil
	}

	err := m.renderHeader(ctx, headerBuf, f, ix)
	if err != nil {
		return false, err
	}

	_, err = io.Copy(w, headerBuf)
	if err != nil {
		return false, err
	}

	_, err = io.Copy(w, bodyBuf)
	if err != nil {
		return false, err
	}

	if ok, _ := strconv.ParseBool(os.Getenv("MCPGW_DUMP_FILE")); ok {
		tdr := os.TempDir()
		_ = os.WriteFile(filepath.Join(tdr, "t.go"), w.Bytes(), 0600)
		_, _ = os.Stderr.WriteString(filepath.Join(tdr, "t.go") + "\n")
	}
	return true, nil
}
