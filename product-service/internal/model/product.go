package model

//商品模型
type Product struct {
	ID        int64         `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string        `gorm:"column:name;type:varchar(255);not null"`
	Avatar    string        `gorm:"column:avatar;type:varchar(500)"`
	Category  string        `gorm:"column:category;type:varchar(100);not null"`
	Price     float64       `gorm:"column:price;type:decimal(10,2);not null"`
	Stock     int32         `gorm:"column:stock;not null;default:0"`
	Status    ProductStatus `gorm:"column:status;not null;default:0"`
	CreatedAt int64         `gorm:"column:created_at;type:bigint;not null"`
	UpdatedAt int64         `gorm:"column:updated_at;type:bigint;not null"`
	Brand     string        `gorm:"column:brand;type:varchar(100);default:''"`
}

//表名
func (Product) TableName() string {
	return "products"
}
