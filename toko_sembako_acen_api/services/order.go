package services

import (
	"errors"
	"log"
	"toko_sembako_acen/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (o *OrderService) GetOrders() ([]*models.Order, error) {
	var orders []*models.Order

	if err := o.db.Select("Distinct").Find(&orders).Error; err != nil {
		log.Println("Error Get order By id : " + err.Error())
		return nil, err
	}

	return orders, nil
}

func (o *OrderService) GetOrderItemById(orderId string) ([]*models.OrderItem, error) {
	var orderItem []*models.OrderItem

	if err := o.db.Preload("Product").Preload("User").Where("order_id = ?", orderId).Find(&orderItem).Error; err != nil {
		log.Println("Error Get order By id : " + err.Error())
		return nil, err
	}

	return orderItem, nil
}

func (o *OrderService) CreateOrderItems(userId string) (*[]models.OrderItem, error) {
	var order models.Order
	var orderItem models.OrderItem
	var orderItems []models.OrderItem
	var totalNetIncome, totalPrice float64

	var cartItems []models.CartItem

	if err := o.db.Where("user_id = ?", userId).Preload("User").Preload("Product").Find(&cartItems).Error; err != nil {
		log.Println("Service Error When getting Cart Items : " + err.Error())
		return nil, err
	}

	tx := o.db.Begin()

	order.UserID = uuid.MustParse(userId)

	if err := tx.Create(&order).Error; err != nil {
		log.Println("Service Error Creating Order : " + err.Error())
		return nil, err
	}

	for _, cartItem := range cartItems {

		if cartItem.Qty > cartItem.Product.Stock {
			return nil, errors.New("Stock is Not Enough")
		}

		if cartItem.Price < cartItem.Product.Capital {
			return nil, errors.New("Price Is invalid ")
		}

		hasilSubnet := cartItem.SubTotal - (cartItem.Product.Capital * float64(cartItem.Qty))

		orderItem.AdminName = cartItem.User.Username
		orderItem.ProductID = cartItem.ProductID
		orderItem.UserID = cartItem.UserID
		orderItem.OrderID = order.Id
		orderItem.SubNetIncome = hasilSubnet
		orderItem.Sub_total = cartItem.SubTotal
		orderItem.Qty = cartItem.Qty
		orderItem.Price = cartItem.Price

		if err := tx.Create(&orderItem).Error; err != nil {
			log.Println("Error creating order Item : " + err.Error())
			tx.Rollback()
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
		orderItem.Id = uuid.Nil

		totalNetIncome += hasilSubnet
		totalPrice += cartItem.SubTotal

		if err := tx.Delete(&cartItem).Error; err != nil {
			log.Println("Error deleting cart item : " + err.Error())
			tx.Rollback()
			return nil, err
		}

		if err := tx.Where("id = ? ", cartItem.Product.Id).
			Model(&models.Product{}).
			Update("stock", cartItem.Product.Stock-orderItem.Qty).
			Error; err != nil {
			log.Println("Error When Updating Product Stock : " + err.Error())
			return nil, err
		}

	}

	tx.Commit()

	if err := o.db.Where("id = ?", order.Id).Updates(&models.Order{
		UserID:         orderItem.UserID,
		TotalNetIncome: &totalNetIncome,
		TotalPrice:     &totalPrice,
	}).Error; err != nil {
		return nil, err
	}

	return &orderItems, nil
}
