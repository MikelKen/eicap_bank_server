package pagination

import "gorm.io/gorm"

type Scope func(*gorm.DB) *gorm.DB

func GormPaginate[T any](db *gorm.DB, params Params, scopes ...Scope) ([]T, int64, error) {
	query := db

	for _, scope := range scopes {
		query = scope(query)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PerPage
	var results []T
	if err := query.Offset(offset).Limit(params.PerPage).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
