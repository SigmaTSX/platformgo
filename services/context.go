package services

import (
	"context"
	"reflect"
)

/*
Чтобы упростить работу с контекстами

Функция NewServiceContext извлекает контекст с помощью
функции WithValue, добавляя карту, в которой хранятся службы,
которые были разрешены
*/
const ServiceKey = "services"

type serviceMap map[reflect.Type]reflect.Value

func NewServiceContext(c context.Context) context.Context {
	if c.Value(ServiceKey) == nil {
		return context.WithValue(c, ServiceKey,
			make(serviceMap))
	} else {
		return c
	}
}
