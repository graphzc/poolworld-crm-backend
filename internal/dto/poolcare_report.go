package dto

import "time"

type CreatePoolCareReportRequest struct {
	PoolCareID string `json:"poolcare_id" validate:"required"`
	Year       int    `json:"year" validate:"required,min=2020,max=2100"`
	Month      int    `json:"month" validate:"required,min=1,max=12"`
}

type BulkCreatePoolCareReportRequest struct {
	Year  int `json:"year" validate:"required,min=2020,max=2100"`
	Month int `json:"month" validate:"required,min=1,max=12"`
}

type BulkCreatePoolCareReportResponse struct {
	Created int      `json:"created"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

type PoolCareReportResponse struct {
	ID               string    `json:"id"`
	PoolCareID       string    `json:"poolcare_id"`
	Year             int       `json:"year"`
	Month            int       `json:"month"`
	Name             string    `json:"name"`
	Address          string    `json:"address"`
	DaysOfService    []string  `json:"days_of_service"`
	PoolSystem       string    `json:"pool_system"`
	ChemicalIncluded bool      `json:"chemical_included"`
	ServiceDates     []int     `json:"service_dates"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PoolCareReportQueryParams struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	PoolCareID string `query:"poolcare_id"`
	Year       int    `query:"year"`
	Month      int    `query:"month"`
	Days       string `query:"days"` // comma-separated, e.g. "MON,WED"
}

type PaginatedPoolCareReportsResponse struct {
	Data  []PoolCareReportResponse `json:"data"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
}
