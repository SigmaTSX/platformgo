package pipeline

import (
	"net/http"
)

/*
Определение интерфейса компонента промежуточного
программного обеспечения
Как следует из названия, интерфейс MiddlewareComponent описывает
функциональные возможности, необходимые компоненту промежуточного
программного обеспечения. Метод Init используется для выполнения любой
одноразовой настройки, а другой метод с именем ProcessRequest отвечает за
обработку HTTP-запросов. Параметры, определенные методом ProcessRequest,
представляют собой указатель на структуру ComponentContext и функцию,
которая передает запрос следующему компоненту в конвейере

Оптимизация разрешения сервиса ServicesMiddlwareComponent
интерфейс, который позволит компонентам указывать, что им
требуется внедрение зависимостей для обработки запросов, как показано в
листинге
*/

type ComponentContext struct {
	*http.Request
	http.ResponseWriter
	error
}

func (mwc *ComponentContext) Error(err error) {
	mwc.error = err
}
func (mwc *ComponentContext) GetError() error {
	return mwc.error
}

type MiddlewareComponent interface {
	Init()
	ProcessRequest(context *ComponentContext, next func(*ComponentContext))
}

type ServicesMiddlwareComponent interface {
	Init()
	ImplementsProcessRequestWithServices()
}
