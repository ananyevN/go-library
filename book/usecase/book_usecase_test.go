package usecase

import (
	"context"
	"github.com/bxcodec/library/domain"
	"github.com/bxcodec/library/message_brocker/rabbit"
	mock_domain "github.com/bxcodec/library/mocks/mock_repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	mockBookRepo := new(mock_domain.BookRepository)
	mockBook := domain.Book{
		ID:      14,
		Title:   "Hello",
		Content: "World",
	}
	mockListOfBook := make([]domain.Book, 0)
	mockListOfBook = append(mockListOfBook, mockBook)

	mockBookRepo.On("Fetch", mock.Anything, mock.AnythingOfType("int")).
		Return(mockListOfBook, nil).Once()

	mockAuthorRepo := new(mock_domain.AuthorRepository)
	rabbitMqService := rabbit.NewRabbitMqService("book")
	usecase := NewBookUseCase(mockBookRepo, mockAuthorRepo, rabbitMqService, time.Second*2)
	fetch, err := usecase.Fetch(context.TODO(), 1)

	assert.NoError(t, err)
	assert.Equal(t, mockListOfBook, fetch)
	mockBookRepo.AssertExpectations(t)
	mockAuthorRepo.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockBookRepo := new(mock_domain.BookRepository)
	mockBook := domain.Book{
		ID:      14,
		Title:   "Hello",
		Content: "World",
	}

	mockBookRepo.On("Update", mock.Anything, &mockBook).Once().Return(nil)

	mockAuthorRepo := new(mock_domain.AuthorRepository)
	rabbitMqService := rabbit.NewRabbitMqService("book")

	usecase := NewBookUseCase(mockBookRepo, mockAuthorRepo, rabbitMqService, time.Second*2)
	err := usecase.Update(context.TODO(), &mockBook)

	assert.NoError(t, err)
	mockBookRepo.AssertExpectations(t)
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

	mockBookRepo := new(mock_domain.BookRepository)
	mockAuthorRepo := new(mock_domain.AuthorRepository)

	mockBookRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(mockBook, nil).Once()
	mockAuthorRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(authorMock, nil).Once()
	mockBookRepo.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()

	rabbitMqService := rabbit.NewRabbitMqService("book")

	usecase := NewBookUseCase(mockBookRepo, mockAuthorRepo, rabbitMqService, time.Second*2)
	err := usecase.Delete(context.TODO(), mockBook.ID)

	assert.NoError(t, err)
	mockBookRepo.AssertExpectations(t)
	mockAuthorRepo.AssertExpectations(t)
}
