package usecase

import (
	"context"
	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/repository"
)

type PoolCareUsecase interface {
	CreatePoolCare(ctx context.Context, poolCare *entity.PoolCare) error
	GetAllPoolCares(ctx context.Context) ([]*entity.PoolCare, error)
	GetPoolCaresWithFilters(ctx context.Context, filter repository.PoolCareFilterOptions) ([]*entity.PoolCare, int64, error)
	GetPoolCareByID(ctx context.Context, id string) (*entity.PoolCare, error)
	GetPoolCaresByCustomerID(ctx context.Context, customerID string) ([]*entity.PoolCare, error)
	UpdatePoolCare(ctx context.Context, poolCare *entity.PoolCare) error
	DeletePoolCare(ctx context.Context, id string) error
}

type poolCareUsecase struct {
	poolCareRepo       repository.PoolCareRepository
	customerRepo       repository.CustomerRepository
	poolCareReportRepo repository.PoolCareReportRepository
}

func NewPoolCareUsecase(poolCareRepo repository.PoolCareRepository, customerRepo repository.CustomerRepository, poolCareReportRepo repository.PoolCareReportRepository) PoolCareUsecase {
	return &poolCareUsecase{
		poolCareRepo:       poolCareRepo,
		customerRepo:       customerRepo,
		poolCareReportRepo: poolCareReportRepo,
	}
}

func (u *poolCareUsecase) CreatePoolCare(ctx context.Context, poolCare *entity.PoolCare) error {
	// Verify customer exists
	_, err := u.customerRepo.FindByID(ctx, poolCare.CustomerID)
	if err != nil {
		return err
	}

	return u.poolCareRepo.Create(ctx, poolCare)
}

func (u *poolCareUsecase) GetAllPoolCares(ctx context.Context) ([]*entity.PoolCare, error) {
	return u.poolCareRepo.FindAll(ctx)
}

func (u *poolCareUsecase) GetPoolCaresWithFilters(ctx context.Context, filter repository.PoolCareFilterOptions) ([]*entity.PoolCare, int64, error) {
	return u.poolCareRepo.FindWithFilters(ctx, filter)
}

func (u *poolCareUsecase) GetPoolCareByID(ctx context.Context, id string) (*entity.PoolCare, error) {
	return u.poolCareRepo.FindByID(ctx, id)
}

func (u *poolCareUsecase) GetPoolCaresByCustomerID(ctx context.Context, customerID string) ([]*entity.PoolCare, error) {
	return u.poolCareRepo.FindByCustomerID(ctx, customerID)
}

func (u *poolCareUsecase) UpdatePoolCare(ctx context.Context, poolCare *entity.PoolCare) error {
	return u.poolCareRepo.Update(ctx, poolCare)
}

// DeletePoolCare cascades to the pool care's reports.
func (u *poolCareUsecase) DeletePoolCare(ctx context.Context, id string) error {
	if _, err := u.poolCareRepo.FindByID(ctx, id); err != nil {
		return err
	}

	if err := u.poolCareReportRepo.DeleteByPoolCareIDs(ctx, []string{id}); err != nil {
		return err
	}

	return u.poolCareRepo.Delete(ctx, id)
}
