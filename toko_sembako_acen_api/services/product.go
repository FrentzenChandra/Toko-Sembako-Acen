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
		log.Println("Error Service When Create Product : " + err.Error())
		return nil, err
	}

	picUrl, err := helpers.UploadToCloudinary(pictureFile)

	if err != nil {
		log.Println("Service Error When Upload Picture To cloudinary : " + err.Error())
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
			log.Println("Error Service When Create Product Category : " + err.Error())
			return err
		}
	}

	return nil
}

func (p *ProductService) GetProducts(page int, limit int) ([]*models.Product, error) {
	var products []*models.Product

	if err := p.db.Where("deleted_at IS NULL").Preload("Categories").Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
		log.Println("Error Service When Find Products : " + err.Error())
		return nil, err
	}

	return products, nil

}

func (p *ProductService) DeleteProduct(productId uuid.UUID) error {

	if err := p.db.Model(&models.Product{}).Where("id = ?", productId).Update("deleted_at", time.Now()).Error; err != nil {
		log.Println("Service error When deleting Product : ", err.Error())
		return err
	}

	return nil

}

func (p *ProductService) GetProductsByCategoryAndSearch(category []string, search string, page int, limit int) ([]*models.Product, error) {
	var categoryQuery string
	var products []*models.Product
	search = "%" + search + "%"

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
			Where(" product.name ILIKE ? and product_category.category_id in ? And deleted_at IS NULL", search, category).
			Scan(&products).Error

		if err != nil {
			log.Println(err)
			return nil, err
		}

		for _, product := range products {
			productIds = append(productIds, product.Id.String())
		}

		if err := p.db.Preload("Categories").Where("id IN ?", productIds).Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
			log.Println("Error Service When Find Products By Category And Search : " + err.Error())
			return nil, err
		}

		if len(products) == 0 {
			return nil, errors.New("Product Not Found")
		}

	} else {
		log.Println(limit)
		log.Println(page)
		if err := p.db.Preload("Categories").Limit(limit).Offset((page-1)*limit).Where("name ILIKE ? AND deleted_at IS NULL", search).Find(&products).Error; err != nil {
			log.Println("Error Service When Find Products By Category And Search : " + err.Error())
			return nil, err
		}
	}

	return products, nil

}

func (p *ProductService) UpdateProduct(product *models.Product, categories []string, pictureFile *multipart.FileHeader) (*models.Product, error) {
	var productBeforeUpdt *models.Product
	var errChan = make(chan error)
	var picUrlChan = make(chan string)

	rowsAffected := p.db.Preload("Categories").Where("id = ? AND deleted_at is NULL", product.Id).First(&productBeforeUpdt).RowsAffected

	if rowsAffected == 0 {
		return nil, errors.New("Product Not Found")
	}

	tx := p.db.Begin()

	for _, category := range productBeforeUpdt.Categories {
		if err := tx.Where("product_id = ? AND category_id = ?", product.Id, category.Id).Delete(models.ProductCategory{}).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("Error Service When Delete Product Category : " + err.Error())
		}
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		log.Println("Error Service Save Product : ", err)
		return nil, err
	}

	if err := p.AddProductCategory(categories, product.Id); err != nil {
		tx.Rollback()
		return nil, err
	}

	go func() {
		err := helpers.DeleteAssetCloudinary(productBeforeUpdt.Picture)

		if err != nil {
			log.Println("Error Service When Delete Asset From Cloudinary : ", err)
			tx.Rollback()
			errChan <- err
			return
		}

		picUrl, err := helpers.UploadToCloudinary(pictureFile)

		if err != nil {
			log.Println("Service Error When Upload Picture To cloudinary : " + err.Error())
			tx.Rollback()
			errChan <- err
			return
		}

		errChan <- nil
		picUrlChan <- picUrl
	}()

	if err := <-errChan; err != nil {
		return nil, err
	}

	picUrl := <-picUrlChan

	if err := tx.Model(&product).Where("id = ?", product.Id).Update("picture", picUrl).Error; err != nil {
		log.Println("Error Repository When Update Product : " + err.Error())
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return product, nil

}
