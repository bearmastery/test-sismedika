package model

import (
	"strconv"
	"testing"
)

func setupStore() BookStore {
	return NewBookStore()
}

func createBook(store BookStore, title string, author string, year int) Book {
	book := Book{
		Title:         title,
		Author:        author,
		PublishedYear: year,
	}
	return store.AddBook(book)
}

func defaultBook(store BookStore) Book {
	return createBook(store, "Unit Testing in Go", "Go Dev", 2023)
}

func TestGetBookStoreEmpty(t *testing.T) {
	store := setupStore()

	if len(store.GetAllBooks()) != 0 {
		t.Error("Expected empty book list")
	}
}

func TestGetBookStore(t *testing.T) {
	store := setupStore()

	defaultBook(store)

	if len(store.GetAllBooks()) == 0 {
		t.Error("Expected non-empty book list")
	}
}

func TestCreateBook(t *testing.T) {
	store := setupStore()

	added := defaultBook(store)

	if added.ID == 0 {
		t.Error("Expected valid book ID")
	}
}

func TestGetBookByID(t *testing.T) {
	store := setupStore()

	added := defaultBook(store)

	got, err := store.GetBookByID(added.ID)
	if err != nil {
		t.Fatalf("Expected book to exist, got error: %v", err)
	}
	if got.Title != added.Title {
		t.Errorf("Expected title %s, got %s", added.Title, got.Title)
	}
}

func TestUpdateBook(t *testing.T) {
	store := setupStore()

	added := defaultBook(store)

	update := Book{
		Title:         "Updated Book",
		Author:        "Tester",
		PublishedYear: 2024,
	}

	updated, err := store.UpdateBook(added.ID, update)
	if err != nil {
		t.Fatalf("Failed to update book: %v", err)
	}

	if updated.Title != "Updated Book" {
		t.Errorf("Update failed, expected title %s, got %s", "Updated Book", updated.Title)
	}
}

func TestUpdateBookNotFoundId(t *testing.T) {
	store := setupStore()

	update := Book{
		Title:         "Updated Book",
		Author:        "Tester",
		PublishedYear: 2024,
	}

	updated, err := store.UpdateBook(1, update)

	if err == nil {
		t.Errorf("got: %v, want: book not found", updated)
	}
}

func TestDeleteBook(t *testing.T) {
	store := setupStore()

	added := defaultBook(store)

	err := store.DeleteBook(added.ID)
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}

	_, err = store.GetBookByID(added.ID)
	if err == nil {
		t.Error("Expected error after deleting book, got none")
	}
}

func TestDeleteBookNotFound(t *testing.T) {
	store := setupStore()

	err := store.DeleteBook(1)
	if err == nil {
		t.Error("got: success deleted, want: book not found")
	}
}

func BenchmarkAddBooks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		store := setupStore()
		for j := 0; j < 1000; j++ {
			store.AddBook(Book{
				Title:         "Book " + strconv.Itoa(j),
				Author:        "Author",
				PublishedYear: 2023,
			})
		}
	}
}

func BenchmarkGetBooks(b *testing.B) {
	store := setupStore()
	for i := 0; i < 1000; i++ {
		store.AddBook(Book{
			Title:         "Book " + strconv.Itoa(i),
			Author:        "Author",
			PublishedYear: 2023,
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := 1
		for pb.Next() {
			_, err := store.GetBookByID(id)
			if err != nil {
				b.Fatalf("Failed to get book ID %d: %v", id, err)
			}
			id = (id % 1000) + 1
		}
	})
}

func BenchmarkUpdateBooks(b *testing.B) {
	store := setupStore()
	for i := 0; i < 1000; i++ {
		store.AddBook(Book{
			Title:         "Book " + strconv.Itoa(i),
			Author:        "Author",
			PublishedYear: 2023,
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := 1
		for pb.Next() {
			_, err := store.UpdateBook(id, Book{
				Title:         "Updated Title",
				Author:        "Updated Author",
				PublishedYear: 2025,
			})
			if err != nil {
				b.Fatalf("Failed to update book ID %d: %v", id, err)
			}
			id = (id % 1000) + 1
		}
	})
}

func BenchmarkDeleteBooks(b *testing.B) {
	store := setupStore()
	ids := make([]int, 1000)

	for i := 0; i < 1000; i++ {
		book := store.AddBook(Book{
			Title:         "Book " + strconv.Itoa(i),
			Author:        "Author",
			PublishedYear: 2023,
		})
		ids[i] = book.ID
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		err := store.DeleteBook(id)
		if err != nil {
			continue
		}
	}
}
