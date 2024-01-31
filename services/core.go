package services

import (
	"context"
	"fmt"
	"reflect"
)

/*
Структура BindingMap представляет собой комбинацию фабричной
функции, выраженной в виде Reflect.Value, и жизненного цикла

Функция addService используется для регистрации службы, что она
делает путем создания BindingMap и добавления к карте, назначенной
переменной services.

Функция resolveServiceFromValue вызывается для разрешения
службы, а ее аргументы — это Context и Value, являющиеся указателем
на переменную, тип которой является интерфейсом, который нужно
решить

Чтобы разрешить службу, функция
getServiceFromValue проверяет, есть ли BindingMap в карте служб,
используя запрошенный тип в качестве ключа. Если есть BindingMap, то
вызывается его фабричная функция, и значение присваивается через
указатель.

Функция invokeFunction отвечает за вызов фабричной функции,
используя функцию resolveFunctionArguments для проверки
параметров фабричной функции и разрешения каждого из них. Эти
функции принимают необязательные дополнительные аргументы,
которые используются, когда функция должна быть вызвана с
сочетанием служб и параметров с обычными значениями

Для служб с заданной областью требуется особое обращение.
ResolveScopedService проверяет, содержит ли Context значение из
предыдущего запроса на разрешение службы. Если нет, служба
разрешается и добавляется в Context, чтобы ее можно было повторно
использовать в той же области.
*/
type BindingMap struct {
	factoryFunc reflect.Value
	lifecycle
}

var services = make(map[reflect.Type]BindingMap)

func addService(life lifecycle, factoryFunc interface{}) (err error) {
	factoryFuncType := reflect.TypeOf(factoryFunc)
	if factoryFuncType.Kind() == reflect.Func && factoryFuncType.NumOut() == 1 {
		services[factoryFuncType.Out(0)] = BindingMap{
			factoryFunc: reflect.ValueOf(factoryFunc),
			lifecycle:   life,
		}
	} else {
		err = fmt.Errorf("Type cannot be used as service: %v",
			factoryFuncType)
	}
	return
}

var contextReference = (*context.Context)(nil)
var contextReferenceType = reflect.TypeOf(contextReference).Elem()

func resolveServiceFromValue(c context.Context, val reflect.Value) (err error) {
	serviceType := val.Elem().Type()
	if serviceType == contextReferenceType {
		val.Elem().Set(reflect.ValueOf(c))
	} else if binding, found := services[serviceType]; found {
		if binding.lifecycle == Scoped {
			resolveScopedService(c, val, binding)
		} else {
			val.Elem().Set(invokeFunction(c,
				binding.factoryFunc)[0])
		}
	} else {
		err = fmt.Errorf("Cannot find service %v",
			serviceType)
	}
	return
}

func resolveScopedService(c context.Context, val reflect.Value, binding BindingMap) (err error) {
	sMap, ok := c.Value(ServiceKey).(serviceMap)
	if ok {
		serviceVal, ok := sMap[val.Type()]
		if !ok {
			serviceVal = invokeFunction(c,
				binding.factoryFunc)[0]
			sMap[val.Type()] = serviceVal
		}
		val.Elem().Set(serviceVal)
	} else {
		val.Elem().Set(invokeFunction(c, binding.factoryFunc)[0])
	}
	return
}

func resolveFunctionArguments(c context.Context, f reflect.Value, otherArgs ...interface{}) []reflect.Value {
	params := make([]reflect.Value, f.Type().NumIn())
	i := 0
	if otherArgs != nil {
		for ; i < len(otherArgs); i++ {
			params[i] = reflect.ValueOf(otherArgs[i])
		}
	}
	for ; i < len(params); i++ {
		pType := f.Type().In(i)
		pVal := reflect.New(pType)
		err := resolveServiceFromValue(c, pVal)
		if err != nil {
			panic(err)
		}
		params[i] = pVal.Elem()
	}
	return params
}

func invokeFunction(c context.Context, f reflect.Value, otherArgs ...interface{}) []reflect.Value {
	return f.Call(resolveFunctionArguments(c, f,
		otherArgs...))
}
