package entity

import "time"

type DayOfWeek string

const (
	DayOfWeekSunday    DayOfWeek = "SUN"
	DayOfWeekMonday    DayOfWeek = "MON"
	DayOfWeekTuesday   DayOfWeek = "TUE"
	DayOfWeekWednesday DayOfWeek = "WED"
	DayOfWeekThursday  DayOfWeek = "THU"
	DayOfWeekFriday    DayOfWeek = "FRI"
	DayOfWeekSaturday  DayOfWeek = "SAT"
)

type PoolSystem string

const (
	PoolSystemChlorine PoolSystem = "CHLORINE"
	PoolSystemSalt     PoolSystem = "SALT"
)

type PoolCare struct {
	ID               string
	CustomerID       string
	Name             string
	Address          string
	DaysOfService    []DayOfWeek
	PoolSystem       PoolSystem
	ChemicalIncluded bool
	Price            float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
