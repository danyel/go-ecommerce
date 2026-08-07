package commonRepository

import (
	Uuid "github.com/google/uuid"
	Database "gorm.io/gorm"
)

type CrudRepository[T any] interface {
	FindAll(searchCriteria SearchCriteria) []*T
	FetchAll() []*T
	FindById(id Uuid.UUID, preloads ...string) (*T, error)
	Create(p *T) error
	Update(p *T) error
	Delete(id Uuid.UUID) error
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
	db *Database.DB
}

func (r *crudRepository[T]) FetchAll() []*T {
	var results []*T
	query := r.db.Model(new(T))
	query.Find(&results)
	return results
}

func (r *crudRepository[T]) FindAll(searchCriteria SearchCriteria) []*T {
	var results []*T
	query := r.db.Model(new(T))

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
func (r *crudRepository[T]) FindById(id Uuid.UUID, preloads ...string) (*T, error) {
	var model T
	var result *Database.DB
	query := r.db.Model(&model)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if result = query.First(&model, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return &model, nil
}
func (r *crudRepository[T]) Create(p *T) error {
	return r.db.Create(p).Error
}

func (r *crudRepository[T]) Update(p *T) error {
	return r.db.Save(p).Error
}

func (r *crudRepository[T]) Delete(id Uuid.UUID) error {
	var result *Database.DB
	if result = r.db.Delete(new(T), id); result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Database.ErrRecordNotFound
	}
	return nil
}

func (r *crudRepository[T]) Paginate(criteria SearchCriteria) ([]T, int64) {
	var results []T
	var total int64

	base := r.db.Model(new(T))

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

func (r *crudRepository[T]) AssocAppend(parent *T, assoc string, values any) error {
	return r.db.Session(&Database.Session{FullSaveAssociations: false}).Model(parent).Association(assoc).Append(values)
}

func (r *crudRepository[T]) AssocReplace(parent *T, assoc string, values any) error {
	return r.db.Model(parent).Association(assoc).Replace(values)
}

func (r *crudRepository[T]) AssocDelete(parent *T, assoc string, values any) error {
	return r.db.Model(parent).Association(assoc).Delete(values)
}

func (r *crudRepository[T]) AssocClear(parent *T, assoc string) error {
	return r.db.Model(parent).Association(assoc).Clear()
}

func (r *crudRepository[T]) AssocCount(parent *T, assoc string) (int64, error) {
	return r.db.Model(parent).Association(assoc).Count(), nil
}

func NewCrudRepository[T any](db *Database.DB) CrudRepository[T] {
	return &crudRepository[T]{db: db}
}
