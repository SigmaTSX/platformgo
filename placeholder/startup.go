package placeholder

import (
	"platform/http"
	"platform/pipeline"
	"platform/pipeline/basic"
	"platform/services"
	"sync"
)

/*
Функция createPipeline создает конвейер с ранее созданными
компонентами промежуточного ПО. Функция Start вызывает createPipeline и
использует результат для настройки и запуска HTTP-сервера.
*/

func createPipeline() pipeline.RequestPipeline {
	return pipeline.CreatePipeline(
		&basic.ServicesComponent{},
		&basic.LoggingComponent{},
		&basic.ErrorComponent{},
		&basic.StaticFileComponent{},
		&SimpleMessageComponent{},
	)
}
func Start() {
	results, err := services.Call(http.Serve, createPipeline())
	if err == nil {
		(results[0].(*sync.WaitGroup)).Wait()
	} else {
		panic(err)
	}
}
