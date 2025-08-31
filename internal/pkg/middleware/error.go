package middleware

import (
	"net/http"

	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appError, ok := err.(*errors.AppError); ok {
				c.JSON(appError.StatusCode, appError)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status_code": http.StatusInternalServerError,
					"message":     "Internal server error",
					"type":        "INTERNAL_SERVER_ERROR",
				})
			}
		}
	}
}
