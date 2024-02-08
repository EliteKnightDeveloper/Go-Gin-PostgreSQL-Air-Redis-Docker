package middleware

import (
	"net/http"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/models"

	"github.com/gin-gonic/gin"
)

func Permission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		User := ctx.MustGet("User").(models.User)

		if User.Role == "Admin" {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission Denied"})
			return
		}

	}
}
