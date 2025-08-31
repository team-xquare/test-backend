package project

import (
	"net/http"
	"strconv"

	"github.com/team-xquare/deployment-platform/internal/pkg/middleware"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	projects := r.Group("/projects")
	projects.Use(middleware.Auth())
	{
		projects.POST("", h.CreateProject)
		projects.GET("", h.GetProjects)
		projects.GET("/:id", h.GetProject)
		projects.DELETE("/:id", h.DeleteProject)
		projects.POST("/:id/applications", h.DeployApplication)
		projects.POST("/:id/addons", h.DeployAddon)
	}
}

func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	project, err := h.service.CreateProject(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *Handler) GetProjects(c *gin.Context) {
	userID := c.GetUint("user_id")
	projects, err := h.service.GetProjects(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *Handler) GetProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	userID := c.GetUint("user_id")
	project, err := h.service.GetProject(c.Request.Context(), userID, uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.DeleteProject(c.Request.Context(), userID, uint(projectID)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func (h *Handler) DeployApplication(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	var req DeployApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.DeployApplication(c.Request.Context(), userID, uint(projectID), req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deployment initiated"})
}

func (h *Handler) DeployAddon(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.BadRequest("Invalid project ID"))
		return
	}

	var req DeployAddonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("Invalid request format"))
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.DeployAddon(c.Request.Context(), userID, uint(projectID), req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Addon deployment initiated"})
}
