package entity

import "time"

type CustomerType string

const (
	CustomerTypeTax    CustomerType = "Tax"
	CustomerTypeNonTax CustomerType = "Non-Tax"
)

type ServiceType string

const (
	ServiceTypePoolCare ServiceType = "POOL_CARE"
)

type Customer struct {
	ID           string
	Name         string
	Address      string
	TaxID        string
	CustomerType CustomerType
	Services     []ServiceType
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
