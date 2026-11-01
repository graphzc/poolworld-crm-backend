package entity

import "time"

type PoolCareReport struct {
	ID               string      `json:"id"`
	PoolCareID       string      `json:"poolcare_id"`
	Year             int         `json:"year"`
	Month            int         `json:"month"`
	Name             string      `json:"name"`
	Address          string      `json:"address"`
	DaysOfService    []DayOfWeek `json:"days_of_service"`
	PoolSystem       PoolSystem  `json:"pool_system"`
	ChemicalIncluded bool        `json:"chemical_included"`
	ServiceDates     []int       `json:"service_dates"` // Days in month that match service schedule
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}
