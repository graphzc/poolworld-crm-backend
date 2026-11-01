package dto

import "time"

type CreateCustomerRequest struct {
	Name         string   `json:"name" validate:"required"`
	Address      string   `json:"address"`
	TaxID        string   `json:"tax_id"`
	CustomerType string   `json:"customer_type" validate:"required,oneof=Tax Non-Tax"`
	Services     []string `json:"services" validate:"required"`
}

type UpdateCustomerRequest struct {
	Name         string   `json:"name" validate:"required"`
	Address      string   `json:"address"`
	TaxID        string   `json:"tax_id"`
	CustomerType string   `json:"customer_type" validate:"required,oneof=Tax Non-Tax"`
	Services     []string `json:"services" validate:"required"`
}

type CustomerResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	TaxID        string    `json:"tax_id,omitempty"`
	CustomerType string    `json:"customer_type"`
	Services     []string  `json:"services"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CustomerQueryParams struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Search       string `query:"search"`
	SortBy       string `query:"sortBy"`
	SortOrder    string `query:"sortOrder"`
	CustomerType string `query:"customerType"`
}

type PaginatedCustomersResponse struct {
	Data       []*CustomerResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"totalPages"`
}
