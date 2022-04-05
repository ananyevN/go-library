package http

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker"
	"github.com/bxcodec/library/domain"
	mocks "github.com/bxcodec/library/mocks/mock_repository"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	var mockBook domain.Book
	err := faker.FakeData(&mockBook)
	assert.NoError(t, err)

	mockUCase := new(mocks.BookUseCase)
	mockListBooks := make([]domain.Book, 0)
	mockListBooks = append(mockListBooks, mockBook)
	num := 1
	mockUCase.On("Fetch", mock.Anything, num).
		Return(mockListBooks, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/book?num=1", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := BookHandler{BUseCase: mockUCase}
	err = handler.FetchBook(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	var mockBook domain.Book
	err := faker.FakeData(&mockBook)
	assert.NoError(t, err)
	j, err := json.Marshal(mockBook)

	mockUCase := new(mocks.BookUseCase)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.Book")).
		Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/book", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := BookHandler{BUseCase: mockUCase}
	err = handler.Update(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestAdd(t *testing.T) {
	var mockBook domain.Book
	err := faker.FakeData(&mockBook)
	assert.NoError(t, err)
	j, err := json.Marshal(mockBook)

	mockUCase := new(mocks.BookUseCase)
	mockUCase.On("Add", mock.Anything, mock.AnythingOfType("*domain.Book")).
		Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/books", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := BookHandler{BUseCase: mockUCase}
	err = handler.Add(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	num := 1
	mockUCase := new(mocks.BookUseCase)
	mockUCase.On("Delete", mock.Anything, num).
		Return(nil)

	e := echo.New()
	id := fmt.Sprintf("%d", num)
	req, err := http.NewRequest(echo.DELETE, "/books"+id, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("books/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := BookHandler{BUseCase: mockUCase}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetById(t *testing.T) {
	num := 1
	var mockBook domain.Book
	err := faker.FakeData(&mockBook)
	assert.NoError(t, err)

	mockUCase := new(mocks.BookUseCase)
	mockUCase.On("GetById", mock.Anything, num).Return(mockBook, nil)

	e := echo.New()
	id := fmt.Sprintf("%d", num)
	req, err := http.NewRequest(echo.GET, "/books"+id, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("books/:id")
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := BookHandler{BUseCase: mockUCase}
	err = handler.GetById(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}
