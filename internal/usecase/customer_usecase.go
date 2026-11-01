package usecase

import (
	"context"
	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/repository"
)

type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, customer *entity.Customer) error
	GetAllCustomers(ctx context.Context) ([]*entity.Customer, error)
	GetCustomersWithFilters(ctx context.Context, filter repository.FilterOptions) ([]*entity.Customer, int64, error)
	GetCustomerByID(ctx context.Context, id string) (*entity.Customer, error)
	UpdateCustomer(ctx context.Context, customer *entity.Customer) error
	DeleteCustomer(ctx context.Context, id string) error
}

type customerUsecase struct {
	customerRepo       repository.CustomerRepository
	poolCareRepo       repository.PoolCareRepository
	poolCareReportRepo repository.PoolCareReportRepository
}

func NewCustomerUsecase(customerRepo repository.CustomerRepository, poolCareRepo repository.PoolCareRepository, poolCareReportRepo repository.PoolCareReportRepository) CustomerUsecase {
	return &customerUsecase{
		customerRepo:       customerRepo,
		poolCareRepo:       poolCareRepo,
		poolCareReportRepo: poolCareReportRepo,
	}
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, customer *entity.Customer) error {
	return u.customerRepo.Create(ctx, customer)
}

func (u *customerUsecase) GetAllCustomers(ctx context.Context) ([]*entity.Customer, error) {
	return u.customerRepo.FindAll(ctx)
}

func (u *customerUsecase) GetCustomersWithFilters(ctx context.Context, filter repository.FilterOptions) ([]*entity.Customer, int64, error) {
	return u.customerRepo.FindWithFilters(ctx, filter)
}

func (u *customerUsecase) GetCustomerByID(ctx context.Context, id string) (*entity.Customer, error) {
	return u.customerRepo.FindByID(ctx, id)
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, customer *entity.Customer) error {
	return u.customerRepo.Update(ctx, customer)
}

// DeleteCustomer cascades to the customer's pool cares and their reports.
// Children are removed first so a partial failure never leaves orphans.
func (u *customerUsecase) DeleteCustomer(ctx context.Context, id string) error {
	poolCares, err := u.poolCareRepo.FindByCustomerID(ctx, id)
	if err != nil {
		return err
	}

	poolCareIDs := make([]string, len(poolCares))
	for i, poolCare := range poolCares {
		poolCareIDs[i] = poolCare.ID
	}

	if err := u.poolCareReportRepo.DeleteByPoolCareIDs(ctx, poolCareIDs); err != nil {
		return err
	}

	if err := u.poolCareRepo.DeleteByCustomerID(ctx, id); err != nil {
		return err
	}

	return u.customerRepo.Delete(ctx, id)
}
