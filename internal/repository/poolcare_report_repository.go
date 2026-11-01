package repository

import (
	"context"
	"time"

	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PoolCareReportRepository interface {
	EnsureIndexes(ctx context.Context) error
	Upsert(ctx context.Context, report *entity.PoolCareReport) error
	FindByID(ctx context.Context, id string) (*entity.PoolCareReport, error)
	FindWithFilters(ctx context.Context, poolCareID string, year, month int, days []string, page, limit int) ([]entity.PoolCareReport, int64, error)
	Delete(ctx context.Context, id string) error
	DeleteByPoolCareIDs(ctx context.Context, poolCareIDs []string) error
}

type poolCareReportRepository struct {
	collection *mongo.Collection
}

func NewPoolCareReportRepository(db *mongo.Database) PoolCareReportRepository {
	return &poolCareReportRepository{
		collection: db.Collection("poolcare_reports"),
	}
}

// EnsureIndexes guarantees at most one report per pool care per month.
func (r *poolCareReportRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "poolcare_id", Value: 1}, {Key: "year", Value: 1}, {Key: "month", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_poolcare_period"),
	})
	return err
}

// Upsert creates the report for the pool care's year/month, or overwrites the existing one.
func (r *poolCareReportRepository) Upsert(ctx context.Context, report *entity.PoolCareReport) error {
	poolCareID, err := bson.ObjectIDFromHex(report.PoolCareID)
	if err != nil {
		return err
	}

	daysOfService := make([]string, len(report.DaysOfService))
	for i, day := range report.DaysOfService {
		daysOfService[i] = string(day)
	}

	now := time.Now()
	filter := bson.M{
		"poolcare_id": poolCareID,
		"year":        report.Year,
		"month":       report.Month,
	}
	update := bson.M{
		"$set": bson.M{
			"name":              report.Name,
			"address":           report.Address,
			"days_of_service":   daysOfService,
			"pool_system":       string(report.PoolSystem),
			"chemical_included": report.ChemicalIncluded,
			"service_dates":     report.ServiceDates,
			"updated_at":        now,
		},
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var modelReport model.PoolCareReport
	if err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&modelReport); err != nil {
		return err
	}

	*report = *r.modelToEntity(&modelReport)
	return nil
}

func (r *poolCareReportRepository) FindByID(ctx context.Context, id string) (*entity.PoolCareReport, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var modelReport model.PoolCareReport
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&modelReport)
	if err != nil {
		return nil, err
	}

	return r.modelToEntity(&modelReport), nil
}

func (r *poolCareReportRepository) FindWithFilters(ctx context.Context, poolCareID string, year, month int, days []string, page, limit int) ([]entity.PoolCareReport, int64, error) {
	filter := bson.M{}

	if poolCareID != "" {
		poolCareOID, err := bson.ObjectIDFromHex(poolCareID)
		if err == nil {
			filter["poolcare_id"] = poolCareOID
		}
	}

	if year > 0 {
		filter["year"] = year
	}

	if month > 0 {
		filter["month"] = month
	}

	if len(days) > 0 {
		filter["days_of_service"] = bson.M{"$in": days}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "year", Value: -1}, {Key: "month", Value: -1}, {Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var modelReports []model.PoolCareReport
	if err = cursor.All(ctx, &modelReports); err != nil {
		return nil, 0, err
	}

	reports := make([]entity.PoolCareReport, len(modelReports))
	for i, modelReport := range modelReports {
		reports[i] = *r.modelToEntity(&modelReport)
	}

	return reports, total, nil
}

func (r *poolCareReportRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *poolCareReportRepository) DeleteByPoolCareIDs(ctx context.Context, poolCareIDs []string) error {
	if len(poolCareIDs) == 0 {
		return nil
	}

	poolCareOIDs := make([]bson.ObjectID, len(poolCareIDs))
	for i, id := range poolCareIDs {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		poolCareOIDs[i] = oid
	}

	_, err := r.collection.DeleteMany(ctx, bson.M{"poolcare_id": bson.M{"$in": poolCareOIDs}})
	return err
}

func (r *poolCareReportRepository) modelToEntity(m *model.PoolCareReport) *entity.PoolCareReport {
	daysOfService := make([]entity.DayOfWeek, len(m.DaysOfService))
	for i, day := range m.DaysOfService {
		daysOfService[i] = entity.DayOfWeek(day)
	}

	return &entity.PoolCareReport{
		ID:               m.ID.Hex(),
		PoolCareID:       m.PoolCareID.Hex(),
		Year:             m.Year,
		Month:            m.Month,
		Name:             m.Name,
		Address:          m.Address,
		DaysOfService:    daysOfService,
		PoolSystem:       entity.PoolSystem(m.PoolSystem),
		ChemicalIncluded: m.ChemicalIncluded,
		ServiceDates:     m.ServiceDates,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
