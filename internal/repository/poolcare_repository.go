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

type PoolCareRepository interface {
	Create(ctx context.Context, poolCare *entity.PoolCare) error
	FindAll(ctx context.Context) ([]*entity.PoolCare, error)
	FindWithFilters(ctx context.Context, filter PoolCareFilterOptions) ([]*entity.PoolCare, int64, error)
	FindByID(ctx context.Context, id string) (*entity.PoolCare, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*entity.PoolCare, error)
	Update(ctx context.Context, poolCare *entity.PoolCare) error
	Delete(ctx context.Context, id string) error
	DeleteByCustomerID(ctx context.Context, customerID string) error
}

type PoolCareFilterOptions struct {
	Page       int
	Limit      int
	Search     string
	SortBy     string
	SortOrder  string
	CustomerID string
	PoolSystem string
	Days       []string
}

type poolCareRepository struct {
	collection *mongo.Collection
}

func NewPoolCareRepository(db *mongo.Database) PoolCareRepository {
	return &poolCareRepository{
		collection: db.Collection("poolcares"),
	}
}

func (r *poolCareRepository) Create(ctx context.Context, poolCare *entity.PoolCare) error {
	customerObjID, err := bson.ObjectIDFromHex(poolCare.CustomerID)
	if err != nil {
		return err
	}

	daysOfService := make([]string, len(poolCare.DaysOfService))
	for i, day := range poolCare.DaysOfService {
		daysOfService[i] = string(day)
	}

	poolCareModel := &model.PoolCare{
		CustomerID:       customerObjID,
		Name:             poolCare.Name,
		Address:          poolCare.Address,
		DaysOfService:    daysOfService,
		PoolSystem:       string(poolCare.PoolSystem),
		ChemicalIncluded: poolCare.ChemicalIncluded,
		Price:            poolCare.Price,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	result, err := r.collection.InsertOne(ctx, poolCareModel)
	if err != nil {
		return err
	}

	poolCare.ID = result.InsertedID.(bson.ObjectID).Hex()
	poolCare.CreatedAt = poolCareModel.CreatedAt
	poolCare.UpdatedAt = poolCareModel.UpdatedAt

	return nil
}

func (r *poolCareRepository) FindAll(ctx context.Context) ([]*entity.PoolCare, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var poolCares []*entity.PoolCare
	for cursor.Next(ctx) {
		var poolCareModel model.PoolCare
		if err := cursor.Decode(&poolCareModel); err != nil {
			return nil, err
		}

		daysOfService := make([]entity.DayOfWeek, len(poolCareModel.DaysOfService))
		for i, day := range poolCareModel.DaysOfService {
			daysOfService[i] = entity.DayOfWeek(day)
		}

		poolCares = append(poolCares, &entity.PoolCare{
			ID:               poolCareModel.ID.Hex(),
			CustomerID:       poolCareModel.CustomerID.Hex(),
			Name:             poolCareModel.Name,
			Address:          poolCareModel.Address,
			DaysOfService:    daysOfService,
			PoolSystem:       entity.PoolSystem(poolCareModel.PoolSystem),
			ChemicalIncluded: poolCareModel.ChemicalIncluded,
			Price:            poolCareModel.Price,
			CreatedAt:        poolCareModel.CreatedAt,
			UpdatedAt:        poolCareModel.UpdatedAt,
		})
	}

	return poolCares, nil
}

func (r *poolCareRepository) FindWithFilters(ctx context.Context, filter PoolCareFilterOptions) ([]*entity.PoolCare, int64, error) {
	query := bson.M{}

	// Search filter
	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"address": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Customer ID filter
	if filter.CustomerID != "" {
		customerObjID, err := bson.ObjectIDFromHex(filter.CustomerID)
		if err == nil {
			query["customer_id"] = customerObjID
		}
	}

	// Pool system filter
	if filter.PoolSystem != "" {
		query["pool_system"] = filter.PoolSystem
	}

	// Days of service filter (matches any of the given days)
	if len(filter.Days) > 0 {
		query["days_of_service"] = bson.M{"$in": filter.Days}
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
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var poolCares []*entity.PoolCare
	for cursor.Next(ctx) {
		var poolCareModel model.PoolCare
		if err := cursor.Decode(&poolCareModel); err != nil {
			return nil, 0, err
		}

		daysOfService := make([]entity.DayOfWeek, len(poolCareModel.DaysOfService))
		for i, day := range poolCareModel.DaysOfService {
			daysOfService[i] = entity.DayOfWeek(day)
		}

		poolCares = append(poolCares, &entity.PoolCare{
			ID:               poolCareModel.ID.Hex(),
			CustomerID:       poolCareModel.CustomerID.Hex(),
			Name:             poolCareModel.Name,
			Address:          poolCareModel.Address,
			DaysOfService:    daysOfService,
			PoolSystem:       entity.PoolSystem(poolCareModel.PoolSystem),
			ChemicalIncluded: poolCareModel.ChemicalIncluded,
			Price:            poolCareModel.Price,
			CreatedAt:        poolCareModel.CreatedAt,
			UpdatedAt:        poolCareModel.UpdatedAt,
		})
	}

	return poolCares, total, nil
}

func (r *poolCareRepository) FindByID(ctx context.Context, id string) (*entity.PoolCare, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var poolCareModel model.PoolCare
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&poolCareModel); err != nil {
		return nil, err
	}

	daysOfService := make([]entity.DayOfWeek, len(poolCareModel.DaysOfService))
	for i, day := range poolCareModel.DaysOfService {
		daysOfService[i] = entity.DayOfWeek(day)
	}

	return &entity.PoolCare{
		ID:               poolCareModel.ID.Hex(),
		CustomerID:       poolCareModel.CustomerID.Hex(),
		Name:             poolCareModel.Name,
		Address:          poolCareModel.Address,
		DaysOfService:    daysOfService,
		PoolSystem:       entity.PoolSystem(poolCareModel.PoolSystem),
		ChemicalIncluded: poolCareModel.ChemicalIncluded,
		Price:            poolCareModel.Price,
		CreatedAt:        poolCareModel.CreatedAt,
		UpdatedAt:        poolCareModel.UpdatedAt,
	}, nil
}

func (r *poolCareRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*entity.PoolCare, error) {
	customerObjID, err := bson.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"customer_id": customerObjID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var poolCares []*entity.PoolCare
	for cursor.Next(ctx) {
		var poolCareModel model.PoolCare
		if err := cursor.Decode(&poolCareModel); err != nil {
			return nil, err
		}

		daysOfService := make([]entity.DayOfWeek, len(poolCareModel.DaysOfService))
		for i, day := range poolCareModel.DaysOfService {
			daysOfService[i] = entity.DayOfWeek(day)
		}

		poolCares = append(poolCares, &entity.PoolCare{
			ID:               poolCareModel.ID.Hex(),
			CustomerID:       poolCareModel.CustomerID.Hex(),
			Name:             poolCareModel.Name,
			Address:          poolCareModel.Address,
			DaysOfService:    daysOfService,
			PoolSystem:       entity.PoolSystem(poolCareModel.PoolSystem),
			ChemicalIncluded: poolCareModel.ChemicalIncluded,
			Price:            poolCareModel.Price,
			CreatedAt:        poolCareModel.CreatedAt,
			UpdatedAt:        poolCareModel.UpdatedAt,
		})
	}

	return poolCares, nil
}

func (r *poolCareRepository) Update(ctx context.Context, poolCare *entity.PoolCare) error {
	objID, err := bson.ObjectIDFromHex(poolCare.ID)
	if err != nil {
		return err
	}

	daysOfService := make([]string, len(poolCare.DaysOfService))
	for i, day := range poolCare.DaysOfService {
		daysOfService[i] = string(day)
	}

	update := bson.M{
		"$set": bson.M{
			"name":              poolCare.Name,
			"address":           poolCare.Address,
			"days_of_service":   daysOfService,
			"pool_system":       string(poolCare.PoolSystem),
			"chemical_included": poolCare.ChemicalIncluded,
			"price":             poolCare.Price,
			"updated_at":        time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *poolCareRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *poolCareRepository) DeleteByCustomerID(ctx context.Context, customerID string) error {
	customerObjID, err := bson.ObjectIDFromHex(customerID)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteMany(ctx, bson.M{"customer_id": customerObjID})
	return err
}
