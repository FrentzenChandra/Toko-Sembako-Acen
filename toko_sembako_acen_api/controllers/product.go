package controllers

import (
	"log"
	"mime/multipart"
	"runtime"
	"strconv"
	"strings"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{db: db}
}

func (p *ProductController) AddProduct(c *gin.Context) {
	runtime.GOMAXPROCS(2)

	var name, picUrl, category string
	var stock int
	var price, capital float64
	var pictureFile *multipart.FileHeader
	var product models.Product

	category = c.PostForm("category")

	if category == "" {
		log.Println("Category Can't be empty")
		c.JSON(422, "Server Error : Data Category Is Empty")
		return
	}

	name = c.PostForm("name")

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

	price, err = strconv.ParseFloat(c.PostForm("price"), 64)

	if err != nil {
		log.Println("Error When Convert Price Text to float : " + err.Error())
		c.JSON(400, "Server Error : Data Price Is not Float")
		return
	}

	capital, err = strconv.ParseFloat(c.PostForm("capital"), 64)

	if err != nil {
		log.Println("Error When Convert Capital Text to float : " + err.Error())
		c.JSON(400, "Server Error : Data Capital Is not Float")
		return
	}

	pictureFile, err = c.FormFile("picture")

	if err != nil {
		log.Println("Error Post Picture : " + err.Error())
		c.JSON(400, "Server Error : Data Picture Failed")
		return
	}

	product.Name = name
	product.Stock = stock
	product.Price = price
	product.Capital = capital

	if err := p.db.Create(&product).Error; err != nil {
		log.Println("Error Create Product : " + err.Error())
		c.JSON(400, "Error Server ")
		return
	}

	picUrl, err = helpers.UploadToCloudinary(pictureFile)

	if err != nil {
		log.Println(err)
		c.JSON(500, "Server Error : Failed to upload picture Err : "+err.Error())
		return
	}

	if err := p.db.Model(&models.Product{}).Where("id = ?", product.Id).Update("picture", picUrl).Error; err != nil {
		log.Println("Error Update Product Picture : " + err.Error())
		c.JSON(500, "Server Error : Failed to Create Product Picture")
		return
	}

	arrCategory := strings.Split(category, ",")

	for _, categoryId := range arrCategory {
		p.db.Create(&models.ProductCategory{
			ProductID:  product.Id,
			CategoryID: uuid.MustParse(categoryId),
		})
	}

	c.JSON(201, "Product Created Successfully!!!")
}
