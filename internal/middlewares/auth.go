package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var tokenStr string

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := c.Cookie("token")
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			tokenStr = cookie
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("secret_key")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["userid"])
		c.Set("access", claims["access"])

		c.Next()
	}
}
