package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"students-crud/internal/config"
	"students-crud/internal/models"
	"students-crud/internal/storage"

	"github.com/gin-gonic/gin"
)

type Storage interface {
	Create(ctx context.Context, student *models.Student) (int, error)
	Read(ctx context.Context, id int) (*models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id int) error
}

func main() {
	cfg := config.MustLoad()

	storage, err := storage.New(&cfg.Storage)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	r := gin.Default()

	// Создание нового студента
	r.POST("/students", func(ctx *gin.Context) {
		var s models.Student
		jsonData, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println("failed to read request body")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		err = json.Unmarshal(jsonData, &s)
		if err != nil {
			log.Println("failed to unmarshal data")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to unmarshal data"})
			return
		}

		id, err := storage.Create(ctx.Request.Context(), &s)
		if err != nil {
			log.Println("failed to create student:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	})

	// Получение студента по ID
	r.GET("/students/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr) // Преобразование строки в int
		if err != nil {
			log.Println("invalid id:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		student, err := storage.Read(ctx.Request.Context(), id)
		if err != nil {
			log.Println("failed to read student:", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}

		ctx.JSON(http.StatusOK, student)
	})

	// Обновление студента
	r.PUT("/students/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("invalid id:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var s models.Student
		jsonData, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println("failed to read request body")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			return
		}

		err = json.Unmarshal(jsonData, &s)
		if err != nil {
			log.Println("failed to unmarshal data")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to unmarshal data"})
			return
		}

		s.ID = id // Устанавливаем ID студента для обновления

		err = storage.Update(ctx.Request.Context(), &s)
		if err != nil {
			log.Println("failed to update student:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "student updated successfully"})
	})

	// Удаление студента по ID
	r.DELETE("/students/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr) // Преобразование строки в int
		if err != nil {
			log.Println("invalid id:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		err = storage.Delete(ctx.Request.Context(), id)
		if err != nil {
			log.Println("failed to delete student:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete student"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "student deleted successfully"})
	})

	r.Run(os.Getenv("ADDRESS"))
}
