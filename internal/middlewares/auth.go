package middlewares

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user, err := utils.ValidateToken(c)

// 		if err != nil {
// 			c.HTML(http.StatusUnauthorized, "errors.html", gin.H{
// 				"Unauthorized": "Authentication required",
// 			})
// 			fmt.Println(err)
// 			c.Abort()
// 			return
// 		}
// 		c.Set("user", user)
// 		c.Next()
// 	}
// }

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
// 			return
// 		}

// 		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

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
		if authHeader != "" {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := c.Cookie("token")
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
				return
			}
			tokenStr = cookie
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("secret_key")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["userid"])
		c.Set("access", claims["access"])

		c.Next()
	}
}
