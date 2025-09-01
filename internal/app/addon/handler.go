package addon

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
	addons := r.Group("/addons")
	addons.Use(middleware.Auth())
	{
		addons.GET("/:id", h.GetAddon)
		addons.PUT("/:id", h.UpdateAddon)
		addons.DELETE("/:id", h.DeleteAddon)
	}

	// Project-specific addon routes
	projects := r.Group("/projects/:id/addons")
	projects.Use(middleware.Auth())
	{
		projects.GET("", h.GetAddonsByProject)
		projects.POST("", h.CreateAddon)
	}
}

func (h *Handler) CreateAddon(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	var req CreateAddonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	_ = c.GetUint("user_id")

	addon, err := h.service.CreateAddon(c.Request.Context(), uint(projectID), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, addon)
}

func (h *Handler) GetAddon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid addon ID"))
		return
	}

	_ = c.GetUint("user_id")

	addon, err := h.service.GetAddon(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, addon)
}

func (h *Handler) GetAddonsByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	_ = c.GetUint("user_id")

	addons, err := h.service.GetAddonsByProject(c.Request.Context(), uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, addons)
}

func (h *Handler) UpdateAddon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid addon ID"))
		return
	}

	var req UpdateAddonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	_ = c.GetUint("user_id")

	addon, err := h.service.UpdateAddon(c.Request.Context(), uint(id), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, addon)
}

func (h *Handler) DeleteAddon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid addon ID"))
		return
	}

	_ = c.GetUint("user_id")

	err = h.service.DeleteAddon(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Addon deleted successfully"})
}