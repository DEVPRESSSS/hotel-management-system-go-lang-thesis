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

		tokenString := ctx.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			cookie, err := ctx.Cookie("token")

			if err != nil {
				ctx.Redirect(http.StatusUnauthorized, "/login")
				return
			}
			tokenString = cookie
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		rawAccess, exists := claims["access"]
		fmt.Println(rawAccess)
		if !exists || rawAccess == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no permissions assigned to role",
			})
			// ctx.HTML(http.StatusForbidden, "404.html", gin.H{
			// 	"code":    403,
			// 	"message": "Access Forbidden",
			// })

			return
		}

		permissions, ok := rawAccess.([]interface{})
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid permissions format",
			})
			return
		}
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
