package controllers

import (
	"log"
	"toko_sembako_acen/models"
	"toko_sembako_acen/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) Login(ctx *gin.Context) {
	var user models.Users

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	if user.Email == "" {
		log.Println("Email Cannot be Empty")
		ctx.JSON(400, gin.H{"status": 400, "message": "Email Cannot be Empty", "data": nil})
		return
	}

	if user.Password == "" {
		log.Println("Password Cannot be Empty")
		ctx.JSON(400, gin.H{"status": 400, "message": "Password Cannot be Empty", "data": nil})
		return
	}

	token, err := u.userService.Login(&user)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 401, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(200, gin.H{"status": 200, "message": "Login Success", "data": token})
}

func (u *UserController) Register(ctx *gin.Context) {
	var user models.Users

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	if user.Email == "" {
		log.Println("Email Cannot be Empty")
		ctx.JSON(400, gin.H{"status": 400, "message": "Email Cannot be Empty", "data": nil})
		return
	}

	if user.Password == "" {
		log.Println("Password Cannot be Empty")
		ctx.JSON(400, gin.H{"status": 400, "message": "Password Cannot be Empty", "data": nil})
		return
	}

	userId, err := u.userService.SignUp(&user)

	if err != nil {
		ctx.JSON(401, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(200, gin.H{"status": 201, "message": "Account Successfully Registered", "data": userId})
}
