package http

import (
	"github.com/gracchi-stdio/barf/internal/domain"
	"github.com/gracchi-stdio/barf/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type BookHandler struct {
	BookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		BookService: bookService,
	}
}

type CreateBookRequest struct {
	Title           string  `json:"title"`
	ISBN            string  `json:"isbn"`
	Author          string  `json:"author"`
	Publisher       string  `json:"publisher"`
	PublicationDate string  `json:"publication_date"`
	InitialQuantity int     `json:"initial_quantity"`
	Price           float64 `json:"price"`
}

func (h *BookHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/v1/books", h.CreateBook)
	e.GET("/api/v1/books/:id", h.GetBook)
	e.GET("/api/v1/books", h.SearchBook)
	e.GET("/api/v1/books/low-stock", h.GetLowStockBooks)
	e.GET("/api/v1/books/:id/inventory", h.GetInventory)
	e.PUT("/api/v1/books/:id/inventory", h.UpdateInventory)
	e.POST("/api/v1/books", h.CreateBook)
	e.PUT("/api/v1/books/:id", h.UpdateBook)
	e.DELETE("/api/v1/books/:id", h.DeleteBook)

}

func (h *BookHandler) CreateBook(c echo.Context) error {
	var req CreateBookRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	book := &domain.Book{
		Title:           req.Title,
		ISBN:            req.ISBN,
		Author:          req.Author,
		Publisher:       req.Publisher,
		PublicationDate: req.PublicationDate,
	}

	if err := h.BookService.CreateBook(c.Request().Context(), book, req.InitialQuantity, req.Price); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c echo.Context) error {
	var book domain.Book
	if err := c.Bind(&book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.BookService.UpdataBook(c.Request().Context(), &book); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c echo.Context) error {
	id := c.Param("id")

	if err := h.BookService.DeleteBook(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *BookHandler) GetBook(c echo.Context) error {
	id := c.Param("id")

	book, err := h.BookService.GetBookByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) SearchBook(c echo.Context) error {
	query := c.QueryParam("q")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	books, total, err := h.BookService.SearchBook(c.Request().Context(), query, page, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"books": books,
		"total": total,
		"page":  page,
	})
}

type UpdateInventoryRequest struct {
	QuantityChange int `json:"quantity_change"`
}

func (h *BookHandler) UpdateInventory(c echo.Context) error {
	id := c.Param("id")

	var req UpdateInventoryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.BookService.UpdateInventory(c.Request().Context(), id, req.QuantityChange); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *BookHandler) GetInventory(c echo.Context) error {
	id := c.Param("id")
	inventory, err := h.BookService.GetInventory(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, inventory)
}

func (h *BookHandler) GetLowStockBooks(c echo.Context) error {
	threshold, _ := strconv.Atoi(c.QueryParam("threshold"))
	if threshold < 1 {
		threshold = 1
	}
	books, err := h.BookService.GetLowStockBooks(c.Request().Context(), threshold)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}
