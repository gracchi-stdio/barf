package worldcat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gracchi-stdio/barf/pkg/bookfetcher"
	"net/http"
	"time"
)

const (
	baseURL      = "https://www.worldcat.org/webservices/catalog/content"
	providerName = "worldCat"
)

type WorldCatProvider struct {
	apiKey     string
	httpClient *http.Client
}

func (w WorldCatProvider) GetBookByISBN(ctx context.Context, isbn string) (*bookfetcher.BookInfo, error) {
	url := fmt.Sprintf("%s/isbn/%s?wskey=%s&format=json", baseURL, isbn, w.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
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

	var wcBook struct {
		Title       string   `json:"title"`
		ISBN        []string `json:"isbn"`
		Author      []string `json:"author"`
		Publisher   string   `json:"publisher"`
		PublishDate string   `json:"publishDate"`
		Description string   `json:"summary"`
		PageCount   int      `json:"pageCount"`
		ID          string   `json:"id"`
		Language    string   `json:"language"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&wcBook); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &bookfetcher.BookInfo{
		Title:           wcBook.Title,
		ISBN:            isbn,
		Authors:         wcBook.Author,
		Publisher:       wcBook.Publisher,
		PublicationDate: wcBook.PublishDate,
		Description:     wcBook.Description,
		PageCount:       wcBook.PageCount,
		Language:        wcBook.Language,
		Provider:        providerName,
		ProviderID:      wcBook.ID,
		RawData:         wcBook,
	}, nil
}

func (w WorldCatProvider) SearchBooks(ctx context.Context, opts bookfetcher.SearchOptions) (*bookfetcher.SearchResult, error) {
	url := fmt.Sprintf(
		"%s/search?wskey=%s&format=json&limit=%d&start=%d",
		baseURL,
		w.apiKey,
		opts.MaxResults,
		opts.StartIndex,
	)

	if opts.Language != "" {
		url += fmt.Sprintf("&language=%s", opts.Language)
	}
	if opts.PrintType != "" {
		url += fmt.Sprintf("&printType=%s", opts.PrintType)
	}
	if opts.OrderBy != "" {
		url += fmt.Sprintf("&orderBy=%s", opts.OrderBy)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", bookfetcher.ErrProviderError, resp.StatusCode)
	}

	var result struct {
		Books []struct {
			Title       string   `json:"title"`
			ISBN        []string `json:"isbn"`
			Author      []string `json:"author"`
			Publisher   string   `json:"publisher"`
			PublishDate string   `json:"publishDate"`
			ID          string   `json:"id"`
		} `json:"books"`
		TotalResults int `json:"total_results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	books := make([]bookfetcher.BookInfo, len(result.Books))
	for i, book := range result.Books {
		isbn := ""
		if len(book.ISBN) > 0 {
			isbn = book.ISBN[0]
		}
		books[i] = bookfetcher.BookInfo{
			Title:           book.Title,
			ISBN:            isbn,
			Authors:         book.Author,
			Publisher:       book.Publisher,
			PublicationDate: book.PublishDate,
			Provider:        providerName,
			ProviderID:      book.ID,
			RawData:         book,
		}
	}

	return &bookfetcher.SearchResult{
		Books:        books,
		TotalResults: result.TotalResults,
		ItemsPerPage: opts.MaxResults,
		StartIndex:   opts.StartIndex,
		HasMore:      result.TotalResults > opts.StartIndex+opts.MaxResults,
	}, nil
}

func (w WorldCatProvider) Name() string {
	return providerName
}

func (w WorldCatProvider) IsHealthy(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/health?wskey=%s", baseURL, w.apiKey), nil)
	if err != nil {
		return false
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

type WorldCatFactory struct{}

func (f WorldCatFactory) CreateFetcher(config map[string]string) (bookfetcher.BookFetcher, error) {
	apiKey, ok := config["apiKey"]
	if !ok {
		return nil, fmt.Errorf("apiKey not found in config")
	}

	timeout := 10 * time.Second
	if timeoutStr, ok := config["timeout"]; ok {
		if t, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = t
		}
	}

	return NewWorldCatProvider(apiKey, timeout), nil
}

func NewWorldCatProvider(apiKey string, timeout time.Duration) *WorldCatProvider {
	return &WorldCatProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: timeout},
	}
}
