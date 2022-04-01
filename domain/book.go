package domain

import (
	"context"
)

type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title" validate:"required"`
	Content   string `json:"content" validate:"required"`
	Author    Author `json:"author"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

// BookUseCase represent the book's use case contract
type BookUseCase interface {
	Fetch(ctx context.Context, num int) ([]Book, error)
	Add(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, book *Book) error
	GetById(ctx context.Context, id int) (Book, error)
}

// BookRepository represent the book's repository contract
type BookRepository interface {
	Fetch(ctx context.Context, num int) ([]Book, error)
	Add(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, book *Book) error
	GetById(ctx context.Context, id int) (Book, error)
}
