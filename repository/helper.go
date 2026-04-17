package repository

import (
	"forgeflow-api/dto"

	"gorm.io/gorm"
)

func Paginate(db *gorm.DB, q dto.PaginateQuery) (*gorm.DB, int64, error) {
	var total int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "created_at"
	if q.SortBy != "" {
		sortBy = q.SortBy
	}

	orderBy := dto.OrderByDesc
	if q.OrderBy == dto.OrderByAsc {
		orderBy = dto.OrderByAsc
	}

	return db.
		Order(sortBy + " " + orderBy).
		Limit(q.Limit).
		Offset(q.Offset()), total, nil
}
