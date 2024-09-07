package controllers

import (
	"log"
	"toko_sembako_acen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	db *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{db: db}
}

func (c *CategoryController) AddCategory(ctx *gin.Context) {
	var category *models.Category

	err := ctx.BindJSON(&category)

	if err != nil {
		log.Println("Error When Bind Json Body Category Controller : " + err.Error())
		ctx.JSON(400, "Server Error")
		return
	}

	if category.Name == "" {
		log.Println("Error Name Cannot Empty")
		ctx.JSON(422, "Server Error : Name Cannot Be Empty")
		return
	}

	if RowsAffected := c.db.Where(&models.Category{
		Name: category.Name,
	}).Take(&models.Category{}).RowsAffected; RowsAffected != 0 {
		log.Println("Error Name Already Exists")
		ctx.JSON(409, "Category Already Exists")
		return
	}

	if err := c.db.Create(category).Error; err != nil {
		log.Println("Error When Create Category : " + err.Error())
		ctx.JSON(500, "Server Error")
		return
	}

	ctx.JSON(201, "Category Successfully Created")
}
