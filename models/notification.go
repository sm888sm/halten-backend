package models

type Notification struct {
	BaseModel
	ActivityLogID uint64
	UserID        uint64
	IsRead        bool
}
