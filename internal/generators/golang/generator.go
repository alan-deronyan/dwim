package golang

import (
	"fmt"
	"go/format"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/deronyan-llc/dwim/internal/common"
	"github.com/deronyan-llc/dwim/internal/generators"
)

type GoSrcGenerator struct {
	generators.SrcGenerator
}

// combine all the imports and types from all the classes together
type GoStructs struct {
	Comment string
	Package string
	Imports []string
	Classes []*common.SchemaClass
}

// generateGoCode generates GoLang code based on the given SchemaClass map.
func (g GoSrcGenerator) Generate(schemaContext *common.SchemaContext) error {
	const tmpl = `package {{.Package}}

/*
{{ .Comment }}
*/

{{- if .Imports }}

import (
{{ .Imports | imports }}
)


{{- end }}

{{- if .Classes }}
	{{- range .Classes }}

// {{ .TermString | localName | sanitizeName | title }} is a generated struct representing the <{{ .TermString }}> class.
type {{ .TermString | localName | sanitizeName | title }} struct {
		{{- if .Properties }}
			{{- range .Properties }}
	{{ .TermString | localName | sanitizeName | title }} {{ .GoLangType }} ` +
		"`json:\"{{ .TermString | localName | sanitizeName | lower }}\"`" +
		"{{ .Comment }}" + `
			{{- end }}
		{{- end }}
}

	{{- end }}
{{- end }}
`
	funcMap := template.FuncMap{
		"title":        strings.Title,
		"lower":        strings.ToLower,
		"localName":    common.LocalName,
		"sanitizeName": common.SanitizeName,
		"imports":      common.FormatImports,
	}

	t, err := template.New("class").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

	goSrcFileName := fmt.Sprintf("%s/%s.go",
		schemaContext.OutputPath,
		schemaContext.SchemaPath.StrippedFileName,
	)

	goStructs := &GoStructs{
		Comment: fmt.Sprintf("Generated code. DO NOT EDIT. %s", time.Now().Format(time.RFC3339)),
		Package: schemaContext.SchemaPath.StrippedFileName,
		Imports: []string{},
		Classes: []*common.SchemaClass{},
	}
	for _, classes := range schemaContext.Classes {
		for _, importName := range classes.Imports {
			if !slices.Contains(goStructs.Imports, importName) {
				goStructs.Imports = append(goStructs.Imports, importName)
			}
		}
		goStructs.Classes = append(goStructs.Classes, classes)
	}

	// Generate a Go source for the classes...

	// execute the template
	stringBuf := strings.Builder{}
	err = t.Execute(&stringBuf, goStructs)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	// create the output directory if it doesn't exist
	os.MkdirAll(schemaContext.OutputPath, os.ModePerm)
	file, err := os.Create(goSrcFileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// write the generated Go source to a buffer, unformatted initially
	_, err = file.Write([]byte(stringBuf.String()))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	// format the Go source code buffer
	formatted, err := format.Source([]byte(stringBuf.String()))
	if err != nil {
		fmt.Println("Error formatting source:", err)
		return err
	}

	// write the formatted Go buffer to a file
	_, err = file.WriteAt(formatted, 0)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}
