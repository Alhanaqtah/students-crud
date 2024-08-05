package storage

import (
	"context"
	"fmt"

	"students-crud/internal/config"
	"students-crud/internal/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

// Create создает нового студента
func (s *Storage) Create(ctx context.Context, student *models.Student) (int, error) {
	const op = "storage.postgres.Create"

	var id int
	err := s.pool.QueryRow(ctx, "INSERT INTO students (name, email) VALUES ($1, $2) RETURNING id", student.Name, student.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Read читает студента по ID
func (s *Storage) Read(ctx context.Context, id int) (*models.Student, error) {
	const op = "storage.postgres.Read"

	student := &models.Student{}
	err := s.pool.QueryRow(ctx, "SELECT id, name, email FROM students WHERE id=$1", id).Scan(&student.ID, &student.Name, &student.Email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return student, nil
}

// Update обновляет информацию о студенте
func (s *Storage) Update(ctx context.Context, student *models.Student) error {
	const op = "storage.postgres.Update"

	_, err := s.pool.Exec(ctx, "UPDATE students SET name=$1, email=$2 WHERE id=$3", student.Name, student.Email, student.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Delete удаляет студента по ID
func (s *Storage) Delete(ctx context.Context, id int) error {
	const op = "storage.postgres.Delete"

	_, err := s.pool.Exec(ctx, "DELETE FROM students WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
