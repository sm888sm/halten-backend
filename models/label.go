// FILE: label.go
package models

import (
	"gorm.io/gorm"
)

type Label struct {
	gorm.Model
	Name    string
	Color   string
	BoardID uint
}
