package service

import (
	"context"
	"errors"
	"github.com/gracchi-stdio/barf/internal/domain"
	"github.com/gracchi-stdio/barf/internal/repository"
	"github.com/gracchi-stdio/barf/pkg/bookfetcher"
	domainErr "github.com/gracchi-stdio/barf/pkg/errors"
	"strings"
)

type BookService struct {
	bookRepo       repository.BookRepository
	inventoryRepo  repository.InventoryRepository
	BookFetchers   map[string]bookfetcher.BookFetcher
	defaultFetcher string
}

func NewBookService(
	bookRepo repository.BookRepository,
	inventoryRepo repository.InventoryRepository,
	fetchers map[string]bookfetcher.BookFetcher,
	defaultFetcher string,
) *BookService {
	return &BookService{
		bookRepo:       bookRepo,
		inventoryRepo:  inventoryRepo,
		BookFetchers:   fetchers,
		defaultFetcher: defaultFetcher,
	}
}

func (s *BookService) FetchBookDetails(ctx context.Context, isbn string, provider string) (*bookfetcher.BookInfo, error) {
	var fetcher bookfetcher.BookFetcher
	var ok bool

	if provider != "" {
		provider = s.defaultFetcher
	}

	if fetcher, ok = s.BookFetchers[provider]; !ok {
		return nil, errors.New("provider not found")
	}

	return fetcher.GetBookByISBN(ctx, isbn)
}

func (s *BookService) CreateBookWithISBN(ctx context.Context, isbn string, initialQuantity int, price float64) (*domain.Book, error) {
	// first check if book exists
	existing, _ := s.bookRepo.GetByISBN(ctx, isbn)
	if existing != nil {
		return existing, nil
	}

	// fetch book details
	bookInfo, err := s.FetchBookDetails(ctx, isbn, "")
	if err != nil {
		return nil, err
	}

	authors := &strings.Builder{}
	for _, v := range m {
		authors.WriteString(v + "; ")
	}
	book := &domain.Book{
		Title:           bookInfo.Title,
		ISBN:            isbn,
		Author:          authors.String(),
		Publisher:       bookInfo.Publisher,
		PublicationDate: bookInfo.PublicationDate,
	}

	// start transaction
	tx, err := s.bookRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	// create book
	if err := s.bookRepo.Create(ctx, tx, book); err != nil {
		tx.Rollback()
		return nil, err
	}

	inventory := &domain.Inventory{
		BookID:   book.ID,
		Quantity: initialQuantity,
		Price:    price,
	}
	if err := s.inventoryRepo.Create(ctx, tx, inventory); err != nil {
		tx.Rollback()
		return nil, err
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return book, nil
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

func (s *BookService) UpdateBook(ctx context.Context, book *domain.Book) error {
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
