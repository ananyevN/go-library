package http

import (
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"

	"github.com/bxcodec/library/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (
	NUM    = "num"
	OFFSET = "offset"
	ID     = "id"
)

type ResponseError struct {
	Message string `json:"message"`
}

type BookHandler struct {
	BUseCase  domain.BookUseCase
	validator *validator.Validate
}

func NewBookHandler(e *echo.Echo, us domain.BookUseCase) {
	handler := &BookHandler{BUseCase: us, validator: validator.New()}
	e.GET("/books", handler.FetchBook)
	e.GET("/books/:id", handler.GetById)
	e.POST("/books", handler.Add)
	e.PUT("/book", handler.Update)
	e.DELETE("/books/:id", handler.Delete)
}

func (h BookHandler) FetchBook(c echo.Context) error {
	nums := c.QueryParam(NUM)
	num, _ := strconv.Atoi(nums)
	offsets := c.QueryParam(OFFSET)
	offset, _ := strconv.Atoi(offsets)

	listBooks, err := h.BUseCase.Fetch(c.Request().Context(), num, offset)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}
	return c.JSON(http.StatusOK, listBooks)
}

func (h BookHandler) GetById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param(ID))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	book, err := h.BUseCase.GetById(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}
	return c.JSON(http.StatusOK, book)
}

func (h BookHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param(ID))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	err = h.BUseCase.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h BookHandler) Add(c echo.Context) error {
	var book domain.Book
	var err error
	if err = c.Bind(&book); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if ok, err := h.isRequestValid(&book); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err = h.BUseCase.Add(c.Request().Context(), &book); err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, book)
}

func (h BookHandler) Update(c echo.Context) error {
	var book domain.Book
	var err error
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if ok, err := h.isRequestValid(&book); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err = h.BUseCase.Update(c.Request().Context(), &book); err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, book)
}

func (h BookHandler) isRequestValid(b *domain.Book) (bool, error) {
	err := h.validator.Struct(b)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
