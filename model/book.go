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

type BookStore struct {
	mu     sync.RWMutex
	books  map[int]Book
	lastID int
}

var (
	storeInstance *BookStore
	once          sync.Once
)

func ResetBookStore() {
	once = sync.Once{}
	storeInstance = &BookStore{
		books:  make(map[int]Book),
		lastID: 0,
	}
}

func GetBookStore() *BookStore {
	once.Do(func() {
		storeInstance = &BookStore{
			books:  make(map[int]Book),
			lastID: 0,
		}
	})

	return storeInstance
}

func (bs *BookStore) AddBook(b Book) Book {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.lastID++
	b.ID = bs.lastID
	bs.books[b.ID] = b
	return b
}

func (bs *BookStore) GetAllBooks() []Book {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	books := []Book{}
	for _, b := range bs.books {
		books = append(books, b)
	}
	return books
}

func (bs *BookStore) GetBookByID(id int) (Book, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	b, ok := bs.books[id]
	if !ok {
		return Book{}, errors.New("book not found")
	}
	return b, nil
}

func (bs *BookStore) UpdateBook(id int, updated Book) (Book, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if _, ok := bs.books[id]; !ok {
		return Book{}, errors.New("book not found")
	}
	updated.ID = id
	bs.books[id] = updated
	return updated, nil
}

func (bs *BookStore) DeleteBook(id int) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if _, ok := bs.books[id]; !ok {
		return errors.New("book not found")
	}
	delete(bs.books, id)
	return nil
}
