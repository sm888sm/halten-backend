package models

import (
	"fmt"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	// Create ENUM types
	if result := db.Exec("DO $$ BEGIN CREATE TYPE role_enum AS ENUM ('owner', 'admin', 'member', 'observer'); EXCEPTION WHEN duplicate_object THEN null; END $$;"); result.Error != nil {
		fmt.Println("Error creating role_enum:", result.Error)
	}
	if result := db.Exec("DO $$ BEGIN CREATE TYPE type_enum AS ENUM ('image', 'video', 'document'); EXCEPTION WHEN duplicate_object THEN null; END $$;"); result.Error != nil {
		fmt.Println("Error creating type_enum:", result.Error)
	}
	if result := db.Exec("DO $$ BEGIN CREATE TYPE visibility_enum AS ENUM ('public', 'private'); EXCEPTION WHEN duplicate_object THEN null; END $$;"); result.Error != nil {
		fmt.Println("Error creating visibility_enum:", result.Error)
	}

	// Migrate the schema
	db.AutoMigrate(
		&ActivityLog{},
		&User{},
		&Board{},
		&List{},
		&Card{},
		&Attachment{},
		&Label{},
		&Notification{},
		&Permission{},
		&Watch{},
	)
}
