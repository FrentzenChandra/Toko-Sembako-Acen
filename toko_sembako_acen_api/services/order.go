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

func (o *OrderService) GetOrders(orderBy string) (*[]models.OrdersRequest, error) {
	var orders []models.OrdersRequest

	err := o.db.Raw("SELECT o.id , o.user_id , o.total_net_income , o.total_price , (SELECT oi2.product_id FROM order_item oi2 WHERE oi2.order_id = o.id limit 1) , o.created_at, COUNT(oi.id) items_count FROM \"order\" o INNER JOIN order_item oi ON o.id  = oi.order_id GROUP BY o.id ORDER BY o.created_at " + orderBy).
		Scan(&orders).
		Error

	for index, OrdersRequest := range orders {

		var productData models.Product
		var userData models.Users

		if err := o.db.Where("id = ? ", OrdersRequest.ProductID).First(&productData).Error; err != nil {
			log.Println("Error service find Product : " + err.Error())
			return nil, err
		}

		if err := o.db.Where("id = ? ", OrdersRequest.UserID).First(&userData).Error; err != nil {
			log.Println("Error service find User : " + err.Error())
			return nil, err
		}

		orders[index].Product = productData
		orders[index].User = userData

	}

	if err != nil {
		log.Println("Error Service Order Find Order : " + err.Error())
		return nil, err
	}

	return &orders, nil
}

func (o *OrderService) GetOrderItemsById(orderId string) ([]*models.OrderItem, error) {
	var orderItem []*models.OrderItem

	rowAffected := o.db.Where("order_id = ?", orderId).Find(&orderItem).RowsAffected

	if rowAffected == 0 {
		return nil, errors.New("Order Data not found")
	}

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
