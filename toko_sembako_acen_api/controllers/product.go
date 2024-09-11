package controllers

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
	"toko_sembako_acen/models"
	"toko_sembako_acen/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (p *ProductController) AddProduct(c *gin.Context) {
	var name, categoryString string
	var stock int
	var price, capital float64

	categoryString = c.PostForm("category")

	category := strings.Split(categoryString, ",")

	// if len(category) == 0 {
	// 	log.Println("Category Can't be empty")
	// 	c.JSON(400, gin.H{"status": 402, "message": "Category Cannot be Empty", "data": nil})
	// 	return
	// }

	name = c.PostForm("name")

	if name == "" {
		log.Println("Name Cannot Be Empty")
		c.JSON(400, gin.H{"status": 402, "message": "Name Cannot be Empty", "data": nil})
		return
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))

	if err != nil {
		if stock <= 0 {
			err = errors.New("Stock input Is Invalid")
		}
		log.Println("Error When Convert Stock Text to integer")
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	price, err = strconv.ParseFloat(c.PostForm("price"), 64)

	if err != nil {
		if price <= 0 {
			err = errors.New("Price input Is Invalid")
		}
		log.Println("Error When Convert Price Text to float : " + err.Error())
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	capital, err = strconv.ParseFloat(c.PostForm("capital"), 64)

	if err != nil {
		if capital <= 0 {
			err = errors.New("Capital input Is Invalid")
		}
		log.Println("Error When Convert Capital Text to float : " + err.Error())
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	pictureFile, err := c.FormFile("picture")

	if err != nil {
		log.Println("Error Post Picture : " + err.Error())
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	product, err := p.productService.AddProduct(&models.Product{
		Name:    name,
		Stock:   stock,
		Capital: capital,
		Price:   price,
	}, category, pictureFile)

	if err != nil {
		c.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Product Created Successfully", "data": product})
}

func (p *ProductController) GetProducts(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	var pageInt, limitInt int

	if page == "" || limit == "" {
		pageInt = 1
		limitInt = 10
	}

	if pageInt, err := strconv.Atoi(page); err != nil || pageInt <= 0 {
		pageInt = 1
	}

	if limitInt, err := strconv.Atoi(limit); err != nil || limitInt <= 0 {
		limitInt = 10
	}

	products, err := p.productService.GetProducts(pageInt, limitInt)

	if err != nil {
		c.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Products Retrieved Successfully", "data": products})
}

func (p *ProductController) DeleteProduct(c *gin.Context) {

	productId := c.Param("id")

	if err := p.productService.DeleteProduct(uuid.MustParse(productId)); err != nil {
		c.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Products Deleted Successfully", "data": productId})

}

func (p *ProductController) GetProductsByCategoryAndSearch(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	var pageInt, limitInt int

	if page == "" || limit == "" {
		pageInt = 1
		limitInt = 10
	}
	pageInt, err := strconv.Atoi(page)

	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err = strconv.Atoi(limit)

	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	search := c.Query("search")
	categoryString := c.Query("category")

	category := strings.Split(categoryString, ",")

	products, err := p.productService.GetProductsByCategoryAndSearch(category, search, pageInt, limitInt)

	if err != nil {
		c.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Products Retrieved Successfully", "data": products})
}

func (p *ProductController) UpdateProduct(c *gin.Context) {
	var name, categoryString string
	var stock int
	var price, capital float64

	productId := c.Param("id")

	if productId == "" {
		log.Println("Invalid Url parameter")
		c.JSON(400, gin.H{"status": 402, "message": "Invalid URL parameter", "data": nil})
		return
	}

	categoryString = c.PostForm("category")

	category := strings.Split(categoryString, ",")

	if len(category) == 0 {
		log.Println("Category Can't be empty")
		c.JSON(400, gin.H{"status": 402, "message": "Category Cannot be Empty", "data": nil})
		return
	}

	name = c.PostForm("name")

	if name == "" {
		log.Println("Name Cannot Be Empty")
		c.JSON(400, gin.H{"status": 402, "message": "Name Cannot be Empty", "data": nil})
		return
	}

	stock, err := strconv.Atoi(c.PostForm("stock"))

	if err != nil {
		if stock <= 0 {
			err = errors.New("Stock input Is Invalid")
		}
		log.Println("Error When Convert Stock Text to integer")
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	price, err = strconv.ParseFloat(c.PostForm("price"), 64)

	if err != nil {
		if price <= 0 {
			err = errors.New("Price input Is Invalid")
		}
		log.Println("Error When Convert Price Text to float : " + err.Error())
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	capital, err = strconv.ParseFloat(c.PostForm("capital"), 64)

	if err != nil {
		if capital <= 0 {
			err = errors.New("Capital input Is Invalid")
		}
		log.Println("Error When Convert Capital Text to float : " + err.Error())
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	pictureFile, err := c.FormFile("picture")

	if err != nil {
		log.Println("Error Post Picture : " + err.Error())
		c.JSON(400, gin.H{"status": 402, "message": err.Error(), "data": nil})
		return
	}

	product, err := p.productService.UpdateProduct(&models.Product{
		Id:        uuid.MustParse(productId),
		Name:      name,
		Stock:     stock,
		Capital:   capital,
		Price:     price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}, category, pictureFile)

	if err != nil {
		c.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"status": 200, "message": "Product Created Successfully", "data": product})
}
