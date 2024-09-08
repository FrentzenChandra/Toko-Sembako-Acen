package services

import (
	"errors"
	"log"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) SignUp(user *models.Users) (*uuid.UUID, error) {

	// if user.Email == "" {
	// 	log.Println("Email Cannot be Empty")
	// 	return nil, errors.New("Email Cannot Be empty")
	// }

	// if user.Password == "" {
	// 	log.Println("Password Cannot be Empty")
	// 	return nil, errors.New("Password Cannot Be empty")
	// }

	hashedPass := helpers.HashPassword(user.Password)

	if err := u.db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&models.Users{}).Error; err != gorm.ErrRecordNotFound {
		log.Println("Service Error Get User Detail : " + err.Error())
		return nil, err
	}

	if err := u.db.Create(&models.Users{
		Email:    user.Email,
		Password: hashedPass,
	}).Error; err != nil {
		log.Println("Service Error Creating User : " + err.Error())
		return nil, err
	}

	return &user.Id, nil
}

func (u *UserService) Login(user *models.Users) (*string, error) {
	
	inputPassword := user.Password

	if err := u.db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&user).Error; err != nil {
		log.Println("Repostiory Error Get user Detail : " + err.Error())
		return nil, err
	}

	isPasswordCorrect := helpers.VerifyPassword(user.Password, inputPassword)

	if !isPasswordCorrect {
		return nil, errors.New("Email Or Password Is Invalid ")
	}

	tokenString, err := helpers.CreateToken(user.Id, user.Username, user.Email)

	if err != nil {
		return nil, errors.New("failed To create Token : " + err.Error())
	}

	return &tokenString, nil
}

// func (u *UserService) UserList(c *gin.Context) {

// 	var users []*models.Users

// 	if err := u.db.Where("deleted_at IS NULL").Find(&users).Error; err != nil {
// 		log.Println("Service Error Get users : " + err.Error())
// 		return
// 	}

// 	c.JSON(200, users)
// }

// func (u *UserService) GoogleSignIn(c *gin.Context) {
// 	c.Redirect(301, "/auth/google")
// }

// func (u *UserService) SignInWithProvider(c *gin.Context) {

// 	provider := c.Param("provider")
// 	q := c.Request.URL.Query()
// 	q.Add("provider", provider)
// 	c.Request.URL.RawQuery = q.Encode()

// 	gothic.BeginAuthHandler(c.Writer, c.Request)
// }

// func (u *UserService) CallbackHandler(c *gin.Context) {

// 	provider := c.Param("provider")
// 	q := c.Request.URL.Query()
// 	q.Add("provider", provider)
// 	c.Request.URL.RawQuery = q.Encode()

// 	c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

// 	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
// 	if err != nil {
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	c.JSON(200, user)
// }
