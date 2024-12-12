package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gracchi-stdio/barf/internal/domain"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *inventoryRepository {
	return &inventoryRepository{
		db: db,
	}
}

func (i inventoryRepository) Create(ctx context.Context, tx *gorm.DB, inventory *domain.Inventory) error {
	db := i.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(inventory)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (i inventoryRepository) Update(ctx context.Context, inventory *domain.Inventory) error {
	result := i.db.WithContext(ctx).Save(inventory)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (i inventoryRepository) GetByBookID(ctx context.Context, bookID string) (*domain.Inventory, error) {
	id, err := uuid.Parse(bookID)
	if err != nil {
		return nil, err
	}

	var inventory domain.Inventory
	// preload book relations
	result := i.db.WithContext(ctx).Preload("Book").Where("book_id = ?", id).First(&inventory)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &inventory, nil
}

func (i inventoryRepository) UpdateQuantity(ctx context.Context, bookID string, quantity int) error {
	id, err := uuid.Parse(bookID)
	if err != nil {
		return err
	}

	result := i.db.WithContext(ctx).Model(&domain.Inventory{}).Where("book_id = ?", id).Update("quantity", gorm.Expr("quantity + ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (i inventoryRepository) ListLowStock(ctx context.Context, threshold int) ([]domain.Inventory, error) {
	var inventories []domain.Inventory
	result := i.db.WithContext(ctx).Joins("Book").Where("quantity <= ?", threshold).Find(&inventories)
	if result.Error != nil {
		return nil, result.Error
	}
	return inventories, nil
}
