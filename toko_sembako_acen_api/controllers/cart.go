package controllers

import (
	"log"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/models"
	"toko_sembako_acen/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartController struct {
	cartService *services.CartItemService
}

func NewCartController(cartService *services.CartItemService) *CartController {
	return &CartController{cartService: cartService}
}

func (c *CartController) AddCartItem(ctx *gin.Context) {
	var input *models.CartItemInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	tokenString := ctx.GetHeader("Authorization")
	jwtToken, err := helpers.GetToken(tokenString)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}
	data, err := helpers.ExtractTokenData(jwtToken)
	log.Println(input)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}

	userid := data["userId"].(string)
	userUUID := uuid.MustParse(userid)

	if err := c.cartService.AddCartItem(input, userUUID); err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(201, gin.H{"status": 201, "message": "Cart Item Added Successfully", "data": nil})
}
