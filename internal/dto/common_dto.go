package dto

import "backend-sembako/internal/constants"

type PaginationRequest struct {
	Page    int    `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	Search  string `form:"search" json:"search"`
	SortBy  string `form:"sort_by" json:"sort_by"`
	SortDir string `form:"sort_dir" json:"sort_dir" binding:"omitempty,oneof=ASC DESC asc desc"`
}

func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = constants.DefaultPageSize
	}
	if p.SortBy == "" {
		p.SortBy = constants.DefaultSortBy
	}
	if p.SortDir == "" {
		p.SortDir = constants.DefaultSortDir
	}
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

type MetaResponse struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Items interface{}  `json:"items"`
	Meta  MetaResponse `json:"meta"`
}

func NewPaginatedResponse(items interface{}, page, limit int, totalRecords int64) PaginatedResponse {
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}
	return PaginatedResponse{
		Items: items,
		Meta: MetaResponse{
			Page:         page,
			Limit:        limit,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
		},
	}
}
