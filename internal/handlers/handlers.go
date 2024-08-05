package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

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

type Handlers struct {
	storage *storage.Storage
}

// NewHandlers создает новый экземпляр Handlers
func NewHandlers(storage *storage.Storage) *Handlers {
	return &Handlers{storage: storage}
}

// Создание нового студента
func (h *Handlers) CreateStudent(ctx *gin.Context) {
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

	id, err := h.storage.Create(ctx.Request.Context(), &s)
	if err != nil {
		log.Println("failed to create student:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// Получение студента по ID
func (h *Handlers) ReadStudent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("invalid id:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	student, err := h.storage.Read(ctx.Request.Context(), id)
	if err != nil {
		log.Println("failed to read student:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

// Обновление студента
func (h *Handlers) UpdateStudent(ctx *gin.Context) {
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

	err = h.storage.Update(ctx.Request.Context(), &s)
	if err != nil {
		log.Println("failed to update student:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "student updated successfully"})
}

// Удаление студента по ID
func (h *Handlers) DeleteStudent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("invalid id:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.storage.Delete(ctx.Request.Context(), id)
	if err != nil {
		log.Println("failed to delete student:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete student"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "student deleted successfully"})
}
