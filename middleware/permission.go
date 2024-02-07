package middleware

import (
	"net/http"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/models"

	"github.com/gin-gonic/gin"
)

func Permission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		User := ctx.MustGet("User").(models.User)

		if User.Role == "admin" {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Permission Denied"})
			return
		}

	}
}