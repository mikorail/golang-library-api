package services

import (
	"errors"
	"fmt"
	"products-api-with-jwt/models"
	"time"

	"gorm.io/gorm"
)

type BookService struct {
	DB *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
	return &BookService{DB: db}
}

// GetAllBooks mengambil semua buku dari database
func (s *BookService) GetAllBooks() ([]models.Book, error) {
	var Books []models.Book
	if err := s.DB.Find(&Books).Error; err != nil {
		return nil, err
	}
	return Books, nil
}

// GetBookByID mengambil buku berdasarkan ID
func (s *BookService) GetBookByID(id int) (*models.Book, error) {
	var Book models.Book
	if err := s.DB.First(&Book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Book Not found")
		}
		return nil, err
	}
	return &Book, nil
}

// CreateBook menambah buku baru ke database
func (s *BookService) CreateBook(Book *models.Book) (models.Book, error) {
	// Menyimpan buku baru ke database
	if err := s.DB.Create(Book).Error; err != nil {
		return models.Book{}, err // Kembalikan error jika terjadi kesalahan
	}
	return *Book, nil // Kembalikan buku yang baru dibuat
}

// UpdateBook memperbarui hanya field yang disediakan dalam permintaan
func (s *BookService) UpdateBook(id int, updatedBook *models.Book) (*models.Book, error) {
	var Book models.Book

	// Cari buku berdasarkan ID
	if err := s.DB.First(&Book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Book Not Found")
		}
		return nil, err
	}

	// Perbarui hanya field yang disediakan
	if updatedBook.Title != "" {
		Book.Title = updatedBook.Title
	}
	if updatedBook.Description != "" {
		Book.Description = updatedBook.Description
	}
	if updatedBook.Stock != 0 {
		Book.Stock = updatedBook.Stock
	}
	if updatedBook.Author != "" {
		Book.Author = updatedBook.Author
	}

	// Simpan perubahan ke database
	if err := s.DB.Save(&Book).Error; err != nil {
		return nil, err
	}
	return &Book, nil
}

// DeleteBook menghapus buku berdasarkan ID
func (s *BookService) DeleteBook(id int) error {
	if err := s.DB.Delete(&models.Book{}, id).Error; err != nil {
		return err
	}
	return nil
}

// BorrowBook memperbarui data pengguna dengan ID buku yang dipinjam dan mengurangi stok buku
func (s *BookService) BorrowBook(userId, bookId int) error {
	var user models.User
	var book models.Book

	if err := s.DB.First(&user, userId).Error; err != nil {
		return err
	}

	if user.BookBorrowed != 0 {
		return fmt.Errorf("User Already Has a Borrowed Book")
	}

	if err := s.DB.First(&book, bookId).Error; err != nil {
		return fmt.Errorf("Book Not Found")
	}

	if book.Stock <= 0 {
		return fmt.Errorf("Book Is Out Of Stock")
	}

	user.BookBorrowed = bookId
	book.Stock -= 1
	book.Borrowed += 1
	now := time.Now()
	user.BorrowDate = &now

	if err := s.DB.Save(&user).Error; err != nil {
		return err
	}
	if err := s.DB.Save(&book).Error; err != nil {
		return err
	}

	return nil
}

// BorrowBook memperbarui data pengguna dengan ID buku yang dipinjam dan mengurangi stok buku
func (s *BookService) ReturnBook(userId, bookId int) error {
	var user models.User
	var book models.Book

	if err := s.DB.First(&user, userId).Error; err != nil {
		return err
	}

	if user.BookBorrowed != bookId {
		return fmt.Errorf("Invalid Book Returned")
	}

	if user.BookBorrowed == 0 {
		return fmt.Errorf("User Has Not Borrowed A Book")
	}

	if err := s.DB.First(&book, bookId).Error; err != nil {
		return fmt.Errorf("Book Not Found")
	}

	user.BookBorrowed = 0
	book.Stock += 1
	book.Borrowed -= 1
	user.BorrowDate = nil

	if err := s.DB.Save(&user).Error; err != nil {
		return err
	}
	if err := s.DB.Save(&book).Error; err != nil {
		return err
	}

	return nil
}
