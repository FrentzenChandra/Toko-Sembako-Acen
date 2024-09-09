package services

import (
	"errors"
	"log"
	"toko_sembako_acen/models"

	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (c *CategoryService) AddCategory(category *models.Category) (*models.Category, error) {

	if category.Name == "" {
		return nil, errors.New("Category Name Cannot Be empty")
	}

	if RowsAffected := c.db.Where(&models.Category{
		Name: category.Name,
	}).Take(&models.Category{}).RowsAffected; RowsAffected != 0 {
		return nil, errors.New("Category Already exists")
	}

	if err := c.db.Create(&category).Error; err != nil {
		log.Println("Category Errors When Creating Category : " + err.Error())
		return nil, err
	}

	return category, nil
}


// func (c *CategoryService) GetCategories() (*[]models.Category, error) {

// }
