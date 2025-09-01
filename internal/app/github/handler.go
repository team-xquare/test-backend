package github

import (
	"io"
	"net/http"

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
	github := r.Group("/github")
	{
		github.POST("/webhook", h.HandleWebhook)

		authenticated := github.Group("")
		authenticated.Use(middleware.Auth())
		{
			authenticated.GET("/installations", h.GetInstallations)
			authenticated.GET("/installations/:id/repositories", h.GetRepositories)
			authenticated.POST("/installations/:id/link", h.LinkInstallation)
		}
	}
}

func (h *Handler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		c.Error(errors.BadRequest("Missing signature"))
		return
	}

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(errors.BadRequest("Failed to read payload"))
		return
	}

	if err := h.service.HandleInstallationWebhook(c.Request.Context(), payload, signature); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) GetInstallations(c *gin.Context) {
	userID := c.GetUint("user_id")

	installations, err := h.service.GetUserInstallations(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, installations)
}

func (h *Handler) GetRepositories(c *gin.Context) {
	installationID := c.Param("id")

	repositories, err := h.service.GetRepositories(c.Request.Context(), installationID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, repositories)
}

func (h *Handler) LinkInstallation(c *gin.Context) {
	installationID := c.Param("id")
	userID := c.GetUint("user_id")

	err := h.service.LinkInstallationToUser(c.Request.Context(), userID, installationID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installation linked successfully"})
}
