package usecase

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"time"
)

type bookUseCase struct {
	bookRepo       domain.BookRepository
	authorRepo     domain.AuthorRepository
	contextTimeout time.Duration
}

func NewBookUseCase(b domain.BookRepository, ar domain.AuthorRepository, timeout time.Duration) domain.BookUseCase {
	return &bookUseCase{
		bookRepo:       b,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

func (b *bookUseCase) Fetch(c context.Context, num int) (res []domain.Book, err error) {
	if num == 0 {
		num = 10
	}
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	res, err = b.bookRepo.Fetch(ctxt, num)
	if err != nil {
		return nil, err
	}
	return
}

func (b *bookUseCase) Add(c context.Context, book *domain.Book) (err error) {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	existingBook, err := b.GetById(ctxt, book.ID)
	if existingBook != (domain.Book{}) {
		return domain.ErrConflict
	}
	err = b.bookRepo.Add(ctxt, book)
	return
}

func (b *bookUseCase) Delete(c context.Context, id int) error {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	existingBook, err := b.GetById(ctxt, id)
	if err != nil {
		return err
	}
	if existingBook == (domain.Book{}) {
		return domain.ErrNotFound
	}
	return b.bookRepo.Delete(ctxt, existingBook.ID)
}

func (b *bookUseCase) Update(c context.Context, book *domain.Book) error {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	return b.bookRepo.Update(ctxt, book)
}

func (b *bookUseCase) GetById(ctx context.Context, id int) (res domain.Book, err error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	res, err = b.bookRepo.GetById(ctx, id)
	if err != nil {
		return
	}
	resAuthor, err := b.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return domain.Book{}, err
	}
	res.Author = resAuthor
	return
}
