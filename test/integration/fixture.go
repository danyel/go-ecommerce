package integration

import (
	Log "log"

	CommonRepository "github.com/danyel/ecommerce/internal/common/repository"
)

type Fixture[T any] struct {
	repository CommonRepository.CrudRepository[T]
}

func (f Fixture[T]) Insert(t *T) {
	e := f.repository.Create(t)
	if e != nil {
		Log.Fatal(e)
	}
}

func Database[T any](repository CommonRepository.CrudRepository[T]) Fixture[T] {
	return Fixture[T]{repository}
}
