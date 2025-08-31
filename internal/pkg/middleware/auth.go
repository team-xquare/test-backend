package middleware

import (
	"strings"

	"github.com/team-xquare/deployment-platform/internal/pkg/utils/errors"
	"github.com/team-xquare/deployment-platform/internal/pkg/utils/jwt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errors.Unauthorized("Authorization header required"))
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.Error(errors.Unauthorized("Invalid authorization header format"))
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(bearerToken[1])
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
