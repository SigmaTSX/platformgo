package templates

import "io"

/*
Реализация выполнения шаблона
Интерфейс TemplateProcessor определяет метод с именем ExecTemplate,
который обрабатывает шаблон с использованием предоставленных значений
данных и записывает содержимое в Writer.
*/
type TemplateExecutor interface {
	ExecTemplate(writer io.Writer, name string, data interface{}) (err error)
}
