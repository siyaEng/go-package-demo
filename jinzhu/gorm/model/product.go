package model

import (
	"github.com/jinzhu/gorm"
	"go-package-demo/jinzhu/gorm/mysql"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func (product *Product) Migrade() {
	// Migrate the schema
	mysql.MysqlPool.AutoMigrate(&Product{})
}

func (product *Product) Create() {
	// Create
	mysql.MysqlPool.Create(&Product{Code: "L1212", Price: 1000})
}

func (product *Product) First() {
	// Read
	mysql.MysqlPool.First(&product, 1)                   // find product with id 1
	mysql.MysqlPool.First(&product, "code = ?", "L1212") // find product with code l1212
}

func (product *Product) Update() {
	// Update - update product's price to 2000
	mysql.MysqlPool.Model(&product).Update("Price", 2000)
}

func (product *Product) Delete() {
	// Delete - delete product
	mysql.MysqlPool.Delete(&product)
}
