package commonRepository

import (
	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

type CrudRepository[T any] interface {
	FindAll(searchCriteria SearchCriteria) []*T
	FetchAll() []*T
	FindById(ID Uuid.UUID, preloads ...string) (*T, error)
	Create(model *T) error
	Update(model *T) error
	Delete(ID Uuid.UUID) error
	Paginate(criteria SearchCriteria) ([]T, int64)
	AssocAppend(parent *T, assoc string, values any) error
	AssocReplace(parent *T, assoc string, values any) error
	AssocDelete(parent *T, assoc string, values any) error
	AssocClear(parent *T, assoc string) error
	AssocCount(parent *T, assoc string) (int64, error)
}

type WhereClause struct {
	Query  string
	Params []any
}

type SearchCriteria struct {
	WhereClause WhereClause
	Limit       *int
	Offset      *int
	OrderBy     *string
	Preloads    []string
}

type crudRepository[T any] struct {
	databaseConnection *Database.DB
}

func (crudRepository *crudRepository[T]) FetchAll() []*T {
	var results []*T
	query := crudRepository.databaseConnection.Model(new(T))
	query.Find(&results)
	return results
}

func (crudRepository *crudRepository[T]) FindAll(searchCriteria SearchCriteria) []*T {
	var results []*T
	query := crudRepository.databaseConnection.Model(new(T))

	if searchCriteria.WhereClause.Query != "" {
		query = query.Where(
			searchCriteria.WhereClause.Query,
			searchCriteria.WhereClause.Params...,
		)
	}

	if searchCriteria.OrderBy != nil {
		query = query.Order(*searchCriteria.OrderBy)
	}

	if searchCriteria.Limit != nil {
		query = query.Limit(*searchCriteria.Limit)
	}

	if searchCriteria.Offset != nil {
		query = query.Offset(*searchCriteria.Offset)
	}

	// TODO only preload the most nearest children
	for _, preload := range searchCriteria.Preloads {
		query = query.Preload(preload)
	}

	query.Find(&results)
	return results
}
func (crudRepository *crudRepository[T]) FindById(ID Uuid.UUID, preloads ...string) (*T, error) {
	var model T
	var result *Database.DB
	query := crudRepository.databaseConnection.Model(&model)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if result = query.First(&model, "id = ?", ID); result.Error != nil {
		return nil, result.Error
	}

	return &model, nil
}
func (crudRepository *crudRepository[T]) Create(p *T) error {
	return crudRepository.databaseConnection.Create(p).Error
}

func (crudRepository *crudRepository[T]) Update(p *T) error {
	return crudRepository.databaseConnection.Save(p).Error
}

func (crudRepository *crudRepository[T]) Delete(ID Uuid.UUID) error {
	var result *Database.DB
	if result = crudRepository.databaseConnection.Delete(new(T), ID); result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Database.ErrRecordNotFound
	}
	return nil
}

func (crudRepository *crudRepository[T]) Paginate(criteria SearchCriteria) ([]T, int64) {
	var results []T
	var total int64

	base := crudRepository.databaseConnection.Model(new(T))

	// Count total
	if criteria.WhereClause.Query != "" {
		base = base.Where(criteria.WhereClause.Query, criteria.WhereClause.Params...)
	}
	base.Count(&total)

	// Apply preloads & pagination
	if criteria.WhereClause.Query != "" {
		base = base.Where(criteria.WhereClause.Query, criteria.WhereClause.Params...)
	}
	for _, preload := range criteria.Preloads {
		base = base.Preload(preload)
	}

	if criteria.OrderBy != nil {
		base = base.Order(*criteria.OrderBy)
	}
	if criteria.Limit != nil {
		base = base.Limit(*criteria.Limit)
	}
	if criteria.Offset != nil {
		base = base.Offset(*criteria.Offset)
	}

	base.Find(&results)
	return results, total
}

func (crudRepository *crudRepository[T]) AssocAppend(parent *T, assoc string, values any) error {
	return crudRepository.databaseConnection.Session(&Database.Session{FullSaveAssociations: false}).Model(parent).Association(assoc).Append(values)
}

func (crudRepository *crudRepository[T]) AssocReplace(parent *T, assoc string, values any) error {
	return crudRepository.databaseConnection.Model(parent).Association(assoc).Replace(values)
}

func (crudRepository *crudRepository[T]) AssocDelete(parent *T, assoc string, values any) error {
	return crudRepository.databaseConnection.Model(parent).Association(assoc).Delete(values)
}

func (crudRepository *crudRepository[T]) AssocClear(parent *T, assoc string) error {
	return crudRepository.databaseConnection.Model(parent).Association(assoc).Clear()
}

func (crudRepository *crudRepository[T]) AssocCount(parent *T, assoc string) (int64, error) {
	return crudRepository.databaseConnection.Model(parent).Association(assoc).Count(), nil
}

func NewCrudRepository[T any](databaseConnection *Database.DB) CrudRepository[T] {
	return &crudRepository[T]{databaseConnection: databaseConnection}
}
