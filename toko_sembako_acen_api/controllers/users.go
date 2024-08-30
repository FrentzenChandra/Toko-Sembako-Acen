package controllers

import (
	"context"
	"log"
	"net/http"
	"toko_sembako_acen/helpers"
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
		log.Println(err.Error(), user)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" {
		log.Println("Email Cannot be Empty")
		return
	}

	if user.Password == "" {
		log.Println("Password Cannot be Empty")
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
		log.Println(err.Error(), user)
		ctx.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, "Successfully Created user ")
}

func (u *UserController) Login(ctx *gin.Context) {

	var user models.Users

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Println(err.Error(), user)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" {
		log.Println("Email Cannot be Empty")
		return
	}

	if user.Password == "" {
		log.Println("Password Cannot be Empty")
		return
	}

	inputPassword := user.Password

	if err := u.db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&user).Error; err != nil {
		log.Println(err.Error())
		ctx.JSON(401, gin.H{"error": "Invalid Email or Password"})
		return
	}

	isPasswordCorrect := helpers.VerifyPassword(user.Password, inputPassword)

	if !isPasswordCorrect {
		ctx.JSON(401, gin.H{"error": "Invalid Email or Password"})
		return
	}

	tokenString, err := helpers.CreateToken(user.Id, user.Username, user.Email)

	if err != nil {
		ctx.JSON(400, err)
		return
	}

	ctx.JSON(200, tokenString)
}

func (u *UserController) UserList(c *gin.Context) {

	var users []*models.Users

	if err := u.db.Where("deleted_at IS NULl").Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
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
