package domain

import (
	"context"
)

type Author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

//go:generate mockgen -source=$GOFILE -destination=../mocks/mock_repository/mock_$GOFILE
// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	GetById(ctx context.Context, id int) (Author, error)
}
