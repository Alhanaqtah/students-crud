package main

import (
	"log"
	"os"

	"students-crud/internal/config"
	"students-crud/internal/handlers"
	"students-crud/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	storage, err := storage.New(&cfg.Storage)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	handlers := handlers.NewHandlers(storage)

	r := gin.Default()

	r.POST("/students", handlers.CreateStudent)
	r.GET("/students/:id", handlers.ReadStudent)
	r.PUT("/students/:id", handlers.UpdateStudent)
	r.POST("/students/:id", handlers.DeleteStudent)

	r.Run(os.Getenv("ADDRESS"))
}
