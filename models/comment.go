package models

type Comment struct {
	BaseModel
	CardID  uint64
	UserID  uint64
	Content string `gorm:"type:varchar(1024)"`
	User    User   `gorm:"foreignKey:UserID"`
}
