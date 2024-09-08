package services

import (
	"log"
	"mime/multipart"
	"toko_sembako_acen/helpers"
	"toko_sembako_acen/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (p *ProductService) AddProduct(product *models.Product, category []string, pictureFile *multipart.FileHeader) (*models.Product, error) {



	if err := p.db.Create(&product).Error; err != nil {
		log.Println("Error Repository When Create Product : " + err.Error())
		return nil, err
	}

	picUrl, err := helpers.UploadToCloudinary(pictureFile)

	if err != nil {
		log.Println("Repository Error When Upload Picture To cloudinary : " + err.Error())
		return nil, err
	}

	if err := p.db.Model(&models.Product{}).Where("id = ?", product.Id).Update("picture", picUrl).Error; err != nil {
		log.Println("Error Repository When Update Product : " + err.Error())
		return nil, err
	}

	for _, categoryId := range category {
		if err := p.db.Create(&models.ProductCategory{
			ProductID:  product.Id,
			CategoryID: uuid.MustParse(categoryId),
		}).Error; err != nil {
			log.Println("Error Repository When Create Product Category : " + err.Error())
			return nil, err
		}
	}

	return product, nil
}
