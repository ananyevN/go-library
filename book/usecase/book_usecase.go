package usecase

import (
	"context"
	"log"
	"time"

	"github.com/bxcodec/library/domain"
	"github.com/bxcodec/library/mail"
	mb "github.com/bxcodec/library/message_broker"
)

type bookUseCase struct {
	bookRepo       domain.BookRepository
	authorRepo     domain.AuthorRepository
	messageBroker  mb.MessageBroker
	contextTimeout time.Duration
	mailService    mail.Sender
}

func NewBookUseCase(b domain.BookRepository, ar domain.AuthorRepository, mb mb.MessageBroker, mail mail.Sender, timeout time.Duration) domain.BookUseCase {
	return &bookUseCase{
		bookRepo:       b,
		authorRepo:     ar,
		messageBroker:  mb,
		mailService:    mail,
		contextTimeout: timeout,
	}
}

func (b *bookUseCase) Fetch(c context.Context, num int, offset int) ([]domain.Book, error) {
	if num == 0 {
		num = 10
	}
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	res, err := b.bookRepo.Fetch(ctxt, num, offset)
	if err != nil {
		return nil, err
	}

	for _, book := range res {
		b.publishToMsgBrokerAndSendEmail(mb.FETCH, book.Content)
	}
	return res, err
}

func (b *bookUseCase) Add(c context.Context, book *domain.Book) error {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	existingBook, err := b.getById(ctxt, book.ID)
	if existingBook != (domain.Book{}) {
		return domain.ErrConflict
	}
	if err = b.bookRepo.Add(ctxt, book); err != nil {
		return err
	}
	b.publishToMsgBrokerAndSendEmail(mb.ADD, book.Content)
	return nil
}

func (b *bookUseCase) Delete(c context.Context, id int) error {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	existingBook, err := b.getById(c, id)
	if err != nil {
		return err
	}
	if existingBook == (domain.Book{}) {
		return domain.ErrNotFound
	}
	if err = b.bookRepo.Delete(ctxt, existingBook.ID); err != nil {
		return err
	}
	b.publishToMsgBrokerAndSendEmail(mb.DELETE, existingBook.Content)
	return nil
}

func (b *bookUseCase) Update(c context.Context, book *domain.Book) error {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	if err := b.bookRepo.Update(ctxt, book); err != nil {
		return err
	}

	b.publishToMsgBrokerAndSendEmail(mb.UPDATE, book.Content)
	return nil
}

func (b *bookUseCase) GetById(ctx context.Context, id int) (domain.Book, error) {
	res, err := b.getById(ctx, id)

	if err != nil {
		return domain.Book{}, err
	}
	b.publishToMsgBrokerAndSendEmail(mb.GetById, res.Content)
	return res, err
}

func (b *bookUseCase) getById(ctx context.Context, id int) (domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	res, err := b.bookRepo.GetById(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}
	resAuthor, err := b.authorRepo.GetById(ctx, res.Author.ID)
	if err != nil {
		return domain.Book{}, err
	}
	res.Author = resAuthor

	return res, err
}

func (b *bookUseCase) publishToMsgBrokerAndSendEmail(eventType mb.EventType, content string) {
	event := mb.Event{Content: content, Subject: string(eventType)}
	err := b.messageBroker.Send(event)
	if err != nil {
		log.Println(err.Error())
	}

	emailChan := make(chan []byte, 1)

	go func() {
		err := b.messageBroker.Receive(emailChan)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	go b.sendToMail(emailChan)
}

func (b *bookUseCase) sendToMail(emailChan chan []byte) {
	for e := range emailChan {
		event := mb.Event{}
		err := b.mailService.SendEmail(*event.Unmarshal(e))
		if err != nil {
			log.Println(err.Error())
		}
	}
}
