package templates

import (
	"errors"
	"html/template"
	"platform/config"
	"sync"
)

/*
Чтобы добавить поддержку загрузки шаблонов и создания значения функции,
которое будет присвоено переменной getTemplates
Функция LoadTemplates загружает шаблон из места, указанного в файле
конфигурации. Существует также параметр конфигурации, который включает
перезагрузку для каждого запроса, что не следует делать в развернутом проекте,
но полезно во время разработки, поскольку это означает, что изменения в
шаблонах можно увидеть без перезапуска приложения.
*/
var once = sync.Once{}

func LoadTemplates(c config.Configuration) (err error) {
	path, ok := c.GetString("templates:path")
	if !ok {
		return errors.New("Cannot load template config")
	}
	reload := c.GetBoolDefault("templates:reload", false)
	once.Do(func() {
		doLoad := func() (t *template.Template) {
			t = template.New("htmlTemplates")
			t.Funcs(map[string]interface{}{
				"body":   func() string { return "" },
				"layout": func() string { return "" },
			})
			t, err = t.ParseGlob(path)
			return
		}
		if reload {
			getTemplates = doLoad
		} else {
			var templates *template.Template
			templates = doLoad()
			getTemplates = func() *template.Template {
				t, _ := templates.Clone()
				return t
			}
		}
	})
	return
}
