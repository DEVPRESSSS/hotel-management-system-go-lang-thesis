package rbac

import (
	"HMS-GO/internal/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("=== RBAC Middleware Debug ===")
		fmt.Println("Request URL:", ctx.Request.URL.Path)
		fmt.Println("All Cookies Received:", ctx.Request.Header.Get("Cookie"))

		tokenString := ctx.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			cookie, err := ctx.Cookie("token")

			if err != nil {
				fmt.Println("ERROR: No token cookie found!")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
				return
			}
			tokenString = cookie
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		permissions := claims["access"].([]interface{})
		for _, p := range permissions {
			if p == permission {
				fmt.Println("SUCCESS: Permission granted")
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
