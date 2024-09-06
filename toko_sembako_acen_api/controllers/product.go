package controllers

import (
	"log"
	"strconv"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{db: db}
}

func (p *ProductController) AddProduct(c *gin.Context) {

	name := c.PostForm("name")

	if name == "" {
		log.Println("Name Cannot Be Empty")
		c.JSON(422, "Server Error : Data Stock Is not Integer")
		return
	}

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

	pictureFile, err := c.FormFile("picture")

	if err != nil {
		log.Println("Error Post Picture : " + err.Error())
		c.JSON(400, "Server Error : Data Picture Failed")
		return
	}

	picUrl, err := helpers.UploadToCloudinary(pictureFile)

	if err != nil {
		log.Println(err)
		c.JSON(500, "Server Error : Failed to upload picture")
		return
	}

	log.Println(picUrl)

	if err := p.db.Create(&models.Product{
		Name:    name,
		Stock:   stock,
		Price:   price,
		Capital: capital,
		Picture: picUrl,
	}).Error; err != nil {
		log.Println("Error Create Product : " + err.Error())
		c.JSON(400, "Error Server ")
		return
	}

	c.JSON(201, "Product Created Successfully!!!")
}
