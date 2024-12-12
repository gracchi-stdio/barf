package googlebooks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gracchi-stdio/barf/pkg/bookfetcher"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL      = "https://www.googleapis.com/books/v1/volumes"
	providerName = "googlebooks"
)

type GoogleBooksProvider struct {
	apiKey     string
	httpClient *http.Client
}

type GoogleBooksFactory struct{}

func NewGoogleBooksProvider(apiKey string, timeout time.Duration) *GoogleBooksProvider {
	return &GoogleBooksProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type volumeInfo struct {
	Title               string   `json:"title"`
	Authors             []string `json:"authors"`
	Publisher           string   `json:"publisher"`
	PublicationDate     string   `json:"publishedDate"`
	Description         string   `json:"description"`
	IndustryIdentifiers []struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"industryIdentifiers"`
	ImageLinks struct {
		SmallThumbnail string `json:"smallThumbnail"`
		Thumbnail      string `json:"thumbnail"`
	} `json:"imageLinks"`
	Language    string   `json:"language"`
	PreviewLink string   `json:"previewLink"`
	PageCount   int      `json:"pageCount"`
	Categories  []string `json:"categories"`
}

type googleBook struct {
	ID         string     `json:"id"`
	VolumeInfo volumeInfo `json:"volumeInfo"`
}

func (p *GoogleBooksProvider) GetBookByISBN(ctx context.Context, isbn string) (*bookfetcher.BookInfo, error) {
	// clean isbn
	isbn = strings.ReplaceAll(isbn, "-", "")

	// build url
	u := fmt.Sprintf("%s/volumes?q=isbn:%s", baseURL, url.QueryEscape(isbn))
	if p.apiKey != "" {
		u += fmt.Sprintf("&key=%s", p.apiKey) // if apiKey is set, use it
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// convert response to BookInfo
	case http.StatusNotFound:
		return nil, bookfetcher.ErrBookNotFound
	case http.StatusTooManyRequests:
		return nil, bookfetcher.ErrRateLimitExceeded
	default:
		return nil, fmt.Errorf("%w: status %d", bookfetcher.ErrProviderError, resp.StatusCode)
	}

	var result struct {
		TotalItems int          `json:"totalItems"`
		Items      []googleBook `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.TotalItems == 0 || len(result.Items) == 0 {
		return nil, bookfetcher.ErrBookNotFound
	}

	book := result.Items[0] // first result

	var isbn10, isbn13 string
	for _, identifier := range book.VolumeInfo.IndustryIdentifiers {
		if identifier.Type == "ISBN_10" {
			isbn10 = identifier.Identifier
		}
		if identifier.Type == "ISBN_13" {
			isbn13 = identifier.Identifier
		}
	}

	return &bookfetcher.BookInfo{
		Title:           book.VolumeInfo.Title,
		ISBN:            isbn,
		ISBN10:          isbn10,
		ISBN13:          isbn13,
		Authors:         book.VolumeInfo.Authors,
		Publisher:       book.VolumeInfo.Publisher,
		PublicationDate: book.VolumeInfo.PublicationDate,
		Description:     book.VolumeInfo.Description,
		PageCount:       book.VolumeInfo.PageCount,
		Categories:      book.VolumeInfo.Categories,
		Language:        book.VolumeInfo.Language,
		PreviewLink:     book.VolumeInfo.PreviewLink,
		ThumbnailURL:    book.VolumeInfo.ImageLinks.Thumbnail,
		Provider:        providerName,
		ProviderID:      book.ID,
		RawData:         book,
	}, nil
}

func (p *GoogleBooksProvider) SearchBooks(ctx context.Context, opts bookfetcher.SearchOptions) (*bookfetcher.SearchResult, error) {
	u := fmt.Sprintf("%s/volumes?q=%s", baseURL, url.QueryEscape(opts.Query))

	if opts.MaxResults > 0 {
		u += fmt.Sprintf("&maxResults=%d", opts.MaxResults)
	}
	if opts.StartIndex > 0 {
		u += fmt.Sprintf("&startIndex=%d", opts.StartIndex)
	}
	if opts.Language != "" {
		u += fmt.Sprintf("&langRestrict=%s", opts.Language)
	}
	if opts.PrintType != "" {
		u += fmt.Sprintf("&printType=%s", opts.PrintType)
	}
	if opts.OrderBy != "" {
		u += fmt.Sprintf("&orderBy=%s", opts.OrderBy)
	}
	if p.apiKey != "" {
		u += fmt.Sprintf("&key=%s", p.apiKey) // if apiKey is set, use it
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// convert response to BookInfo
	case http.StatusNotFound:
		return nil, bookfetcher.ErrBookNotFound
	case http.StatusTooManyRequests:
		return nil, bookfetcher.ErrRateLimitExceeded
	default:
		return nil, fmt.Errorf("%w: status %d", bookfetcher.ErrProviderError, resp.StatusCode)
	}

	var result struct {
		TotalItems int          `json:"totalItems"`
		Items      []googleBook `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	books := make([]bookfetcher.BookInfo, 0, len(result.Items))

	for _, item := range result.Items {
		var isbn10, isbn13 string
		for _, identifier := range item.VolumeInfo.IndustryIdentifiers {
			if identifier.Type == "ISBN_10" {
				isbn10 = identifier.Identifier
			}
			if identifier.Type == "ISBN_13" {
				isbn13 = identifier.Identifier
			}
		}

		books = append(books, bookfetcher.BookInfo{
			Title:           item.VolumeInfo.Title,
			ISBN10:          isbn10,
			ISBN13:          isbn13,
			Authors:         item.VolumeInfo.Authors,
			Publisher:       item.VolumeInfo.Publisher,
			PublicationDate: item.VolumeInfo.PublicationDate,
			Description:     item.VolumeInfo.Description,
			PageCount:       item.VolumeInfo.PageCount,
			Categories:      item.VolumeInfo.Categories,
			Language:        item.VolumeInfo.Language,
			PreviewLink:     item.VolumeInfo.PreviewLink,
			ThumbnailURL:    item.VolumeInfo.ImageLinks.Thumbnail,
			Provider:        providerName,
			ProviderID:      item.ID,
			RawData:         item,
		})
	}

	return &bookfetcher.SearchResult{
		Books:        books,
		TotalResults: result.TotalItems,
		ItemsPerPage: len(result.Items),
		StartIndex:   opts.StartIndex,
		HasMore:      opts.StartIndex+len(books) < result.TotalItems,
	}, nil
}

func (p *GoogleBooksProvider) Name() string {
	return providerName
}

func (p *GoogleBooksProvider) IsHealthy(ctx context.Context) bool {
	u := fmt.Sprintf("%s/volumes?q=test&maxResults=1", baseURL)
	if p.apiKey != "" {
		u += fmt.Sprintf("&key=%s", p.apiKey) // if apiKey is set, use it
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return false
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
