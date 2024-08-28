package controllers

import (
	"net/http"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/infra/logger"
	"toko_sembako_acen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (u *UserController) SignUp(ctx *gin.Context) {
	var user models.Users

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		logger.Infof(err.Error(), user)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" {
		logger.Infof("Email Cannot be Empty")
		return
	}

	if user.Password == "" {
		logger.Infof("Email Cannot be Empty")
		return
	}

	hashedPass := helpers.HashPassword(user.Password)

	if err := u.db.Create(&models.Users{
		Email:    user.Email,
		Password: hashedPass,
	}).Error; err != nil {
		logger.Infof(err.Error(), user)
		ctx.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, "Successfully Created user ")
}
