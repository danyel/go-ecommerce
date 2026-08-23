package factory

import (
	Reflect "reflect"
	Sync "sync"

	Logger "github.com/danyel/ecommerce/cmd/logger"
)

var (
	applicationConnectionFactoryInstance ApplicationConnectionFactory
	databaseConnectionFactoryInstance    DatabaseConnectionFactory
	brokerConnectionFactoryInstance      MessageBrokerConnectionFactory
	brokerConfigurationOnce              Sync.Once
	applicationConnectionOnce            Sync.Once
	databaseConnectionOnce               Sync.Once
)

func getInstanceOfType[T any](value *T, factoryMethod func() T) T {
	typeName := Reflect.TypeFor[T]().String()
	Logger.Log.Debug("[getInstanceOfType] value: %v", value)
	if isNil(any(*value)) {
		Logger.Log.Debug("[getInstanceOfType] Creating new instance for type: %s", typeName)
		*value = factoryMethod()
		return *value
	}
	Logger.Log.Debug("[getInstanceOfType] Returning existing instance of type: %s (pointer: %p)", typeName, any(*value))

	return *value
}

func isNil(value any) bool {
	if value == nil {
		return true
	}

	valueOf := Reflect.ValueOf(value)
	switch valueOf.Kind() {
	case Reflect.Chan, Reflect.Func, Reflect.Map, Reflect.Pointer, Reflect.UnsafePointer, Reflect.Interface, Reflect.Slice:
		return valueOf.IsNil()
	default:
		return false
	}
}
