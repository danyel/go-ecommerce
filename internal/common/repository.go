package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CrudRepository[T any] struct {
	DB *gorm.DB
}

type WhereClause struct {
	Query  string
	Params []interface{}
}

type SearchCriteria struct {
	WhereClause WhereClause
	Limit       *int
	Offset      *int
	OrderBy     *string
	Preloads    []string
}

func NewCrudRepository[T any](db *gorm.DB) CrudRepository[T] {
	return CrudRepository[T]{db}
}

func (r *CrudRepository[T]) FindAll(searchCriteria SearchCriteria) []*T {
	var results []*T
	query := r.DB.Model(new(T))

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

func (r *CrudRepository[T]) FindById(id uuid.UUID, preloads ...string) (*T, error) {
	var model T
	query := r.DB.Model(&model)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	result := query.First(&model, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model, nil
}

func (r *CrudRepository[T]) Create(p *T) error {
	return r.DB.Create(p).Error
}

func (r *CrudRepository[T]) Update(p *T) error {
	return r.DB.Save(p).Error
}

func (r *CrudRepository[T]) Delete(id uuid.UUID) error {
	result := r.DB.Delete(new(T), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *CrudRepository[T]) Paginate(criteria SearchCriteria) ([]T, int64) {
	var results []T
	var total int64

	base := r.DB.Model(new(T))

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

func (r *CrudRepository[T]) AssocAppend(parent *T, assoc string, values interface{}) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: false}).Debug().Model(parent).Association(assoc).Append(values)
}

func (r *CrudRepository[T]) AssocReplace(parent *T, assoc string, values interface{}) error {
	return r.DB.Model(parent).Association(assoc).Replace(values)
}

func (r *CrudRepository[T]) AssocDelete(parent *T, assoc string, values interface{}) error {
	return r.DB.Model(parent).Association(assoc).Delete(values)
}

func (r *CrudRepository[T]) AssocClear(parent *T, assoc string) error {
	return r.DB.Model(parent).Association(assoc).Clear()
}

func (r *CrudRepository[T]) AssocCount(parent *T, assoc string) (int64, error) {
	return r.DB.Model(parent).Association(assoc).Count(), nil
}
