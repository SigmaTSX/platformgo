package templates

import (
	"html/template"
	"io"
	"strings"
)

type LayoutTemplateProcessor struct{}

/*
Реализация метода ExecTemplate выполняет шаблон и сохраняет
содержимое в файле strings.Builder.
*/
func (proc *LayoutTemplateProcessor) ExecTemplate(writer io.Writer,
	name string, data interface{}) (err error) {
	var sb strings.Builder
	layoutName := ""
	localTemplates := getTemplates()
	localTemplates.Funcs(map[string]interface{}{
		"body":   insertBodyWrapper(&sb),
		"layout": setLayoutWrapper(&layoutName),
	})
	err = localTemplates.ExecuteTemplate(&sb, name, data)
	if layoutName != "" {
		localTemplates.ExecuteTemplate(writer, layoutName, data)
	} else {
		io.WriteString(writer, sb.String())
	}
	return
}

var getTemplates func() (t *template.Template)

func insertBodyWrapper(body *strings.Builder) func() template.HTML {
	return func() template.HTML {
		return template.HTML(body.String())
	}
}
func setLayoutWrapper(val *string) func(string) string {
	return func(layout string) string {
		*val = layout
		return ""
	}
}
