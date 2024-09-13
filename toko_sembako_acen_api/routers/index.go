package routers

import (
	"net/http"
	"toko_sembako_acen/controllers"
	"toko_sembako_acen/infra/database"
	"toko_sembako_acen/routers/middlewares"
	"toko_sembako_acen/services"

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
	userService := services.NewUserService(DB)
	userControllers := controllers.NewUserController(userService)

	route.POST("/users/signup", userControllers.Register)
	route.POST("/users/signin", userControllers.Login)

	route.Use(middlewares.JwtMiddleware)
	route.GET("/users", userControllers.GetUsers)
}

func ProductRoutes(route *gin.Engine) {
	productService := services.NewProductService(DB)
	ProductsControllers := controllers.NewProductController(productService)
	route.Use(middlewares.JwtMiddleware)

	route.POST("/product", ProductsControllers.AddProduct)

	// route.GET("/product", ProductsControllers.GetProducts)
	route.GET("/product", ProductsControllers.GetProductsByCategoryAndSearch)

	route.DELETE("/product/:id", ProductsControllers.DeleteProduct)

	route.PATCH("/product/:id", ProductsControllers.UpdateProduct)
}

func CategoryRoutes(route *gin.Engine) {
	categoryService := services.NewCategoryService(DB)
	categoryControllers := controllers.NewCategoryController(categoryService)
	route.Use(middlewares.JwtMiddleware)

	route.POST("/category", categoryControllers.Create)

	route.GET("/category", categoryControllers.GetCategories)
}

func CartRoutes(route *gin.Engine) {
	cartItemService := services.NewCartItemService(DB)
	cartItemController := controllers.NewCartController(cartItemService)

	route.Use(middlewares.JwtMiddleware)

	route.POST("/cart", cartItemController.AddCartItem)

	route.GET("/cart", cartItemController.GetCartItems)

	route.PATCH("/cart/:productId", cartItemController.UpdateCartItem)
}

func OrderRoutes(route *gin.Engine) {
	orderService := services.NewOrderService(DB)
	orderController := controllers.NewOrderController(orderService)

	route.Use(middlewares.JwtMiddleware)

	route.POST("/order", orderController.CreateOrderItems)

	route.GET("/order/:orderId", orderController.GetOrderItemsById)
}

//func LocalHostRoute(route *gin.Engine) {
// 	userControllers := controllers.NewUserController(DB)

// 	route.GET("/users/google/signin", userControllers.GoogleSignIn)
// 	route.GET("/auth/:provider", userControllers.SignInWithProvider)
// 	route.GET("/auth/:provider/callback", userControllers.CallbackHandler)
// }
