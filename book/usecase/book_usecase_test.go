package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/bxcodec/library/domain"
	"github.com/bxcodec/library/message_broker"
	"github.com/bxcodec/library/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	mockBook := domain.Book{
		ID:      14,
		Title:   "Hello",
		Content: "World",
	}

	mockListOfBook := make([]domain.Book, 0)
	mockListOfBook = append(mockListOfBook, mockBook)

	mockEvent := message_broker.Event{
		Content: mockBook.Content,
		Subject: "fetch.sql",
	}

	mockRabbitMq := new(mocks.MessageBroker)
	mockBookRepo := new(mocks.BookRepository)
	mockMailUseCase := new(mocks.MailService)

	mockBookRepo.On("Fetch", mock.Anything, mock.AnythingOfType("int")).
		Return(mockListOfBook, nil).Once()
	mockRabbitMq.On("Send", mockEvent).Return(nil).Once()

	mockAuthorRepo := new(mocks.AuthorRepository)
	usecase := NewBookUseCase(mockBookRepo, mockAuthorRepo, mockRabbitMq, mockMailUseCase, time.Second*2)
	fetch, err := usecase.Fetch(context.TODO(), 1, 0)

	assert.NoError(t, err)
	assert.Equal(t, mockListOfBook, fetch)
	mockBookRepo.AssertExpectations(t)
	mockAuthorRepo.AssertExpectations(t)
	mockRabbitMq.AssertExpectations(t)
	mockMailUseCase.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockBook := domain.Book{
		ID:      14,
		Title:   "Hello",
		Content: "World",
	}

	mockEvent := message_broker.Event{
		Content: mockBook.Content,
		Subject: "update.sql",
	}

	mockBookRepo := new(mocks.BookRepository)
	mockMailUseCase := new(mocks.MailService)
	mockRabbitMq := new(mocks.MessageBroker)
	mockAuthorRepo := new(mocks.AuthorRepository)

	mockBookRepo.On("Update", mock.Anything, &mockBook).Once().Return(nil)
	mockRabbitMq.On("Send", mockEvent).Return(nil).Once()

	usecase := NewBookUseCase(mockBookRepo, mockAuthorRepo, mockRabbitMq, mockMailUseCase, time.Second*2)
	err := usecase.Update(context.TODO(), &mockBook)

	assert.NoError(t, err)
	mockBookRepo.AssertExpectations(t)
	mockRabbitMq.AssertExpectations(t)
	mockMailUseCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockBook := domain.Book{
		ID:      14,
		Title:   "Hello",
		Content: "World",
	}
	authorMock := domain.Author{
		ID:   0,
		Name: "TestAuthor",
	}

	mockEvent := message_broker.Event{
		Content: mockBook.Content,
		Subject: "delete.sql",
	}

	mockBookRepo := new(mocks.BookRepository)
	mockAuthorRepo := new(mocks.AuthorRepository)
	mockRabbitMq := new(mocks.MessageBroker)
	mockMailUseCase := new(mocks.MailService)

	mockBookRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(mockBook, nil).Once()
	mockAuthorRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(authorMock, nil).Once()
	mockBookRepo.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()
	mockRabbitMq.On("Send", mockEvent).Return(nil).Once()

	usecase := NewBookUseCase(mockBookRepo, mockAuthorRepo, mockRabbitMq, mockMailUseCase, time.Second*2)
	err := usecase.Delete(context.TODO(), mockBook.ID)

	assert.NoError(t, err)
	mockBookRepo.AssertExpectations(t)
	mockAuthorRepo.AssertExpectations(t)
	mockRabbitMq.AssertExpectations(t)
	mockMailUseCase.AssertExpectations(t)
}
