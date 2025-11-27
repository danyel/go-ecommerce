package integration

import (
	"log"

	commonRepository "github.com/danyel/ecommerce/internal/common/repository"
)

type Fixture[T any] struct {
	repository commonRepository.CrudRepository[T]
}

func (f Fixture[T]) Insert(t *T) {
	e := f.repository.Create(t)
	if e != nil {
		log.Fatal(e)
	}
}

func Database[T any](repository commonRepository.CrudRepository[T]) Fixture[T] {
	return Fixture[T]{repository}
}
