package repository

import (
	"context"
	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) error
	FindAll(ctx context.Context) ([]*entity.Customer, error)
	FindWithFilters(ctx context.Context, filter FilterOptions) ([]*entity.Customer, int64, error)
	FindByID(ctx context.Context, id string) (*entity.Customer, error)
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id string) error
}

type FilterOptions struct {
	Page         int
	Limit        int
	Search       string
	SortBy       string
	SortOrder    string
	CustomerType string
}

type customerRepository struct {
	collection *mongo.Collection
}

func NewCustomerRepository(db *mongo.Database) CustomerRepository {
	return &customerRepository{
		collection: db.Collection("customers"),
	}
}

func (r *customerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	services := make([]string, len(customer.Services))
	for i, s := range customer.Services {
		services[i] = string(s)
	}

	customerModel := &model.Customer{
		Name:         customer.Name,
		Address:      customer.Address,
		TaxID:        customer.TaxID,
		CustomerType: string(customer.CustomerType),
		Services:     services,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, customerModel)
	if err != nil {
		return err
	}

	customer.ID = result.InsertedID.(bson.ObjectID).Hex()
	customer.CreatedAt = customerModel.CreatedAt
	customer.UpdatedAt = customerModel.UpdatedAt
	return nil
}

func (r *customerRepository) FindAll(ctx context.Context) ([]*entity.Customer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []*entity.Customer
	for cursor.Next(ctx) {
		var customerModel model.Customer
		if err := cursor.Decode(&customerModel); err != nil {
			return nil, err
		}

		services := make([]entity.ServiceType, len(customerModel.Services))
		for i, s := range customerModel.Services {
			services[i] = entity.ServiceType(s)
		}

		customers = append(customers, &entity.Customer{
			ID:           customerModel.ID.Hex(),
			Name:         customerModel.Name,
			Address:      customerModel.Address,
			TaxID:        customerModel.TaxID,
			CustomerType: entity.CustomerType(customerModel.CustomerType),
			Services:     services,
			CreatedAt:    customerModel.CreatedAt,
			UpdatedAt:    customerModel.UpdatedAt,
		})
	}

	return customers, nil
}

func (r *customerRepository) FindWithFilters(ctx context.Context, filter FilterOptions) ([]*entity.Customer, int64, error) {
	// Build filter query
	query := bson.M{}

	// Search filter
	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"address": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"tax_id": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Customer type filter
	if filter.CustomerType != "" {
		query["customer_type"] = filter.CustomerType
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Build options
	opts := options.Find()

	// Pagination
	if filter.Page > 0 && filter.Limit > 0 {
		skip := int64((filter.Page - 1) * filter.Limit)
		opts.SetSkip(skip)
		opts.SetLimit(int64(filter.Limit))
	}

	// Sorting
	if filter.SortBy != "" {
		sortOrder := 1
		if filter.SortOrder == "desc" {
			sortOrder = -1
		}
		opts.SetSort(bson.D{{Key: filter.SortBy, Value: sortOrder}})
	} else {
		// Default sort by created_at desc
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var customers []*entity.Customer
	for cursor.Next(ctx) {
		var customerModel model.Customer
		if err := cursor.Decode(&customerModel); err != nil {
			return nil, 0, err
		}

		services := make([]entity.ServiceType, len(customerModel.Services))
		for i, s := range customerModel.Services {
			services[i] = entity.ServiceType(s)
		}

		customers = append(customers, &entity.Customer{
			ID:           customerModel.ID.Hex(),
			Name:         customerModel.Name,
			Address:      customerModel.Address,
			TaxID:        customerModel.TaxID,
			CustomerType: entity.CustomerType(customerModel.CustomerType),
			Services:     services,
			CreatedAt:    customerModel.CreatedAt,
			UpdatedAt:    customerModel.UpdatedAt,
		})
	}

	return customers, total, nil
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*entity.Customer, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var customerModel model.Customer
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&customerModel)
	if err != nil {
		return nil, err
	}

	services := make([]entity.ServiceType, len(customerModel.Services))
	for i, s := range customerModel.Services {
		services[i] = entity.ServiceType(s)
	}

	return &entity.Customer{
		ID:           customerModel.ID.Hex(),
		Name:         customerModel.Name,
		Address:      customerModel.Address,
		TaxID:        customerModel.TaxID,
		CustomerType: entity.CustomerType(customerModel.CustomerType),
		Services:     services,
		CreatedAt:    customerModel.CreatedAt,
		UpdatedAt:    customerModel.UpdatedAt,
	}, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	objectID, err := bson.ObjectIDFromHex(customer.ID)
	if err != nil {
		return err
	}

	services := make([]string, len(customer.Services))
	for i, s := range customer.Services {
		services[i] = string(s)
	}

	update := bson.M{
		"$set": bson.M{
			"name":          customer.Name,
			"address":       customer.Address,
			"tax_id":        customer.TaxID,
			"customer_type": string(customer.CustomerType),
			"services":      services,
			"updated_at":    time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
