package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"products-api-with-jwt/models"
	"products-api-with-jwt/services"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	BookService *services.BookService
	AuthService *services.AuthService
}

// NewBookController menginisialisasi BookController baru
func NewBookController(bookService *services.BookService, authService *services.AuthService) *BookController {
	return &BookController{BookService: bookService, AuthService: authService}
}

// Security definition for Bearer token
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// GetBooks godoc
// @Summary Get all books
// @Description Get a list of all books
// @Tags books
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /books [get]
func (pc *BookController) GetBooks(c *gin.Context) {
	books, err := pc.BookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not retrieve products",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Books retrieved successfully",
		Data:    books,
		Count:   len(books), // Optional count of items
	})
}

// GetBookByID godoc
// @Summary Get product by ID
// @Description Get details of a product by its ID
// @Tags books
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 404 {object} models.ApiResponse
// @Router /books/{id} [get]
func (pc *BookController) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid product ID",
			Data:    nil,
		})
		return
	}

	book, err := pc.BookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusNotFound,
			Message: "Book not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Book retrieved successfully",
		Data:    book,
	})
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with the given details
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param book body models.Book true "Book"
// @Success 201 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /books [post]
func (pc *BookController) CreateBook(c *gin.Context) {
	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	book, err := pc.BookService.CreateBook(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not create book",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: "Book created successfully",
		Data:    book,
	})
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update a book's information by its ID
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Book"
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 404 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /books/{id} [put]
func (pc *BookController) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid book ID",
			Data:    nil,
		})
		return
	}

	var input models.Book
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	updatedBook, err := pc.BookService.UpdateBook(id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not update book",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Book updated successfully",
		Data:    updatedBook,
	})
}

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Delete a book by its ID
// @Tags books
// @Security BearerAuth
// @Param id path int true "Book ID"
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /books/{id} [delete]
func (pc *BookController) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid book ID",
			Data:    nil,
		})
		return
	}

	if err := pc.BookService.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not delete book",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Book deleted successfully",
		Data:    nil,
	})
}

// BorrowBook godoc
// @Summary Borrow a book
// @Description Borrow a book by its ID for the authenticated user
// @Tags books
// @Security BearerAuth
// @Param bookId path int true "Book ID"
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 401 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /books/borrow/{bookId} [post]
func (pc *BookController) BorrowBook(c *gin.Context) {
	// Ambil token dari header Authorization
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Authorization token is required",
			Data:    nil,
		})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	// Dapatkan userID dari token
	userID, err := pc.AuthService.GetUserIDFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Ambil bookId dari path parameter
	bookId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid book ID",
			Data:    nil,
		})
		return
	}

	// Panggil service untuk meminjam buku
	if err := pc.BookService.BorrowBook(int(userID), bookId); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Book borrowed successfully",
		Data:    nil,
	})
}

// ReturnBook godoc
// @Summary Return a borrowed book
// @Description Return the borrowed book for the authenticated user
// @Tags books
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 401 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /books/return [post]
func (pc *BookController) ReturnBook(c *gin.Context) {
	// Ambil token dari header Authorization
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Authorization token is required",
			Data:    nil,
		})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	// Dapatkan userID dari token
	userID, err := pc.AuthService.GetUserIDFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	//Dapatkan id dari POST
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid book ID",
			Data:    nil,
		})
		return
	}

	// Panggil service untuk mengembalikan buku
	if err := pc.BookService.ReturnBook(int(userID), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Book returned successfully",
		Data:    nil,
	})
}
