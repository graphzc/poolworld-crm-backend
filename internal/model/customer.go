package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Customer struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	Name         string        `bson:"name"`
	Address      string        `bson:"address"`
	TaxID        string        `bson:"tax_id,omitempty"`
	CustomerType string        `bson:"customer_type"`
	Services     []string      `bson:"services"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}
