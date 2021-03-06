package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bxcodec/library/domain"
	"github.com/sirupsen/logrus"
)

type postgresBookRepository struct {
	Conn *sql.DB
}

func NewPostgresBookRepository(Conn *sql.DB) domain.BookRepository {
	return &postgresBookRepository{Conn}
}

func (p *postgresBookRepository) fetch(c context.Context, query string) ([]domain.Book, error) {
	rows, err := p.Conn.QueryContext(c, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result := make([]domain.Book, 0)
	for rows.Next() {
		book := domain.Book{}
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Content,
			&book.Author.ID,
			&book.Author.Name,
			&book.Author.CreatedAt,
			&book.Author.UpdatedAt,
			&book.UpdatedAt,
			&book.CreatedAt,
		)

		if err != nil {
			logrus.Error(rows)
			return nil, err
		}
		result = append(result, book)
	}
	return result, nil
}

func (p *postgresBookRepository) Fetch(ctx context.Context, num int, offset int) ([]domain.Book, error) {
	query := fmt.Sprintf("SELECT b.id, b.title, b.content, a.id, a.name, a.created_at, a.updated_at, b.created_at, b.updated_at " +
		"FROM book as b " +
		"INNER JOIN author a on b.author_id = a.id " +
		"ORDER BY b.id")
	if num != 0 {
		query = fmt.Sprintf(query+" LIMIT %d", num)
	}
	if offset != 0 {
		query = fmt.Sprintf(query+" OFFSET %d", offset)
	}

	res, err := p.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (p *postgresBookRepository) Add(ctx context.Context, book *domain.Book) error {
	res, err := p.Conn.ExecContext(ctx, `INSERT INTO book (title, content, author_id, updated_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
		book.Title, book.Content, book.Author.ID, book.UpdatedAt, book.CreatedAt)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
		return err
	}

	return err
}

func (p *postgresBookRepository) Delete(ctx context.Context, id int) error {
	res, err := p.Conn.ExecContext(ctx, `DELETE FROM book WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
		return err
	}
	return err
}

func (p *postgresBookRepository) Update(ctx context.Context, book *domain.Book) error {
	res, err := p.Conn.ExecContext(ctx, `UPDATE  book SET title = $1, content = $2, author_id = $3, updated_at = $4, created_at = $5 WHERE id = $6`,
		book.Title, book.Content, book.Author.ID, book.UpdatedAt, book.CreatedAt, book.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
		return err
	}
	return err
}

func (p *postgresBookRepository) GetById(ctx context.Context, id int) (domain.Book, error) {
	query := fmt.Sprintf("SELECT b.id, b.title, b.content, a.id, a.name, a.created_at, a.updated_at, b.created_at, b.updated_at "+
		"FROM book as b "+
		"INNER JOIN author a on b.author_id = a.id "+
		"WHERE b.id = %d", id)
	list, err := p.fetch(ctx, query)
	if err != nil {
		return domain.Book{}, err
	}
	if len(list) == 0 {
		return domain.Book{}, domain.ErrNotFound
	}
	return list[0], nil
}
