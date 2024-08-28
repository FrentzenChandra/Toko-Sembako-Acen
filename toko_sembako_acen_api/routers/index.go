package routers

import (
	"net/http"
	"toko_sembako_acen/controllers"
	"toko_sembako_acen/infra/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	DB = database.GetDB()

	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	//Add All route
	//TestRoutes(route)
}

func UserRoutes(route *gin.Engine) {

	UserControllers := controllers.NewUserController(DB)
	route.POST("/users/signup", UserControllers.SignUp)
}
