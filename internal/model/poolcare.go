package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PoolCare struct {
	ID               bson.ObjectID `bson:"_id,omitempty"`
	CustomerID       bson.ObjectID `bson:"customer_id"`
	Name             string        `bson:"name"`
	Address          string        `bson:"address"`
	DaysOfService    []string      `bson:"days_of_service"`
	PoolSystem       string        `bson:"pool_system"`
	ChemicalIncluded bool          `bson:"chemical_included"`
	Price            float64       `bson:"price"`
	CreatedAt        time.Time     `bson:"created_at"`
	UpdatedAt        time.Time     `bson:"updated_at"`
}
