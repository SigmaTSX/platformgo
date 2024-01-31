package basic

import (
	"platform/pipeline"
	"platform/services"
)

/*
Создание сервисов компонента ПО промежуточного слоя

Этот компонент промежуточного программного обеспечения изменяет
Context, связанный с запросом, чтобы во время обработки запроса можно было
использовать контекстно-зависимые службы. Метод http.Request.Context
используется для получения стандартного Context, созданного с помощью
запроса, который подготавливается для служб, а затем обновляется с помощью
метода WithContext.
После подготовки контекста запрос передается по конвейеру путем вызова
функции, полученной через параметр с именем next
*/
type ServicesComponent struct{}

func (c *ServicesComponent) Init() {}
func (c *ServicesComponent) ProcessRequest(ctx *pipeline.ComponentContext,
	next func(*pipeline.ComponentContext)) {
	reqContext := ctx.Request.Context()
	ctx.Request.WithContext(services.NewServiceContext(reqContext))
	next(ctx)
}
