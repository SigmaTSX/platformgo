package services

import (
	"context"
	"errors"
	"reflect"
)

/*
Следующий набор функций включает в себя функции, позволяющие
разрешать службы.
GetServiceForContext принимает контекст и указатель на значение,
которое можно установить с помощью рефлексии. Для удобства
функция GetService разрешает службу, используя фоновый контекст.
*/
func GetService(target interface{}) error {
	return GetServiceForContext(context.Background(), target)
}
func GetServiceForContext(c context.Context, target interface{}) (err error) {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() == reflect.Ptr &&
		targetValue.Elem().CanSet() {
		err = resolveServiceFromValue(c, targetValue)
	} else {
		err = errors.New("Type cannot be used as target")
	}
	return
}
