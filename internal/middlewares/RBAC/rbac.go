package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RBACMiddleware(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		rawAccess, exists := ctx.Get("access")
		if !exists || rawAccess == nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no permissions assigned to role",
			})
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
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
