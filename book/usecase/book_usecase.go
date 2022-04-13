package usecase

import (
	"context"
	"github.com/bxcodec/library/domain"
	"github.com/bxcodec/library/mail"
	mb "github.com/bxcodec/library/message_broker"
	"github.com/bxcodec/library/message_broker/rabbit"
	"log"
	"time"
)

type bookUseCase struct {
	bookRepo       domain.BookRepository
	authorRepo     domain.AuthorRepository
	messageBroker  mb.MessageBroker
	contextTimeout time.Duration
	mailService    mail.MailService
}

func NewBookUseCase(b domain.BookRepository, ar domain.AuthorRepository, mb mb.MessageBroker, mail mail.MailService, timeout time.Duration) domain.BookUseCase {
	return &bookUseCase{
		bookRepo:       b,
		authorRepo:     ar,
		messageBroker:  mb,
		mailService:    mail,
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

	for _, book := range res {
		b.publishToMsgBrokerAndSendEmail(mb.FETCH, book.Content)
	}
	return
}

func (b *bookUseCase) Add(c context.Context, book *domain.Book) (err error) {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	existingBook, err := b.getById(ctxt, book.ID)
	if existingBook != (domain.Book{}) {
		return domain.ErrConflict
	}
	err = b.bookRepo.Add(ctxt, book)

	b.publishToMsgBrokerAndSendEmail(mb.ADD, book.Content)
	return
}

func (b *bookUseCase) Delete(c context.Context, id int) (err error) {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	existingBook, err := b.getById(c, id)
	if err != nil {
		return err
	}
	if existingBook == (domain.Book{}) {
		return domain.ErrNotFound
	}
	err = b.bookRepo.Delete(ctxt, existingBook.ID)

	b.publishToMsgBrokerAndSendEmail(mb.DELETE, existingBook.Content)
	return
}

func (b *bookUseCase) Update(c context.Context, book *domain.Book) (err error) {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	err = b.bookRepo.Update(ctxt, book)

	b.publishToMsgBrokerAndSendEmail(mb.UPDATE, book.Content)
	return
}

func (b *bookUseCase) GetById(ctx context.Context, id int) (res domain.Book, err error) {
	res, err = b.getById(ctx, id)

	b.publishToMsgBrokerAndSendEmail(mb.GetById, res.Content)
	return
}

func (b *bookUseCase) getById(ctx context.Context, id int) (res domain.Book, err error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	res, err = b.bookRepo.GetById(ctx, id)
	if err != nil {
		return
	}
	resAuthor, err := b.authorRepo.GetById(ctx, res.Author.ID)
	if err != nil {
		return domain.Book{}, err
	}
	res.Author = resAuthor

	return
}

func (b *bookUseCase) publishToMsgBrokerAndSendEmail(eventType mb.EventType, content string) {
	event := mb.Event{Content: content, Subject: string(eventType)}
	err := b.messageBroker.Send(event)
	if err != nil {
		log.Println(rabbit.FailedPublishing)
	}

	emailChan := make(chan []byte, 1)

	go func() {
		err := b.messageBroker.Receive(emailChan)
		if err != nil {
			log.Println(rabbit.FailedReceiving)
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
