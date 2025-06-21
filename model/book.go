package model

import (
	"errors"
	"sync"
)

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
}

type BookStore interface {
	AddBook(book Book) Book
	GetAllBooks() []Book
	GetBookByID(id int) (Book, error)
	UpdateBook(id int, updated Book) (Book, error)
	DeleteBook(id int) error
}

type bookStore struct {
	mu     sync.RWMutex
	books  map[int]Book
	lastID int
}

// NewBookStore membuat instance BookStore baru dengan inisialisasi map dan ID terakhir.
func NewBookStore() BookStore {
	return &bookStore{
		books:  make(map[int]Book),
		lastID: 0,
	}
}

// AddBook menambahkan buku baru ke dalam store dan memberikan ID secara otomatis.
//
// Parameters:
//   - book: Book tanpa ID (akan diisi otomatis)
//
// Returns:
//   - Book yang sudah memiliki ID
func (bs *bookStore) AddBook(book Book) Book {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.lastID++
	book.ID = bs.lastID
	bs.books[book.ID] = book
	return book
}

// GetAllBooks mengembalikan semua buku dalam bentuk slice.
//
// Returns:
//   - Slice dari semua Book yang tersimpan
func (bs *bookStore) GetAllBooks() []Book {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	books := []Book{}
	for _, b := range bs.books {
		books = append(books, b)
	}
	return books
}

// GetBookByID mencari buku berdasarkan ID.
//
// Parameters:
//   - id: ID buku yang dicari
//
// Returns:
//   - Book jika ditemukan
//   - error jika tidak ditemukan
func (bs *bookStore) GetBookByID(id int) (Book, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	b, ok := bs.books[id]
	if !ok {
		return Book{}, errors.New("book not found")
	}
	return b, nil
}

// UpdateBook memperbarui data buku berdasarkan ID.
//
// Parameters:
//   - id: ID buku yang ingin diperbarui
//   - updated: data baru untuk buku
//
// Returns:
//   - Book hasil update
//   - error jika ID tidak ditemukan
func (bs *bookStore) UpdateBook(id int, updated Book) (Book, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if _, ok := bs.books[id]; !ok {
		return Book{}, errors.New("book not found")
	}
	updated.ID = id
	bs.books[id] = updated
	return updated, nil
}

// DeleteBook menghapus buku berdasarkan ID.
//
// Parameters:
//   - id: ID buku yang akan dihapus
//
// Returns:
//   - error jika ID tidak ditemukan
func (bs *bookStore) DeleteBook(id int) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if _, ok := bs.books[id]; !ok {
		return errors.New("book not found")
	}
	delete(bs.books, id)
	return nil
}
