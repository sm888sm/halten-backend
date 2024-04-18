package models

type BoardMember struct {
	BaseModel
	BoardID uint64 `gorm:"uniqueIndex:user_board_idx"`
	UserID  uint64 `gorm:"uniqueIndex:user_board_idx"`
	Role    string `gorm:"type:role_enum;default:'observer'"`
}
