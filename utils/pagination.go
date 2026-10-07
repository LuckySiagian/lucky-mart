package utils

import (
	"fmt"

	"backend-sembako/internal/dto"

	"gorm.io/gorm"
)

func ApplyPagination(db *gorm.DB, req dto.PaginationRequest) *gorm.DB {
	req.SetDefaults()

	// Apply Search Filter if provided
	if req.Search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", req.Search)
		db = db.Where("name ILIKE ? OR full_name ILIKE ? OR phone_number ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Apply Order By
	sortOrder := fmt.Sprintf("%s %s", req.SortBy, req.SortDir)
	db = db.Order(sortOrder)

	// Apply Offset & Limit
	return db.Offset(req.GetOffset()).Limit(req.Limit)
}
