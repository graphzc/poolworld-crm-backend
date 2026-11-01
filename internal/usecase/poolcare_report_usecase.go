package usecase

import (
	"context"
	"errors"
	"time"

	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/repository"
)

type PoolCareReportUsecase struct {
	poolCareRepo       repository.PoolCareRepository
	poolCareReportRepo repository.PoolCareReportRepository
}

func NewPoolCareReportUsecase(poolCareRepo repository.PoolCareRepository, poolCareReportRepo repository.PoolCareReportRepository) *PoolCareReportUsecase {
	return &PoolCareReportUsecase{
		poolCareRepo:       poolCareRepo,
		poolCareReportRepo: poolCareReportRepo,
	}
}

// CreateReport upserts the report: a pool care has at most one report per year/month.
func (uc *PoolCareReportUsecase) CreateReport(ctx context.Context, poolCareID string, year, month int) (*entity.PoolCareReport, error) {
	poolCare, err := uc.poolCareRepo.FindByID(ctx, poolCareID)
	if err != nil {
		return nil, errors.New("pool care not found")
	}

	serviceDates := calculateServiceDates(year, month, poolCare.DaysOfService)

	report := &entity.PoolCareReport{
		PoolCareID:       poolCareID,
		Year:             year,
		Month:            month,
		Name:             poolCare.Name,
		Address:          poolCare.Address,
		DaysOfService:    poolCare.DaysOfService,
		PoolSystem:       poolCare.PoolSystem,
		ChemicalIncluded: poolCare.ChemicalIncluded,
		ServiceDates:     serviceDates,
	}

	if err := uc.poolCareReportRepo.Upsert(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (uc *PoolCareReportUsecase) GetReportByID(ctx context.Context, id string) (*entity.PoolCareReport, error) {
	return uc.poolCareReportRepo.FindByID(ctx, id)
}

func (uc *PoolCareReportUsecase) GetReportsWithFilters(ctx context.Context, poolCareID string, year, month int, days []string, page, limit int) ([]entity.PoolCareReport, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return uc.poolCareReportRepo.FindWithFilters(ctx, poolCareID, year, month, days, page, limit)
}

func (uc *PoolCareReportUsecase) DeleteReport(ctx context.Context, id string) error {
	return uc.poolCareReportRepo.Delete(ctx, id)
}

func (uc *PoolCareReportUsecase) BulkCreateReports(ctx context.Context, year, month int) (int, int, []string) {
	// Get all pool cares
	allPoolCares, err := uc.poolCareRepo.FindAll(ctx)
	if err != nil {
		return 0, 0, []string{"Failed to fetch pool cares: " + err.Error()}
	}

	created := 0
	failed := 0
	errors := []string{}

	for _, poolCare := range allPoolCares {
		serviceDates := calculateServiceDates(year, month, poolCare.DaysOfService)

		report := &entity.PoolCareReport{
			PoolCareID:       poolCare.ID,
			Year:             year,
			Month:            month,
			Name:             poolCare.Name,
			Address:          poolCare.Address,
			DaysOfService:    poolCare.DaysOfService,
			PoolSystem:       poolCare.PoolSystem,
			ChemicalIncluded: poolCare.ChemicalIncluded,
			ServiceDates:     serviceDates,
		}

		if err := uc.poolCareReportRepo.Upsert(ctx, report); err != nil {
			failed++
			errors = append(errors, "Failed to upsert report for "+poolCare.Name+": "+err.Error())
		} else {
			created++
		}
	}

	return created, failed, errors
}

func calculateServiceDates(year, month int, daysOfService []entity.DayOfWeek) []int {
	serviceDates := []int{}

	dayMap := map[entity.DayOfWeek]time.Weekday{
		entity.DayOfWeekSunday:    time.Sunday,
		entity.DayOfWeekMonday:    time.Monday,
		entity.DayOfWeekTuesday:   time.Tuesday,
		entity.DayOfWeekWednesday: time.Wednesday,
		entity.DayOfWeekThursday:  time.Thursday,
		entity.DayOfWeekFriday:    time.Friday,
		entity.DayOfWeekSaturday:  time.Saturday,
	}

	serviceWeekdays := make(map[time.Weekday]bool)
	for _, day := range daysOfService {
		if weekday, ok := dayMap[day]; ok {
			serviceWeekdays[weekday] = true
		}
	}

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	for day := 1; day <= lastDay.Day(); day++ {
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
		if serviceWeekdays[date.Weekday()] {
			serviceDates = append(serviceDates, day)
		}
	}

	return serviceDates
}
