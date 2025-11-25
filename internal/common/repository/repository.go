package commonRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CrudRepository[T any] struct {
	FindAll      func(searchCriteria SearchCriteria) []*T
	FindById     func(id uuid.UUID, preloads ...string) (*T, error)
	Create       func(p *T) error
	Update       func(p *T) error
	Delete       func(id uuid.UUID) error
	Paginate     func(criteria SearchCriteria) ([]T, int64)
	AssocAppend  func(parent *T, assoc string, values interface{}) error
	AssocReplace func(parent *T, assoc string, values interface{}) error
	AssocDelete  func(parent *T, assoc string, values interface{}) error
	AssocClear   func(parent *T, assoc string) error
	AssocCount   func(parent *T, assoc string) (int64, error)
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
	return CrudRepository[T]{
		FindAll: func(searchCriteria SearchCriteria) []*T {
			var results []*T
			query := db.Model(new(T))

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
		},
		FindById: func(id uuid.UUID, preloads ...string) (*T, error) {
			var model T
			query := db.Model(&model)

			for _, preload := range preloads {
				query = query.Preload(preload)
			}

			result := query.First(&model, "id = ?", id)
			if result.Error != nil {
				return nil, result.Error
			}

			return &model, nil
		},
		Create: func(p *T) error {
			return db.Create(p).Error
		},
		Update: func(p *T) error {
			return db.Save(p).Error
		},
		Delete: func(id uuid.UUID) error {
			result := db.Delete(new(T), id)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
			return nil
		},
		Paginate: func(criteria SearchCriteria) ([]T, int64) {
			var results []T
			var total int64

			base := db.Model(new(T))

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
		},
		AssocAppend: func(parent *T, assoc string, values interface{}) error {
			return db.Session(&gorm.Session{FullSaveAssociations: false}).Model(parent).Association(assoc).Append(values)
		},
		AssocReplace: func(parent *T, assoc string, values interface{}) error {
			return db.Model(parent).Association(assoc).Replace(values)
		},
		AssocDelete: func(parent *T, assoc string, values interface{}) error {
			return db.Model(parent).Association(assoc).Delete(values)
		},
		AssocClear: func(parent *T, assoc string) error {
			return db.Model(parent).Association(assoc).Clear()
		},
		AssocCount: func(parent *T, assoc string) (int64, error) {
			return db.Model(parent).Association(assoc).Count(), nil
		},
	}
}
