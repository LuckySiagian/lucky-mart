package utils

import (
	"fmt"
	"strings"

	"backend-sembako/internal/dto"

	"gorm.io/gorm"
)

// ApplyPagination applies ordering, offset and limit. Sorting is restricted to
// the columns listed in allowedSort (plus created_at) so a user-supplied
// sort_by value can never be injected into the SQL ORDER BY clause.
// Searching is left to each repository because searchable columns differ per table.
func ApplyPagination(db *gorm.DB, req dto.PaginationRequest, allowedSort ...string) *gorm.DB {
	req.SetDefaults()

	sortBy := "created_at"
	for _, col := range allowedSort {
		if req.SortBy == col {
			sortBy = col
			break
		}
	}

	sortDir := "DESC"
	if strings.EqualFold(req.SortDir, "ASC") {
		sortDir = "ASC"
	}

	db = db.Order(fmt.Sprintf("%s %s", sortBy, sortDir))
	return db.Offset(req.GetOffset()).Limit(req.Limit)
}

// SearchTerm wraps a search string for use in an ILIKE clause.
func SearchTerm(s string) string {
	return "%" + strings.TrimSpace(s) + "%"
}
