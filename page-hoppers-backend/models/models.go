package models

import (
	"time"
	"gorm.io/gorm"
)

// User model - represents both parent and child users
type User struct {
	gorm.Model
	Name         string    `json:"name"` // Child's real name
	Age          int       `json:"age"`
	Password     string    `json:"-"` // Password hash, not exposed in JSON
	Email        string    `json:"email,omitempty" gorm:"uniqueIndex;default:null"`
	Role         string    `json:"role"` // "parent" or "child"
	ParentID     *uint     `json:"parent_id,omitempty"`
	Parent       *User     `json:"-" gorm:"foreignKey:ParentID"`
	Children     []User    `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	PIN          string    `json:"-"` // Optional PIN for child login
	LastLoginAt  time.Time `json:"last_login_at"`
	ReadingLogs  []ReadingLog `json:"reading_logs,omitempty" gorm:"foreignKey:ChildID"`
}

// ReadingLog model - represents a book reading activity by a child
type ReadingLog struct {
	gorm.Model
	Title       string    `json:"title"`
	Author      string    `json:"author,omitempty"`
	Status      string    `json:"status"` // "started" or "completed"
	Date        time.Time `json:"date"`
	ChildID     uint      `json:"child_id"`
	Child       User      `json:"-" gorm:"foreignKey:ChildID"`
	OpenLibraryKey string `json:"open_library_key,omitempty"` // For books found via Open Library API
	CoverID     *int      `json:"cover_id,omitempty"` // Open Library cover ID
}