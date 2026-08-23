package factory

import (
	Log "log"
	Reflect "reflect"
	Sync "sync"
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
	Log.Printf("[getInstanceOfType] value: %v", value)
	if isNil(any(*value)) {
		Log.Printf("[getInstanceOfType] Creating new instance for type: %s", typeName)
		*value = factoryMethod()
		return *value
	}
	Log.Printf("[getInstanceOfType] Returning existing instance of type: %s (pointer: %p)", typeName, any(*value))

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
