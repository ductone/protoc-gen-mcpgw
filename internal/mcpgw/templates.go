package mcpgw

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"text/template"
)

// templateFuncs are registered on every parsed template. goquote wraps
// arbitrary strings (e.g. proto leading comments used as MCP tool
// descriptions) into a Go double-quoted string literal with embedded
// quotes, backslashes, and control characters escaped. Templates should
// invoke it as `{{ .Field | goquote }}` and omit surrounding double
// quotes; goquote provides them.
var templateFuncs = template.FuncMap{
	"goquote": strconv.Quote,
}

var (
	//go:embed templates/*.tmpl
	templateFiles embed.FS
	templates     map[string]*template.Template
)

//nolint:gochecknoinits // compling templates from embed // This linter is broken currently https://github.com/golangci/golangci-lint/issues/4697
func init() {
	err := loadTemplates()
	if err != nil {
		panic(fmt.Errorf("pgdb.loadTemplates failed; %w", err))
	}
}

func templateExecToString(name string, c interface{}) (string, error) {
	buf := bytes.Buffer{}
	err := templates[name].Execute(&buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func loadTemplates() error {
	templates = make(map[string]*template.Template)
	tmplFiles, err := fs.ReadDir(templateFiles, "templates")
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.New(tmpl.Name()).Funcs(templateFuncs).ParseFS(templateFiles, path.Join("templates", tmpl.Name()))
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

type varNamer struct {
	prefix string
	offset int
}

func (fn *varNamer) Next() *varNamer {
	fn.offset += 1
	return fn
}

func (fn *varNamer) String() string {
	return fmt.Sprintf("%s%d", fn.prefix, fn.offset)
}
