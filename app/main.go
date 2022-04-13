package main

import (
	"database/sql"
	"github.com/bxcodec/library/author/repository/postgres"
	"github.com/bxcodec/library/book/delivery/http"
	_book_repository "github.com/bxcodec/library/book/repository/postgres"
	"github.com/bxcodec/library/book/usecase"
	"github.com/bxcodec/library/mail"
	"github.com/bxcodec/library/message_broker/rabbit"
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

	e := echo.New()

	authorRepo := postgres.NewPostgresAuthorRepository(dbConn)
	bookRepo := _book_repository.NewPostgresBookRepository(dbConn)
	rabbitMqService := rabbit.NewRabbitMqService("crud_exchange")
	mailUseCase := mail.NewMailUseCase()

	timeoutContext := time.Second
	bookUseCase := usecase.NewBookUseCase(bookRepo, authorRepo, rabbitMqService, mailUseCase, timeoutContext)

	http.NewBookHandler(e, bookUseCase)

	log.Fatal(e.Start(":9000"))
}
