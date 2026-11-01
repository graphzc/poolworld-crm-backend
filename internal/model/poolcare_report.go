package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PoolCareReport struct {
	ID               bson.ObjectID `bson:"_id,omitempty"`
	PoolCareID       bson.ObjectID `bson:"poolcare_id"`
	Year             int           `bson:"year"`
	Month            int           `bson:"month"`
	Name             string        `bson:"name"`
	Address          string        `bson:"address"`
	DaysOfService    []string      `bson:"days_of_service"`
	PoolSystem       string        `bson:"pool_system"`
	ChemicalIncluded bool          `bson:"chemical_included"`
	ServiceDates     []int         `bson:"service_dates"`
	CreatedAt        time.Time     `bson:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at"`
}
