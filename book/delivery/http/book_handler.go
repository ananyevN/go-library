package http

import (
	"github.com/bxcodec/library/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

const NAME = "num"
const ID = "id"

type ResponseError struct {
	Message string `json:"message"`
}

type BookHandler struct {
	BUseCase domain.BookUseCase
}

func NewBookHandler(e *echo.Echo, us domain.BookUseCase) {
	handler := &BookHandler{BUseCase: us}
	e.GET("/books", handler.FetchBook)
	e.GET("/books/:id", handler.GetById)
	e.POST("/books", handler.Add)
	e.POST("/book", handler.Update)
	e.DELETE("/books/:id", handler.Delete)
}

func (h BookHandler) FetchBook(c echo.Context) error {
	nums := c.QueryParam(NAME)
	num, _ := strconv.Atoi(nums)
	ctx := c.Request().Context()

	listBooks, err := h.BUseCase.Fetch(ctx, num)
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

	ctxt := c.Request().Context()
	book, err := h.BUseCase.GetById(ctxt, id)
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

	ctxt := c.Request().Context()
	err = h.BUseCase.Delete(ctxt, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h BookHandler) Add(c echo.Context) error {
	var book domain.Book
	err := c.Bind(&book)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	var ok bool
	if ok, err = isRequestValid(&book); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = h.BUseCase.Add(ctx, &book)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, book)
}

func (h BookHandler) Update(c echo.Context) error {
	var book domain.Book
	err := c.Bind(&book)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	var ok bool
	if ok, err = isRequestValid(&book); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = h.BUseCase.Update(ctx, &book)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, book)
}

func isRequestValid(b *domain.Book) (bool, error) {
	validate := validator.New()
	err := validate.Struct(b)
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
