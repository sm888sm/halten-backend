package models

type ActivityLog struct {
	BaseModel
	BoardID       uint64
	UserID        uint64
	ActionType    string `gorm:"type:varchar(50)"`
	Details       string `gorm:"type:text"`
	Notifications []Notification
}
