package basic

import (
	"net/http"
	"platform/logging"
	"platform/pipeline"
	//"platform/services"
)

/*
Создание компонента ПО промежуточного слоя ведения журналов

Этот компонент регистрирует основные сведения о запросе и ответе с
помощью службы Logger, созданной в главе 32. Интерфейс ResponseWriter не
предоставляет доступ к коду состояния, отправленному в ответе, поэтому
создается LoggingResponseWriter и передается следующему компоненту в
конвейере.
Этот компонент выполняет действия до и после вызова функции next,
регистрируя сообщение перед передачей запроса и регистрируя другое
сообщение, в котором выводится код состояния после обработки запроса.
Этот компонент получает службу Logger при обработке запроса. Я мог бы
получить Logger только один раз, но это работает только потому, что я знаю, что
Logger был зарегистрирован как одноэлементная служба. Вместо этого я
предпочитаю не делать предположений о жизненном цикле Logger, а это значит,
что я не получу неожиданных результатов, если жизненный цикл изменится в
будущем.
*/

type LoggingResponseWriter struct {
	statusCode int
	http.ResponseWriter
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

type LoggingComponent struct{}

func (lc *LoggingComponent) ImplementsProcessRequestWithServices() {}
func (lc *LoggingComponent) Init()                                 {}
func (lc *LoggingComponent) ProcessRequestWithServices(ctx *pipeline.ComponentContext,
	next func(*pipeline.ComponentContext), logger logging.Logger) {
	// var logger logging.Logger
	// err := services.GetServiceForContext(ctx.Request.Context(),
	// &logger)
	// if (err != nil) {
	// ctx.Error(err)
	// return
	// }

	loggingWriter := LoggingResponseWriter{0, ctx.ResponseWriter}
	ctx.ResponseWriter = &loggingWriter
	logger.Infof("REQ --- %v - %v", ctx.Request.Method,
		ctx.Request.URL)
	next(ctx)
	logger.Infof("RSP %v %v", loggingWriter.statusCode,
		ctx.Request.URL)
}
