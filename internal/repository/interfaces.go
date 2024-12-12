package repository

import (
	"context"
	"github.com/gracchi-stdio/barf/internal/domain"

	"gorm.io/gorm"
)

type Transaction interface {
	Commit() error
	Rollback() error
}

type BookRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	Create(ctx context.Context, tx *gorm.DB, book *domain.Book) error
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.Book, error)
	GetByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	List(ctx context.Context, limit, offset int) ([]domain.Book, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]domain.Book, int64, error)
}

type InventoryRepository interface {
	Create(ctx context.Context, tx *gorm.DB, inventory *domain.Inventory) error
	Update(ctx context.Context, inventory *domain.Inventory) error
	GetByBookID(ctx context.Context, bookID string) (*domain.Inventory, error)
	UpdateQuantity(ctx context.Context, bookID string, quantity int) error
	ListLowStock(ctx context.Context, threshold int) ([]domain.Inventory, error)
}
