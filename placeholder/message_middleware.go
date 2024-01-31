package placeholder

import (
	//"errors"
	//"io"

	"platform/config"
	"platform/pipeline"
	"platform/templates"
)

/*
Создание компонента ответа-заполнителя
на данный момент мне нужен компонентзаполнитель, который будет генерировать простые ответы по мере разработки
других функций.
Этот компонент выдает простой текстовый ответ, которого достаточно,
чтобы убедиться, что конвейер работает должным образом.

Компонент теперь реализует метод ProcessRequestWithServices и получает
службы посредством внедрения зависимостей. Одна из запрашиваемых служб
— это реализация интерфейса TemplateExecutor, который используется для
отображения шаблона simple_message.html.
*/
type SimpleMessageComponent struct {
	Message string
	config.Configuration
}

func (lc *SimpleMessageComponent) ImplementsProcessRequestWithServices() {}
func (c *SimpleMessageComponent) Init() {
	c.Message = c.Configuration.GetStringDefault("main:message",
		"Default Message")
}
func (c *SimpleMessageComponent) ProcessRequestWithServices(
	ctx *pipeline.ComponentContext,
	next func(*pipeline.ComponentContext), executor templates.TemplateExecutor) {
	err := executor.ExecTemplate(ctx.ResponseWriter,
		"simple_message.html", c.Message)
	if err != nil {
		ctx.Error(err)
	} else {
		next(ctx)
	}
}
