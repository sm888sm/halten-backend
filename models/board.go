package models

type Board struct {
	BaseModel
	UserID     uint64
	Name       string `gorm:"type:varchar(50)"`
	Visibility string `gorm:"type:visibility_enum;default:'private'"`
	IsArchived bool
	Members    []BoardMember `gorm:"foreignKey:BoardID"`
	Lists      []List        `gorm:"foreignKey:BoardID"`
	Cards      []Card        `gorm:"foreignKey:BoardID"`
	Watches    []Watch       `gorm:"foreignKey:BoardID"`
	Labels     []Label       `gorm:"foreignKey:BoardID"`
}
