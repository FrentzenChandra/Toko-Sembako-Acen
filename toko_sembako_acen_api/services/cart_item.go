package services

import (
	"errors"
	"log"
	"toko_sembako_acen/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItemService struct {
	db *gorm.DB
}

func NewCartItemService(db *gorm.DB) *CartItemService {
	return &CartItemService{db: db}
}

func (c *CartItemService) AddCartItem(input *models.CartItemInput, userId uuid.UUID) error {
	var product *models.Product
	var cartItem *models.CartItem
	if input.Qty == 0 {
		input.Qty = 1
	}

	if err := c.db.Where("id = ? ", input.ProductID).First(&product).Error; err != nil {
		log.Println("Service Erorr finding Product with id : " + err.Error())

		if err == gorm.ErrRecordNotFound {
			return errors.New("Product Not found")
		}

		return err
	}

	if input.Price == 0 {
		input.Price = product.Price
	}



	if err := c.db.Where("product_id = ? AND user_id = ?", input.ProductID, userId).First(&cartItem).Error; err != nil {
		

		if err := c.db.Create(&models.CartItem{
			ProductID: input.ProductID,
			UserID:    userId,
			Qty:       input.Qty,
			Price:     input.Price,
			SubTotal:  input.Price * float64(input.Qty),
		}).Error; err != nil {
			log.Println("Service Error Creating Cart Item : " + err.Error())
			return err
		}
	} else {
		if err := c.db.Where("id = ?", cartItem.Id).
			Updates(models.CartItem{Qty: input.Qty + cartItem.Qty, SubTotal: float64(cartItem.Qty+input.Qty) * input.Price}).
			Error; err != nil {
			log.Println("Service Error Updating Cart Item : " + err.Error())
			return err
		}
	}

	return nil
}
