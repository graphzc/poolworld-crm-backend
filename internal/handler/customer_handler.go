package handler

import (
	"net/http"
	"poolworld-pattaya-backend/internal/dto"
	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/repository"
	"poolworld-pattaya-backend/internal/usecase"

	"github.com/labstack/echo/v5"
)

type CustomerHandler struct {
	customerUsecase usecase.CustomerUsecase
}

func NewCustomerHandler(customerUsecase usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

func (h *CustomerHandler) CreateCustomer(c *echo.Context) error {
	var req dto.CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	services := make([]entity.ServiceType, len(req.Services))
	for i, s := range req.Services {
		services[i] = entity.ServiceType(s)
	}

	customer := &entity.Customer{
		Name:         req.Name,
		Address:      req.Address,
		TaxID:        req.TaxID,
		CustomerType: entity.CustomerType(req.CustomerType),
		Services:     services,
	}

	if err := h.customerUsecase.CreateCustomer(c.Request().Context(), customer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, h.toResponse(customer))
}

func (h *CustomerHandler) GetAllCustomers(c *echo.Context) error {
	// Parse query parameters
	var params dto.CustomerQueryParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Get customers with filters
	customers, total, err := h.customerUsecase.GetCustomersWithFilters(
		c.Request().Context(),
		repository.FilterOptions{
			Page:         params.Page,
			Limit:        params.Limit,
			Search:       params.Search,
			SortBy:       params.SortBy,
			SortOrder:    params.SortOrder,
			CustomerType: params.CustomerType,
		},
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responses := make([]*dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		responses[i] = h.toResponse(customer)
	}

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, dto.PaginatedCustomersResponse{
		Data:       responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	})
}

func (h *CustomerHandler) GetCustomerByID(c *echo.Context) error {
	id := c.Param("id")

	customer, err := h.customerUsecase.GetCustomerByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "customer not found")
	}

	return c.JSON(http.StatusOK, h.toResponse(customer))
}

func (h *CustomerHandler) UpdateCustomer(c *echo.Context) error {
	id := c.Param("id")

	var req dto.UpdateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	services := make([]entity.ServiceType, len(req.Services))
	for i, s := range req.Services {
		services[i] = entity.ServiceType(s)
	}

	customer := &entity.Customer{
		ID:           id,
		Name:         req.Name,
		Address:      req.Address,
		TaxID:        req.TaxID,
		CustomerType: entity.CustomerType(req.CustomerType),
		Services:     services,
	}

	if err := h.customerUsecase.UpdateCustomer(c.Request().Context(), customer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, h.toResponse(customer))
}

func (h *CustomerHandler) DeleteCustomer(c *echo.Context) error {
	id := c.Param("id")

	if err := h.customerUsecase.DeleteCustomer(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CustomerHandler) toResponse(customer *entity.Customer) *dto.CustomerResponse {
	services := make([]string, len(customer.Services))
	for i, s := range customer.Services {
		services[i] = string(s)
	}

	return &dto.CustomerResponse{
		ID:           customer.ID,
		Name:         customer.Name,
		Address:      customer.Address,
		TaxID:        customer.TaxID,
		CustomerType: string(customer.CustomerType),
		Services:     services,
		CreatedAt:    customer.CreatedAt,
		UpdatedAt:    customer.UpdatedAt,
	}
}
