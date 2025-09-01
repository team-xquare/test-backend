package application

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/team-xquare/deployment-platform/internal/pkg/middleware"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	applications := r.Group("/applications")
	applications.Use(middleware.Auth())
	{
		applications.GET("/:id", h.GetApplication)
		applications.PUT("/:id", h.UpdateApplication)
		applications.DELETE("/:id", h.DeleteApplication)
	}

	// Project-specific application routes
	projects := r.Group("/projects/:id/applications")
	projects.Use(middleware.Auth())
	{
		projects.GET("", h.GetApplicationsByProject)
		projects.POST("", h.CreateApplication)
	}
}

func (h *Handler) CreateApplication(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	_ = c.GetUint("user_id")

	app, err := h.service.CreateApplication(c.Request.Context(), uint(projectID), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, app)
}

func (h *Handler) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid application ID"))
		return
	}

	_ = c.GetUint("user_id")

	app, err := h.service.GetApplication(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, app)
}

func (h *Handler) GetApplicationsByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	_ = c.GetUint("user_id")

	apps, err := h.service.GetApplicationsByProject(c.Request.Context(), uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, apps)
}

func (h *Handler) UpdateApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid application ID"))
		return
	}

	var req UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	_ = c.GetUint("user_id")

	app, err := h.service.UpdateApplication(c.Request.Context(), uint(id), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, app)
}

func (h *Handler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid application ID"))
		return
	}

	_ = c.GetUint("user_id")

	err = h.service.DeleteApplication(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}