package handler

import (
	"net/http"
	"poolworld-pattaya-backend/internal/dto"
	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/repository"
	"poolworld-pattaya-backend/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type PoolCareHandler struct {
	poolCareUsecase usecase.PoolCareUsecase
	customerUsecase usecase.CustomerUsecase
	validator       *validator.Validate
}

func NewPoolCareHandler(poolCareUsecase usecase.PoolCareUsecase, customerUsecase usecase.CustomerUsecase) *PoolCareHandler {
	return &PoolCareHandler{
		poolCareUsecase: poolCareUsecase,
		customerUsecase: customerUsecase,
		validator:       validator.New(),
	}
}

func (h *PoolCareHandler) CreatePoolCare(c *echo.Context) error {
	var req dto.CreatePoolCareRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	daysOfService := make([]entity.DayOfWeek, len(req.DaysOfService))
	for i, day := range req.DaysOfService {
		daysOfService[i] = entity.DayOfWeek(day)
	}

	poolCare := &entity.PoolCare{
		CustomerID:       req.CustomerID,
		Name:             req.Name,
		Address:          req.Address,
		DaysOfService:    daysOfService,
		PoolSystem:       entity.PoolSystem(req.PoolSystem),
		ChemicalIncluded: req.ChemicalIncluded,
		Price:            req.Price,
	}

	if err := h.poolCareUsecase.CreatePoolCare(c.Request().Context(), poolCare); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	daysOfServiceStr := make([]string, len(poolCare.DaysOfService))
	for i, day := range poolCare.DaysOfService {
		daysOfServiceStr[i] = string(day)
	}

	response := dto.PoolCareResponse{
		ID:               poolCare.ID,
		CustomerID:       poolCare.CustomerID,
		Name:             poolCare.Name,
		Address:          poolCare.Address,
		DaysOfService:    daysOfServiceStr,
		PoolSystem:       string(poolCare.PoolSystem),
		ChemicalIncluded: poolCare.ChemicalIncluded,
		Price:            poolCare.Price,
		CreatedAt:        poolCare.CreatedAt,
		UpdatedAt:        poolCare.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *PoolCareHandler) GetAllPoolCares(c *echo.Context) error {
	var params dto.PoolCareQueryParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query parameters"})
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	filter := repository.PoolCareFilterOptions{
		Page:       params.Page,
		Limit:      params.Limit,
		Search:     params.Search,
		SortBy:     params.SortBy,
		SortOrder:  params.SortOrder,
		CustomerID: params.CustomerID,
		PoolSystem: params.PoolSystem,
		Days:       parseDays(params.Days),
	}

	poolCares, total, err := h.poolCareUsecase.GetPoolCaresWithFilters(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]dto.PoolCareResponse, len(poolCares))
	for i, poolCare := range poolCares {
		daysOfService := make([]string, len(poolCare.DaysOfService))
		for j, day := range poolCare.DaysOfService {
			daysOfService[j] = string(day)
		}

		customer, err := h.customerUsecase.GetCustomerByID(c.Request().Context(), poolCare.CustomerID)
		customerName := ""
		if err == nil && customer != nil {
			customerName = customer.Name
		}

		responses[i] = dto.PoolCareResponse{
			ID:               poolCare.ID,
			CustomerID:       poolCare.CustomerID,
			CustomerName:     customerName,
			Name:             poolCare.Name,
			Address:          poolCare.Address,
			DaysOfService:    daysOfService,
			PoolSystem:       string(poolCare.PoolSystem),
			ChemicalIncluded: poolCare.ChemicalIncluded,
			Price:            poolCare.Price,
			CreatedAt:        poolCare.CreatedAt,
			UpdatedAt:        poolCare.UpdatedAt,
		}
	}

	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}

	paginatedResponse := dto.PaginatedPoolCaresResponse{
		Data:       responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, paginatedResponse)
}

func (h *PoolCareHandler) GetPoolCareByID(c *echo.Context) error {
	id := c.Param("id")

	poolCare, err := h.poolCareUsecase.GetPoolCareByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Pool care not found"})
	}

	daysOfService := make([]string, len(poolCare.DaysOfService))
	for i, day := range poolCare.DaysOfService {
		daysOfService[i] = string(day)
	}

	customer, err := h.customerUsecase.GetCustomerByID(c.Request().Context(), poolCare.CustomerID)
	customerName := ""
	if err == nil && customer != nil {
		customerName = customer.Name
	}

	response := dto.PoolCareResponse{
		ID:               poolCare.ID,
		CustomerID:       poolCare.CustomerID,
		CustomerName:     customerName,
		Name:             poolCare.Name,
		Address:          poolCare.Address,
		DaysOfService:    daysOfService,
		PoolSystem:       string(poolCare.PoolSystem),
		ChemicalIncluded: poolCare.ChemicalIncluded,
		Price:            poolCare.Price,
		CreatedAt:        poolCare.CreatedAt,
		UpdatedAt:        poolCare.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PoolCareHandler) UpdatePoolCare(c *echo.Context) error {
	id := c.Param("id")

	var req dto.UpdatePoolCareRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	daysOfService := make([]entity.DayOfWeek, len(req.DaysOfService))
	for i, day := range req.DaysOfService {
		daysOfService[i] = entity.DayOfWeek(day)
	}

	poolCare := &entity.PoolCare{
		ID:               id,
		Name:             req.Name,
		Address:          req.Address,
		DaysOfService:    daysOfService,
		PoolSystem:       entity.PoolSystem(req.PoolSystem),
		ChemicalIncluded: req.ChemicalIncluded,
		Price:            req.Price,
	}

	if err := h.poolCareUsecase.UpdatePoolCare(c.Request().Context(), poolCare); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	updatedPoolCare, err := h.poolCareUsecase.GetPoolCareByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	daysOfServiceStr := make([]string, len(updatedPoolCare.DaysOfService))
	for i, day := range updatedPoolCare.DaysOfService {
		daysOfServiceStr[i] = string(day)
	}

	response := dto.PoolCareResponse{
		ID:               updatedPoolCare.ID,
		CustomerID:       updatedPoolCare.CustomerID,
		Name:             updatedPoolCare.Name,
		Address:          updatedPoolCare.Address,
		DaysOfService:    daysOfServiceStr,
		PoolSystem:       string(updatedPoolCare.PoolSystem),
		ChemicalIncluded: updatedPoolCare.ChemicalIncluded,
		Price:            updatedPoolCare.Price,
		CreatedAt:        updatedPoolCare.CreatedAt,
		UpdatedAt:        updatedPoolCare.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *PoolCareHandler) DeletePoolCare(c *echo.Context) error {
	id := c.Param("id")

	if err := h.poolCareUsecase.DeletePoolCare(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *PoolCareHandler) GetPoolCaresByCustomer(c *echo.Context) error {
	customerID := c.Param("customerId")

	poolCares, err := h.poolCareUsecase.GetPoolCaresByCustomerID(c.Request().Context(), customerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]dto.PoolCareResponse, len(poolCares))
	for i, poolCare := range poolCares {
		daysOfService := make([]string, len(poolCare.DaysOfService))
		for j, day := range poolCare.DaysOfService {
			daysOfService[j] = string(day)
		}

		responses[i] = dto.PoolCareResponse{
			ID:               poolCare.ID,
			CustomerID:       poolCare.CustomerID,
			Name:             poolCare.Name,
			Address:          poolCare.Address,
			DaysOfService:    daysOfService,
			PoolSystem:       string(poolCare.PoolSystem),
			ChemicalIncluded: poolCare.ChemicalIncluded,
			Price:            poolCare.Price,
			CreatedAt:        poolCare.CreatedAt,
			UpdatedAt:        poolCare.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, responses)
}

func (h *PoolCareHandler) RegisterRoutes(e *echo.Group) {
	poolCares := e.Group("/poolcares")
	poolCares.POST("", h.CreatePoolCare)
	poolCares.GET("", h.GetAllPoolCares)
	poolCares.GET("/:id", h.GetPoolCareByID)
	poolCares.PUT("/:id", h.UpdatePoolCare)
	poolCares.DELETE("/:id", h.DeletePoolCare)
	poolCares.GET("/customer/:customerId", h.GetPoolCaresByCustomer)
}
