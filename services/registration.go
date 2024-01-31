package services

import (
	"reflect"
	"sync"
)

/*
Тут функции, которые будут использоваться
в остальной части приложения для регистрации служб.

Функции AddTransient и AddScoped просто передают фабричную
функцию функции addService. Для жизненного цикла синглтона
требуется немного больше работы, и функция AddSingleton создает
оболочку вокруг фабричной функции, которая гарантирует, что она
выполняется только один раз, для первого запроса на разрешение
службы. Это гарантирует, что создан только один экземпляр структуры
реализации и что он не будет создан до тех пор, пока он не понадобится
в первый раз.
*/
func AddTransient(factoryFunc interface{}) (err error) {
	return addService(Transient, factoryFunc)
}

func AddScoped(factoryFunc interface{}) (err error) {
	return addService(Scoped, factoryFunc)
}

func AddSingleton(factoryFunc interface{}) (err error) {
	factoryFuncVal := reflect.ValueOf(factoryFunc)
	if factoryFuncVal.Kind() == reflect.Func && factoryFuncVal.Type().NumOut() == 1 {
		var results []reflect.Value
		once := sync.Once{}
		wrapper := reflect.MakeFunc(factoryFuncVal.Type(), func([]reflect.Value) []reflect.Value {
			once.Do(func() {
				results = invokeFunction(nil,
					factoryFuncVal)
			})
			return results
		})
		err = addService(Singleton, wrapper.Interface())
	}
	return
}
