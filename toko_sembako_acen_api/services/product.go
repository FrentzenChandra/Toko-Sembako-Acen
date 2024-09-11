package services

import (
	"errors"
	"log"
	"mime/multipart"
	"time"
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

	if err := p.db.Model(&models.Product{}).Where("id = ?", productId).Update("deleted_at", time.Now()).Error; err != nil {
		log.Println("Repository error When deleting Product : ", err.Error())
		return err
	}

	return nil

}

func (p *ProductService) GetProductsByCategoryAndSearch(category []string, search string) ([]*models.Product, error) {
	var categoryQuery string
	var products []*models.Product
	search = "%" + search + "%"
	log.Println(len(category))

	if len(category) > 0 && category[0] != "" {
		for index, _ := range category {

			if index == 0 {
				categoryQuery += "id = ? "
			} else {
				categoryQuery += " OR id = ? "
			}
		}
	}

	// Preload Orders with conditions
	if len(category) > 0 && category[0] != "" {
		var productIds []string
		err := p.db.Model(&models.Product{}).
			Preload("Categories").
			Select("Distinct product.id ").
			Joins("inner join product_category on product.id = product_category.product_id").
			Where(" product.name ILIKE ? and product_category.category_id in ? ", search, category).
			Scan(&products).Error

		if err != nil {
			log.Println(err)
			return nil, err
		}

		for _, product := range products {
			productIds = append(productIds, product.Id.String())
		}

		if err := p.db.Preload("Categories").Where("id IN ?", productIds).Find(&products).Error; err != nil {
			log.Println("Error Repository When Find Products By Category And Search : " + err.Error())
			return nil, err
		}

	} else {
		if err := p.db.Preload("Categories").Where("name ILIKE ?", search).Find(&products).Error; err != nil {
			log.Println("Error Repository When Find Products By Category And Search : " + err.Error())
			return nil, err
		}
	}

	return products, nil

}
