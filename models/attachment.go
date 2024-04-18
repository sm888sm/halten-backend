package models

type Attachment struct {
	BaseModel
	BoardID   uint64 `gorm:"foreignKey:CardID"`
	CardID    uint64 `gorm:"foreignKey:CardID"`
	FileName  string
	FilePath  string
	Type      string `gorm:"type:type_enum;default:'document'"`
	Thumbnail string
}
