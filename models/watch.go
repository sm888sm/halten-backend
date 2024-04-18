package models

type Watch struct {
	BaseModel
	UserID  uint64
	BoardID *uint64
	ListID  *uint64
	CardID  *uint64
}
