package controllers

import (
	"toko_sembako_acen/models"
	"toko_sembako_acen/services"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (c *CategoryController) Create(ctx *gin.Context) {

	var category *models.Category

	if err := ctx.ShouldBindBodyWithJSON(&category); err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	category, err := c.categoryService.AddCategory(category)

	if err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(201, gin.H{"status": 201, "message": "Category created successfully", "data": category})

}

func (c *CategoryController) GetCategories(ctx *gin.Context) {
	categories, err := c.categoryService.GetCategories()

	if err != nil {
		ctx.JSON(400, gin.H{"status": 400, "message": err.Error(), "data": nil})
		return
	}

	ctx.JSON(200, gin.H{"status": 200, "message": "Categories retrieved successfully", "data": categories})
}
