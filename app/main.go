package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/bxcodec/library/author/repository/postgres"
	"github.com/bxcodec/library/book/delivery/http"
	_book_repository "github.com/bxcodec/library/book/repository/postgres"
	"github.com/bxcodec/library/book/usecase"
	"github.com/bxcodec/library/mail"
	"github.com/bxcodec/library/message_broker/rabbit"
)

func init() {
	viper.SetConfigFile(`config.json`)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetInt(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	url := fmt.Sprintf("%s://%s:%s@%s:%d/postgres?sslmode=disable", dbName, dbUser, dbPass, dbHost, dbPort)

	dbConn, err := sql.Open(`postgres`, url)

	if err != nil {
		log.Fatal(err)
	}
	if err = dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()

	authorRepo := postgres.NewPostgresAuthorRepository(dbConn)
	bookRepo := _book_repository.NewPostgresBookRepository(dbConn)
	rabbitMqService := rabbit.NewRabbitMqService("crud_exchange")
	mailUseCase := mail.NewSender()

	timeoutContext := time.Second
	bookUseCase := usecase.NewBookUseCase(bookRepo, authorRepo, rabbitMqService, mailUseCase, timeoutContext)

	http.NewBookHandler(e, bookUseCase)

	log.Fatal(e.Start(":9000"))
}
