package dto

import "time"

type CreatePoolCareRequest struct {
	CustomerID       string   `json:"customer_id" validate:"required"`
	Name             string   `json:"name" validate:"required"`
	Address          string   `json:"address"`
	DaysOfService    []string `json:"days_of_service" validate:"required,min=1"`
	PoolSystem       string   `json:"pool_system" validate:"required,oneof=CHLORINE SALT"`
	ChemicalIncluded bool     `json:"chemical_included"`
	Price            float64  `json:"price" validate:"min=0"`
}

type UpdatePoolCareRequest struct {
	Name             string   `json:"name" validate:"required"`
	Address          string   `json:"address"`
	DaysOfService    []string `json:"days_of_service" validate:"required,min=1"`
	PoolSystem       string   `json:"pool_system" validate:"required,oneof=CHLORINE SALT"`
	ChemicalIncluded bool     `json:"chemical_included"`
	Price            float64  `json:"price" validate:"min=0"`
}

type PoolCareResponse struct {
	ID               string    `json:"id"`
	CustomerID       string    `json:"customer_id"`
	CustomerName     string    `json:"customer_name,omitempty"`
	Name             string    `json:"name"`
	Address          string    `json:"address"`
	DaysOfService    []string  `json:"days_of_service"`
	PoolSystem       string    `json:"pool_system"`
	ChemicalIncluded bool      `json:"chemical_included"`
	Price            float64   `json:"price"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PoolCareQueryParams struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Search     string `query:"search"`
	SortBy     string `query:"sortBy"`
	SortOrder  string `query:"sortOrder"`
	CustomerID string `query:"customerId"`
	PoolSystem string `query:"poolSystem"`
	Days       string `query:"days"` // comma-separated, e.g. "MON,WED"
}

type PaginatedPoolCaresResponse struct {
	Data       []PoolCareResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}
