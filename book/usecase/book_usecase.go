package usecase

import (
	"context"
	"fmt"
	"github.com/bxcodec/library/domain"
	"github.com/bxcodec/library/mail"
	"github.com/bxcodec/library/message_broker"
	"github.com/bxcodec/library/message_broker/rabbit"
	"log"
	"time"
)

type bookUseCase struct {
	bookRepo       domain.BookRepository
	authorRepo     domain.AuthorRepository
	messageBroker  message_broker.MessageBroker
	contextTimeout time.Duration
	mailService    mail.MailService
}

func NewBookUseCase(b domain.BookRepository, ar domain.AuthorRepository, mb message_broker.MessageBroker, mail mail.MailService, timeout time.Duration) domain.BookUseCase {
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
		b.publishToMsgBroker(message_broker.FETCH, book.Content)
	}
	//b.sendToMail(message_broker.FETCH)
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

	b.publishToMsgBroker(message_broker.ADD, book.Content)
	//b.sendToMail(message_broker.ADD)
	return
}

func (b *bookUseCase) Delete(c context.Context, id int) (err error) {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	existingBook, err := b.GetById(ctxt, id)
	if err != nil {
		return err
	}
	if existingBook == (domain.Book{}) {
		return domain.ErrNotFound
	}
	err = b.bookRepo.Delete(ctxt, existingBook.ID)

	b.publishToMsgBroker(message_broker.DELETE, existingBook.Content)
	//b.sendToMail(message_broker.DELETE)
	return
}

func (b *bookUseCase) Update(c context.Context, book *domain.Book) (err error) {
	ctxt, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	err = b.bookRepo.Update(ctxt, book)

	b.publishToMsgBroker(message_broker.UPDATE, book.Content)
	//b.sendToMail(message_broker.UPDATE)
	return
}

func (b *bookUseCase) GetById(ctx context.Context, id int) (res domain.Book, err error) {
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

	b.publishToMsgBroker(message_broker.GetById, res.Content)
	go b.sendToMail(message_broker.GetById)
	return
}

func (b *bookUseCase) publishToMsgBroker(eventType message_broker.EventType, content string) {
	err := b.messageBroker.Send(eventType, content)
	if err != nil {
		log.Println(rabbit.FailedPublishing)
	}
}

func (b *bookUseCase) sendToMail(eventType message_broker.EventType) {

	emailChan, _ := b.messageBroker.Receive(eventType)
	go send(eventType, emailChan)
}

func send(eventType message_broker.EventType, emailChan chan string) error {

	//fmt.Println("DATA ", <-emailChan)

	for email := range emailChan {
		fmt.Println("DATA ", email)

		//body := fmt.Sprintf("%s", email)
		//event := message_broker.Event{Content: body}
		//event.Subject = string(eventType)
		//useCase := mail.NewMailUseCase()
		//useCase.SendEmail(event)
	}

	return nil
}
