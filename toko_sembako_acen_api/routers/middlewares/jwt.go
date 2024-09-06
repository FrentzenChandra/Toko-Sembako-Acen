package middlewares

import (
	"log"
	"net/http"
	"toko_sembako_acen/helpers"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	tokenString := ctx.GetHeader("Authorization")

	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, "You are not Authorized to access this")
		ctx.Abort()
		return
	}

	if err := helpers.VerifyToken(tokenString); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusUnauthorized, "Invalid Token")
		ctx.Abort()
		return
	}

	ctx.Next()
}
