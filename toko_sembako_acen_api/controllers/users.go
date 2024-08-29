package controllers

import (
	"context"
	"net/http"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/infra/logger"
	"toko_sembako_acen/models"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
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
		logger.Infof("Password Cannot be Empty")
		return
	}

	hashedPass := helpers.HashPassword(user.Password)

	if err := u.db.Where("email = ? ", user.Email).First(&models.Users{}).Error; err != gorm.ErrRecordNotFound {
		ctx.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

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

func (u *UserController) SignIn(ctx *gin.Context) {
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
		logger.Infof("Password Cannot be Empty")
		return
	}

	inputPassword := user.Password

	if err := u.db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&user).Error; err != nil {
		logger.Infof(err.Error())
		ctx.JSON(401, gin.H{"error": "Invalid Email or Password"})
		return
	}

	isPasswordCorrect := helpers.VerifyPassword(user.Password, inputPassword)

	if !isPasswordCorrect {
		logger.Infof("3333333333333")
		ctx.JSON(401, gin.H{"error": "Invalid Email or Password"})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (u *UserController) GoogleSignIn(c *gin.Context) {
	c.Redirect(301, "/auth/google")
}

func (u *UserController) SignInWithProvider(c *gin.Context) {

	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (u *UserController) CallbackHandler(c *gin.Context) {

	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, user)
}
