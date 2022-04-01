package main

import (
	"database/sql"
	"fmt"
	"github.com/bxcodec/go-clean-arch/author/repository/postgres"
	"github.com/bxcodec/go-clean-arch/book/delivery/http"
	_book_repository "github.com/bxcodec/go-clean-arch/book/repository/postgres"
	"github.com/bxcodec/go-clean-arch/book/usecase"
	"github.com/labstack/echo"
	_ "github.com/lib/pq" // <------------ here
	"log"
	"time"
)

func main() {
	dbConn, err := sql.Open(`postgres`,
		"postgres://postgres:password@host.docker.internal:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("INIT is OK")
	e := echo.New()

	authorRepo := postgres.NewPostgresAuthorRepository(dbConn)
	bookRepo := _book_repository.NewPostgresBookRepository(dbConn)

	timeoutContext := time.Duration(2) * time.Second
	bookUseCase := usecase.NewBookUseCase(bookRepo, authorRepo, timeoutContext)

	http.NewBookHandler(e, bookUseCase)

	log.Fatal(e.Start(":9000"))
}
