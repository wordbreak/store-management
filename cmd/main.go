package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"store-management/internal/datasource"
	"store-management/internal/middleware"
	"store-management/internal/repository"
	"store-management/internal/response"
	"store-management/internal/router"
	"store-management/internal/service"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, response.New(http.StatusOK, response.MessageOK, nil))
	})

	sqlWriter := datasource.NewMySQL(&datasource.MySQLConfig{
		User:   os.Getenv("DB_WRITER_USER"),
		Passwd: os.Getenv("DB_WRITER_PASS"),
		Host:   os.Getenv("DB_WRITER_HOST"),
		Port:   os.Getenv("DB_WRITER_PORT"),
		DBName: os.Getenv("DB_WRITER_DBNAME"),
	})
	defer sqlWriter.Close()
	sqlReader := datasource.NewMySQL(&datasource.MySQLConfig{
		User:   os.Getenv("DB_READER_USER"),
		Passwd: os.Getenv("DB_READER_PASS"),
		Host:   os.Getenv("DB_READER_HOST"),
		Port:   os.Getenv("DB_READER_PORT"),
		DBName: os.Getenv("DB_READER_DBNAME"),
	})
	defer sqlReader.Close()

	repository.Init(sqlWriter, sqlReader)
	service.Init(repository.Get())
	r.Use(middleware.JwtMiddleware())
	router.Init(r, service.Get())

	go func() {
		err = r.Run()
		if err != nil {
			panic(err)
		}
	}()

	<-shutdownChan
}
