package middlewares

import (
	"HMS-GO/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.ValidateToken(c)

		if err != nil {
			c.HTML(http.StatusUnauthorized, "errors.html", gin.H{
				"Unauthorized": "Authentication required",
			})
			fmt.Println(err)
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
