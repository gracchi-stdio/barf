package service

import (
	"context"
	"github.com/gracchi-stdio/barf/internal/domain"
	"github.com/gracchi-stdio/barf/internal/repository"
	domainErr "github.com/gracchi-stdio/barf/pkg/errors"
)

type BookService struct {
	bookRepo      repository.BookRepository
	inventoryRepo repository.InventoryRepository
}

func NewBookService(bookRepo repository.BookRepository, inventoryRepo repository.InventoryRepository) *BookService {
	return &BookService{
		bookRepo:      bookRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *BookService) CreateBook(ctx context.Context, book *domain.Book, initialQuantity int, price float64) error {
	// start transaction
	tx, err := s.bookRepo.BeginTx(ctx)
	if err != nil {
		return err
	}

	// create book
	if err := s.bookRepo.Create(ctx, tx, book); err != nil {
		tx.Rollback()
		return err
	}

	inventory := &domain.Inventory{
		BookID:   book.ID,
		Quantity: initialQuantity,
		Price:    price,
	}
	if err := s.inventoryRepo.Create(ctx, tx, inventory); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *BookService) UpdataBook(ctx context.Context, book *domain.Book) error {
	return s.bookRepo.Update(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id string) error {
	return s.bookRepo.Delete(ctx, id)
}

func (s *BookService) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id)
}

func (s *BookService) SearchBook(ctx context.Context, query string, page, pageSize int) ([]domain.Book, int64, error) {
	offset := (page - 1) * pageSize
	return s.bookRepo.Search(ctx, query, offset, pageSize)
}

func (s *BookService) UpdateInventory(ctx context.Context, bookID string, quantityChange int) error {
	// verify the book exists
	_, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return err
	}

	// get current inventory
	inventory, err := s.inventoryRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return err
	}

	// check if we have enough stock to sell
	if quantityChange < 0 && (inventory.Quantity+quantityChange) < 0 {
		return domainErr.ErrInsufficientStock
	}

	return s.inventoryRepo.UpdateQuantity(ctx, bookID, quantityChange)
}

func (s *BookService) GetInventory(ctx context.Context, bookID string) (*domain.Inventory, error) {
	return s.inventoryRepo.GetByBookID(ctx, bookID)
}

func (s *BookService) GetLowStockBooks(ctx context.Context, threshold int) ([]domain.Inventory, error) {
	return s.inventoryRepo.ListLowStock(ctx, threshold)
}
