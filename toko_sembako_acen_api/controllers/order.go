package controllers

import (
	"log"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/services"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (o *OrderController) CreateOrderItems(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	jwtToken, err := helpers.GetToken(tokenString)

	if err != nil {
		c.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}

	data, err := helpers.ExtractTokenData(jwtToken)

	if err != nil {
		c.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}

	userid := data["userId"].(string)
	log.Println("id User : " + userid)

	orderItems, err := o.orderService.CreateOrderItems(userid)

	if err != nil {
		c.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(201, gin.H{"status": 201, "message": "Order Items Created Successfully", "data": orderItems})
}

func (o *OrderController) GetOrderItemsById(c *gin.Context) {

	orderId := c.Param("orderId")

	orderItems, err := o.orderService.GetOrderItemById(orderId)

	if err != nil {
		c.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Order Items Retrieved Successfully", "data": orderItems})

}
