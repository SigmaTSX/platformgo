package pipeline

import (
	"net/http"
	"platform/services"
	"reflect"
)

/*
Создание конвейера запросов

CreatePipeline принимает ряд компонентов и соединяет их для создания
функции, которая принимает указатель на структуру ComponentContext. Эта
функция вызывает метод ProcessRequest первого компонента в конвейере со
следующим аргументом, который вызывает метод ProcessRequest следующего
компонента. Эта цепочка передает структуру ComponentContext всем
компонентам по очереди, если только один из них не вызывает метод Error.
Запросы обрабатываются с помощью метода ProcessRequest, который создает
значение ComponentContext и использует его для запуска обработки запроса
*/
type RequestPipeline func(*ComponentContext)

var emptyPipeline RequestPipeline = func(*ComponentContext) { /* do
	nothing */
}

func CreatePipeline(components ...interface{}) RequestPipeline {
	f := emptyPipeline
	for i := len(components) - 1; i >= 0; i-- {
		currentComponent := components[i]
		services.Populate(currentComponent)
		nextFunc := f
		if servComp, ok := currentComponent.(ServicesMiddlwareComponent); ok {
			f = createServiceDependentFunction(currentComponent,
				nextFunc)
			servComp.Init()
		} else if stdComp, ok := currentComponent.(MiddlewareComponent); ok {
			f = func(context *ComponentContext) {
				if context.error == nil {
					stdComp.ProcessRequest(context, nextFunc)
				}
			}
			stdComp.Init()
		} else {
			panic("Value is not a middleware component")
		}
	}
	return f
}

func createServiceDependentFunction(component interface{}, nextFunc RequestPipeline) RequestPipeline {
	method := reflect.ValueOf(component).MethodByName("ProcessRequestWithServices")
	if method.IsValid() {
		return func(context *ComponentContext) {
			if context.error == nil {
				_, err :=
					services.CallForContext(context.Request.Context(),
						method.Interface(), context, nextFunc)
				if err != nil {
					context.Error(err)
				}
			}
		}
	} else {
		panic("No ProcessRequestWithServices method defined")
	}
}

func (pl RequestPipeline) ProcessRequest(req *http.Request, resp http.ResponseWriter) error {
	ctx := ComponentContext{
		Request:        req,
		ResponseWriter: resp,
	}
	pl(&ctx)
	return ctx.error
}
