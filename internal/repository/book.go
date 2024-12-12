package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gracchi-stdio/barf/internal/domain"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) Create(ctx context.Context, tx *gorm.DB, book *domain.Book) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	result := r.db.WithContext(ctx).Save(book)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *bookRepository) Delete(ctx context.Context, id string) error {
	bookID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	result := r.db.WithContext(ctx).Delete(&domain.Book{}, bookID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *bookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	bookID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var book domain.Book
	result := r.db.WithContext(ctx).First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func (r *bookRepository) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	var book domain.Book
	result := r.db.WithContext(ctx).Where("isbn = ?", isbn).First(&book)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &book, nil
}

func (r *bookRepository) List(ctx context.Context, limit, offset int) ([]domain.Book, int64, error) {
	var books []domain.Book
	var count int64

	// get total count
	if err := r.db.WithContext(ctx).Model(&domain.Book{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&books)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return books, count, nil
}

func (r *bookRepository) Search(ctx context.Context, query string, offset, limit int) ([]domain.Book, int64, error) {
	var books []domain.Book
	var count int64

	searchQuery := "%" + query + "%"
	baseQuery := r.db.WithContext(ctx).Where("title ILIKE ? OR author ILIKE ? OR isbn LIKE ?",
		searchQuery, searchQuery, searchQuery)

	if err := baseQuery.Model(&domain.Book{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	result := baseQuery.Limit(limit).Offset(offset).Order("created_at DESC").Find(&books)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return books, count, nil
}

func (r *bookRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
