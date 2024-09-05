package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{db: db}
}

func (p *ProductController) AddProduct(cloud *cloudinary.Cloudinary, c *gin.Context) {

	name := c.PostForm("name")
	stock, err := strconv.Atoi(c.PostForm("stock"))

	if err != nil {
		log.Println("Error When Convert Stock Text to integer")
		c.JSON(400, "Server Error : Data Stock Is not Integer")
		return
	}

	price, err := strconv.ParseFloat(c.PostForm("price"), 64)

	if err != nil {
		log.Println("Error When Convert Price Text to float : " + err.Error())
		c.JSON(400, "Server Error : Data Price Is not Float")
		return
	}

	capital, err := strconv.ParseFloat(c.PostForm("capital"), 64)

	if err != nil {
		log.Println("Error When Convert Capital Text to float : " + err.Error())
		c.JSON(400, "Server Error : Data Capital Is not Float")
		return
	}

	nowTime := time.Now()
	pictureFileHeader, err := c.FormFile("picture")

	if err != nil {
		log.Println("Error Post Picture : " + err.Error())
		c.JSON(400, "Server Error : Data Picture Failed")
		return
	}

	if name == "" {
		log.Println("Name Cannot Be Empty")
		c.JSON(422, "Server Error : Data Stock Is not Integer")
		return
	}

	

}
