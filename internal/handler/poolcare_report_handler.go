package handler

import (
	"net/http"

	"poolworld-pattaya-backend/internal/dto"
	"poolworld-pattaya-backend/internal/entity"
	"poolworld-pattaya-backend/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type PoolCareReportHandler struct {
	usecase   *usecase.PoolCareReportUsecase
	validator *validator.Validate
}

func NewPoolCareReportHandler(usecase *usecase.PoolCareReportUsecase) *PoolCareReportHandler {
	return &PoolCareReportHandler{
		usecase:   usecase,
		validator: validator.New(),
	}
}

func (h *PoolCareReportHandler) RegisterRoutes(e *echo.Group) {
	reports := e.Group("/poolcare-reports")
	reports.POST("", h.CreateReport)
	reports.POST("/bulk", h.BulkCreateReports)
	reports.GET("", h.GetAllReports)
	reports.GET("/:id", h.GetReportByID)
	reports.DELETE("/:id", h.DeleteReport)
}

func (h *PoolCareReportHandler) CreateReport(c *echo.Context) error {
	var req dto.CreatePoolCareReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	report, err := h.usecase.CreateReport(c.Request().Context(), req.PoolCareID, req.Year, req.Month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, h.entityToResponse(report))
}

func (h *PoolCareReportHandler) BulkCreateReports(c *echo.Context) error {
	var req dto.BulkCreatePoolCareReportRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	created, failed, errors := h.usecase.BulkCreateReports(c.Request().Context(), req.Year, req.Month)

	return c.JSON(http.StatusOK, dto.BulkCreatePoolCareReportResponse{
		Created: created,
		Failed:  failed,
		Errors:  errors,
	})
}

func (h *PoolCareReportHandler) GetAllReports(c *echo.Context) error {
	params := new(dto.PoolCareReportQueryParams)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	reports, total, err := h.usecase.GetReportsWithFilters(
		c.Request().Context(),
		params.PoolCareID,
		params.Year,
		params.Month,
		parseDays(params.Days),
		params.Page,
		params.Limit,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	responses := make([]dto.PoolCareReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = *h.entityToResponse(&report)
	}

	return c.JSON(http.StatusOK, dto.PaginatedPoolCareReportsResponse{
		Data:  responses,
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	})
}

func (h *PoolCareReportHandler) GetReportByID(c *echo.Context) error {
	id := c.Param("id")

	report, err := h.usecase.GetReportByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "report not found"})
	}

	return c.JSON(http.StatusOK, h.entityToResponse(report))
}

func (h *PoolCareReportHandler) DeleteReport(c *echo.Context) error {
	id := c.Param("id")

	if err := h.usecase.DeleteReport(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *PoolCareReportHandler) entityToResponse(e *entity.PoolCareReport) *dto.PoolCareReportResponse {
	daysOfService := make([]string, len(e.DaysOfService))
	for i, day := range e.DaysOfService {
		daysOfService[i] = string(day)
	}

	return &dto.PoolCareReportResponse{
		ID:               e.ID,
		PoolCareID:       e.PoolCareID,
		Year:             e.Year,
		Month:            e.Month,
		Name:             e.Name,
		Address:          e.Address,
		DaysOfService:    daysOfService,
		PoolSystem:       string(e.PoolSystem),
		ChemicalIncluded: e.ChemicalIncluded,
		ServiceDates:     e.ServiceDates,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}
