package integration

import (
	Log "log"

	Repository "github.com/danyel/ecommerce/internal/common/repository"
)

type Fixture[T any] struct {
	repository Repository.CrudRepository[T]
}

func (f Fixture[T]) Insert(t *T) {
	e := f.repository.Create(t)
	if e != nil {
		Log.Fatal(e)
	}
}

func Database[T any](repository Repository.CrudRepository[T]) Fixture[T] {
	return Fixture[T]{repository}
}
