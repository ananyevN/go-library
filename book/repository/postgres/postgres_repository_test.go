package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/library/domain"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error %s occured while testing", err)
	}
	books := make([]domain.Book, 2)
	books[0] = domain.Book{
		ID: 1, Title: "title 1", Content: "content 1",
		Author: domain.Author{ID: 1}, UpdatedAt: time.Now().String(), CreatedAt: time.Now().String(),
	}
	books[1] = domain.Book{
		ID: 2, Title: "title 2", Content: "content 2",
		Author: domain.Author{ID: 1}, UpdatedAt: time.Now().String(), CreatedAt: time.Now().String(),
	}
	rows := sqlmock.NewRows([]string{"b.id", "b.title", "b.content", "a.id", "a.name", "a.created_at", "a.updated_at", "b.created_at", "b.updated_at"}).
		AddRow(books[0].ID, books[0].Title, books[0].Content, books[0].Author.ID, books[0].Author.Name, books[0].Author.CreatedAt, books[0].Author.UpdatedAt, books[0].UpdatedAt, books[0].CreatedAt).
		AddRow(books[1].ID, books[1].Title, books[1].Content, books[1].Author.ID, books[1].Author.Name, books[1].Author.CreatedAt, books[1].Author.UpdatedAt, books[1].UpdatedAt, books[1].CreatedAt)

	var num = 2
	var offset = 1
	query := fmt.Sprintf("SELECT b.id, b.title, b.content, a.id, a.name, a.created_at, a.updated_at, b.created_at, b.updated_at "+
		"FROM book as b "+
		"INNER JOIN author a on b.author_id = a.id "+
		"ORDER BY b.id LIMIT %d OFFSET %d", num, offset)
	mock.ExpectQuery(query).WillReturnRows(rows)

	bookRepository := NewPostgresBookRepository(db)
	result, err := bookRepository.Fetch(context.TODO(), num, offset)
	assert.Equal(t, num, len(result))
}

func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error %s occured while testing", err)
	}

	author := domain.Author{
		ID:        1,
		Name:      "Author X",
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}
	book := &domain.Book{
		Title:     "KISS",
		Content:   "That Girl",
		Author:    author,
		UpdatedAt: time.Now().String(),
		CreatedAt: time.Now().String(),
	}

	mock.ExpectExec("INSERT INTO book (title, content, author_id, updated_at, created_at)*").
		WithArgs(book.Title, book.Content, author.ID, book.UpdatedAt, book.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	bookRepository := NewPostgresBookRepository(db)
	err = bookRepository.Add(context.TODO(), book)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error %s occured while testing", err)
	}
	var num = 1
	mock.ExpectExec("DELETE FROM book WHERE id = ").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(int64(num), 1))
	bookRepository := NewPostgresBookRepository(db)
	err = bookRepository.Delete(context.TODO(), int(num))
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error %s occured while testing", err)
	}
	var num = 1
	author := domain.Author{
		ID:        1,
		Name:      "Author X",
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}
	book := &domain.Book{
		Title:     "KISS",
		Content:   "That Girl",
		Author:    author,
		UpdatedAt: time.Now().String(),
		CreatedAt: time.Now().String(),
	}

	mock.ExpectExec(`UPDATE book SET *`).
		WithArgs(book.Title, book.Content, book.Author.ID, book.UpdatedAt, book.CreatedAt, book.ID).
		WillReturnResult(sqlmock.NewResult(int64(num), 1))
	bookRepository := NewPostgresBookRepository(db)
	err = bookRepository.Update(context.TODO(), book)
	assert.NoError(t, err)
}

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error %s occured while testing", err)
	}
	book := domain.Book{
		ID: 1, Title: "title 1", Content: "content 1",
		Author: domain.Author{ID: 1}, UpdatedAt: time.Now().String(), CreatedAt: time.Now().String(),
	}
	rows := sqlmock.NewRows([]string{"b.id", "b.title", "b.content", "a.id", "a.name", "a.created_at", "a.updated_at", "b.created_at", "b.updated_at"}).
		AddRow(book.ID, book.Title, book.Content, book.Author.ID, book.Author.Name, book.Author.CreatedAt, book.Author.UpdatedAt, book.UpdatedAt, book.CreatedAt)

	var id = 1
	query := fmt.Sprintf("SELECT b.id, b.title, b.content, a.id, a.name, a.created_at, a.updated_at, b.created_at, b.updated_at "+
		"FROM book as b "+
		"INNER JOIN author a on b.author_id = a.id "+
		"WHERE b.id = %d", id)
	mock.ExpectQuery(query).WillReturnRows(rows)

	bookRepository := NewPostgresBookRepository(db)
	bookActual, err := bookRepository.GetById(context.TODO(), id)
	assert.Equal(t, book, bookActual)
}
