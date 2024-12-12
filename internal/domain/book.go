package domain

import (
	"github.com/google/uuid"
	"time"
)

type Book struct {
	ID              uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Title           string     `json:"title" gorm:"not null"`
	ISBN            string     `json:"isbn" gorm:"not null"`
	Author          string     `json:"author" gorm:"not null"`
	Publisher       string     `json:"publisher"`
	PublicationDate string     `json:"publication_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

type Inventory struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	BookID    uuid.UUID  `json:"book_id" gorm:"type:uuid;not null"`
	Book      Book       `gorm:"foreignkey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity  int        `json:"quantity" gorm:"not null"`
	Price     float64    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
