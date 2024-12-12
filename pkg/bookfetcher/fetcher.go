package bookfetcher

import (
	"context"
	"errors"
)

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrInvalidISBN       = errors.New("invalid isbn")
	ErrProviderError     = errors.New("provide error")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

type BookInfo struct {
	Title           string   `json:"title"`
	ISBN            string   `json:"isbn"`
	ISBN13          string   `json:"isbn_13,omitempty"`
	ISBN10          string   `json:"isbn_10,omitempty"`
	Authors         []string `json:"author"`
	Publisher       string   `json:"publisher"`
	PublicationDate string   `json:"publication_date"`
	Description     string   `json:"description,omitempty"`
	PageCount       int      `json:"page_count,omitempty"`
	Categories      []string `json:"category,omitempty"`
	Language        string   `json:"language,omitempty"`
	PreviewLink     string   `json:"preview_link,omitempty"`
	ThumbnailURL    string   `json:"thumbnail_url,omitempty"`
	Provider        string   `json:"provider"`
	ProviderID      string   `json:"provider_id"`
	RawData         any      `json:"raw_data"`
}

// SearchResult represents a search response from providers
type SearchResult struct {
	Books        []BookInfo `json:"books"`
	TotalResults int        `json:"total_results"`
	ItemsPerPage int        `json:"items_per_page"`
	StartIndex   int        `json:"start_index"`
	HasMore      bool       `json:"has_more"`
}

type SearchOptions struct {
	Query      string
	MaxResults int
	StartIndex int
	Language   string
	PrintType  string
	OrderBy    string
}

type BookFetcher interface {
	// GetBookByISBN returns a book by its ISBN
	GetBookByISBN(ctx context.Context, isbn string) (*BookInfo, error)

	// SearchBooks searches for books based on the provided options
	SearchBooks(ctx context.Context, opts SearchOptions) (*SearchResult, error)

	// Name returns the name of the book fetcher
	Name() string

	// IsHealthy returns true if the book fetcher is healthy
	IsHealthy(ctx context.Context) bool
}

type FetcherFactory interface {
	CreateFetcher(config map[string]string) (BookFetcher, error)
}
