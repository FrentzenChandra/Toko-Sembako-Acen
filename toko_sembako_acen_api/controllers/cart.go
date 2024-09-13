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

func (c *CartController) GetCartItems(ctx *gin.Context) {

	tokenString := ctx.GetHeader("Authorization")
	jwtToken, err := helpers.GetToken(tokenString)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}
	data, err := helpers.ExtractTokenData(jwtToken)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}

	userid := data["userId"].(string)
	userUUID := uuid.MustParse(userid)
	log.Println(userUUID)

	cartItems, err := c.cartService.GetCartItems(userUUID)

	if err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(200, gin.H{"status": 200, "message": "Cart Item Retrieved Successfully", "data": cartItems})
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

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}

	userid := data["userId"].(string)
	userUUID := uuid.MustParse(userid)

	if err := c.cartService.AddCartItem(input, userUUID); err != nil {

		if err.Error() == "Not Enough Stock" {
			ctx.JSON(200, gin.H{"status": 200, "message": err.Error(), "data": nil})
			return
		}

		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(201, gin.H{"status": 201, "message": "Cart Item Added Successfully", "data": nil})
}

func (c *CartController) UpdateCartItem(ctx *gin.Context) {
	var input *models.CartItemInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	productId := ctx.Param("productId")

	tokenString := ctx.GetHeader("Authorization")
	jwtToken, err := helpers.GetToken(tokenString)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}
	data, err := helpers.ExtractTokenData(jwtToken)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}

	userid := data["userId"].(string)
	userUUID := uuid.MustParse(userid)

	if err := c.cartService.UpdateCartItem(input, userUUID, productId); err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(200, gin.H{"status": 200, "message": "Cart Item Updated Successfully", "data": nil})
}
