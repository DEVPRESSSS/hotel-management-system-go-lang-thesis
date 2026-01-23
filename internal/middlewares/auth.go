package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var tokenStr string

// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader != "" {
// 			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
// 			if strings.Contains(c.GetHeader("Accept"), "text/html") {

// 				c.Redirect(http.StatusFound, "/login")
// 				c.Abort()
// 				return
// 			}
// 		} else {
// 			cookie, err := c.Cookie("token")
// 			if err != nil {
// 				c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
// 				return
// 			}
// 			tokenStr = cookie
// 		}

// 		// Parse and validate token
// 		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("secret_key")), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
// 			return
// 		}

// 		claims := token.Claims.(jwt.MapClaims)
// 		c.Set("user_id", claims["userid"])
// 		c.Set("access", claims["access"])

// 		c.Next()
// 	}
// }

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
