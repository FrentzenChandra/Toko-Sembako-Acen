package services

import (
	"errors"
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
	tx := p.db.Begin()

	if err := tx.Create(&product).Error; err != nil {
		log.Println("Error Repository When Create Product : " + err.Error())
		return nil, err
	}

	picUrl, err := helpers.UploadToCloudinary(pictureFile)

	if err != nil {
		log.Println("Repository Error When Upload Picture To cloudinary : " + err.Error())
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&product).Where("id = ?", product.Id).Update("picture", picUrl).Error; err != nil {
		log.Println("Error Repository When Update Product : " + err.Error())
		tx.Rollback()
		return nil, err
	}

	// jalankan codingan
	tx.Commit()

	if err := p.AddProductCategory(category, product.Id); err != nil {

		if err := p.db.Delete(&product).Error; err != nil {
			log.Println("Error Repository When Delete Product After Error : " + err.Error())
			return nil, err
		}

		return nil, err
	}

	return product, nil
}

func (p *ProductService) AddProductCategory(category []string, productId uuid.UUID) error {

	for _, categoryId := range category {
		if RowsAffected := p.db.Where(&models.ProductCategory{
			ProductID:  productId,
			CategoryID: uuid.MustParse(categoryId),
		}).Take(&models.Category{}).RowsAffected; RowsAffected != 0 {
			return errors.New("Duplicate Product Category")
		}

		if err := p.db.Create(&models.ProductCategory{
			ProductID:  productId,
			CategoryID: uuid.MustParse(categoryId),
		}).Error; err != nil {
			log.Println("Error Repository When Create Product Category : " + err.Error())
			return err
		}
	}

	return nil
}

func (p *ProductService) GetProducts() ([]*models.Product, error) {
	var products []*models.Product

	if err := p.db.Where("deleted_at IS NULL").Preload("Categories").Find(&products).Error; err != nil {
		log.Println("Error Repository When Find Products : " + err.Error())
		return nil, err
	}

	return products, nil

}

func (p *ProductService) DeleteProduct(productId uuid.UUID) error {

	if err := p.db.Delete(&models.Product{
		Id: productId,
	}).Error; err != nil {
		log.Println("Repository error When deleting Product : ", err.Error())
		return err
	}
	

	return nil

}
