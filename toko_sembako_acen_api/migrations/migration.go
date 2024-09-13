package migrations

import (
	"toko_sembako_acen/infra/database"
	"toko_sembako_acen/models"
)

// Migrate Add list of model add for migrations
// TODO later separate migration each models
func Migrate() {
	var migrationModels = []interface{}{
		&models.Users{},
		&models.CartItem{},
		&models.Product{},
		&models.ProductCategory{},
		&models.Category{},
		&models.Order{},
		&models.OrderItem{},
	}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
