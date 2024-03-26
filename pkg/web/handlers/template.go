package handlers

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"unterlagen/pkg/config"
	"unterlagen/views"
)

type templateMap = map[string]*template.Template

type TemplateExecutor interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

func NewTemplateExecutor() TemplateExecutor {
	if config.Get().Development {
		return newDevelopmentTemplateExecutor()
	}
	return newProductionTemplateExecutor()
}

type DevelopmentTemplateExecutor struct {
}

func newDevelopmentTemplateExecutor() *DevelopmentTemplateExecutor {
	return &DevelopmentTemplateExecutor{}
}

func (e DevelopmentTemplateExecutor) ExecuteTemplate(writer io.Writer, name string, data any) error {
	templates := processTemplates(os.DirFS("./views"))
	return get(templates, name).ExecuteTemplate(writer, name, data)
}

type ProductionTemplateExecutor struct {
	Templates map[string]*template.Template
}

func newProductionTemplateExecutor() *ProductionTemplateExecutor {
	templates := processTemplates(views.Templates)
	return &ProductionTemplateExecutor{Templates: templates}
}

func (e ProductionTemplateExecutor) ExecuteTemplate(wr io.Writer, name string, data any) error {
	return get(e.Templates, name).ExecuteTemplate(wr, name, data)
}

func load(root string, filesystem fs.FS) []string {
	var files []string
	err := fs.WalkDir(filesystem, root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

func processTemplates(filesystem fs.FS) templateMap {
	var files []string
	files = append(files, load("templates/fixtures", filesystem)...)
	files = append(files, load("templates/icons", filesystem)...)
	files = append(files, load("templates/layouts", filesystem)...)
	commonTemplate := template.Must(template.New("common").ParseFS(filesystem, files...))

	pages := load("templates/pages", filesystem)
	templates := make(map[string]*template.Template)
	for _, page := range pages {
		templates[page] = template.Must(template.Must(commonTemplate.Clone()).ParseFS(filesystem, page))
	}

	return templates
}

func get(templates templateMap, name string) *template.Template {
	page := fmt.Sprintf("templates/pages/%s", name)
	return templates[page]
}
