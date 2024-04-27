package models

type List struct {
	BaseModel
	BoardID    uint64 `gorm:"foreign_key"`
	Name       string `gorm:"type:varchar(50)"`
	Position   int64
	IsArchived bool
	Cards      []Card  `gorm:"foreignKey:ListID"`
	Watches    []Watch `gorm:"foreignKey:ListID"`
}
