// FILE: label.go
package models

type Label struct {
	BaseModel
	BoardID uint64 `gorm:"not null"`
	Name    string `gorm:"type:varchar(50);not null"`
	Color   string `gorm:"type:varchar(7);not null"`
	Cards   []Card `gorm:"many2many:card_labels"` // New field
}
