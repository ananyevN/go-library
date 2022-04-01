package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bxcodec/go-clean-arch/domain"
)

type postgresAuthorRepo struct {
	DB *sql.DB
}

func NewPostgresAuthorRepository(db *sql.DB) domain.AuthorRepository {
	return &postgresAuthorRepo{
		DB: db,
	}
}

func (p postgresAuthorRepo) GetByID(ctx context.Context, id int) (domain.Author, error) {
	query := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM author WHERE id=%d", id)
	return p.getOne(ctx, query)
}

func (p postgresAuthorRepo) getOne(ctx context.Context, query string) (res domain.Author, err error) {
	stmt, err := p.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Author{}, err
	}

	row := stmt.QueryRow()
	res = domain.Author{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}
