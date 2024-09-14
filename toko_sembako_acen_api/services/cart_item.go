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

func (c *CartItemService) GetCartItems(userId uuid.UUID) ([]*models.CartItem, error) {
	var cartItems []*models.CartItem

	if err := c.db.Preload("Product").Where("user_id = ?", userId).Find(&cartItems).Error; err != nil {
		log.Println("Service Error When getting Cart Items : " + err.Error())
		return nil, err
	}

	return cartItems, nil

}

func (c *CartItemService) AddCartItem(input *models.CartItemInput, userId uuid.UUID) error {
	var product *models.Product
	var cartItem *models.CartItem

	if input.ProductID == nil {
		return errors.New("Invalid Input")
	}

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

	if input.Price < product.Capital {
		return errors.New("Price Is invalid ")
	}

	if input.Price == 0 {
		input.Price = product.Price
	}

	if err := c.db.Where("product_id = ? AND user_id = ?", input.ProductID, userId).First(&cartItem).Error; err != nil {

		if input.Qty > product.Stock {
			return errors.New("Not enough stock")
		}

		if err := c.db.Create(&models.CartItem{
			ProductID: *input.ProductID,
			UserID:    userId,
			Qty:       input.Qty,
			Price:     input.Price,
			SubTotal:  input.Price * float64(input.Qty),
		}).Error; err != nil {
			log.Println("Service Error Creating Cart Item : " + err.Error())
			return err
		}
	} else {
		hasilQty := input.Qty + cartItem.Qty
		hasilSubTotal := float64(hasilQty) * input.Price

		if hasilQty > product.Stock {
			return errors.New("Not Enough Stock")
		}

		if err := c.db.Where("id = ?", cartItem.Id).
			Updates(models.CartItem{Qty: hasilQty, SubTotal: hasilSubTotal, Price: input.Price}).
			Error; err != nil {
			log.Println("Service Error Updating Cart Item : " + err.Error())
			return err
		}
	}

	return nil
}

func (c *CartItemService) UpdateCartItem(input *models.CartItemInput, userId uuid.UUID, productId string) error {
	var cartItem *models.CartItem
	var product *models.Product

	if input.Qty == 0 {
		input.Qty = 1
	}

	if err := c.db.Where("id = ? ", productId).First(&product).Error; err != nil {
		log.Println("Service Erorr finding Product with id : " + err.Error())

		if err == gorm.ErrRecordNotFound {
			return errors.New("Product Not found")
		}

		return err
	}

	if input.Price == 0 {
		input.Price = product.Price
	}

	if err := c.db.Where("user_id = ? AND product_id = ?", userId, productId).First(&cartItem).Error; err != nil {
		log.Println("Service Error When getting Cart Item : " + err.Error())
		return err
	}

	hasilQty := input.Qty
	hasilSubTotal := float64(input.Qty) * input.Price

	if hasilQty > product.Stock {
		return errors.New("Not Enough Stock")
	}

	if err := c.db.Where("id = ?", cartItem.Id).
		Updates(models.CartItem{Qty: hasilQty, SubTotal: hasilSubTotal, Price: input.Price}).
		Error; err != nil {
		log.Println("Service Error Updating Cart Item : " + err.Error())
		return err
	}

	return nil
}
