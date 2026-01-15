package storage

import (
	"database/sql"
	"log/slog"
	"task-api/models"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil

}

func (s *PostgresStore) CreateTask(description string) (*models.Task, error) {
	query := `
		INSERT into tasks (description)
		VALUES ($1)
		RETURNING id, description, completed, created_at, completed_at
	`

	var task models.Task
	var completedAt sql.NullTime

	err := s.db.QueryRow(query, description).Scan(
		&task.ID,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&completedAt,
	)

	// Always check validity even on new tasks (though it will be NULL)
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, err
}

func (s *PostgresStore) GetAllTasks() ([]*models.Task, error) {
	query := `SELECT id, description, completed, created_at, completed_at from tasks`

	row, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	var tasks []*models.Task

	for row.Next() {
		var task models.Task
		var completedAt sql.NullTime

		err := row.Scan(
			&task.ID,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&completedAt,
		)

		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, &task)
	}
	return tasks, row.Err()
}

func (s *PostgresStore) CompletedTaskById(id int) (*models.Task, error) {
	query := `
        UPDATE tasks 
        SET completed = true, completed_at = $1 
        WHERE id = $2
        RETURNING id, description, completed, created_at, completed_at`

	var task models.Task
	var completedAt sql.NullTime // Use NullTime for safety
	now := time.Now()

	err := s.db.QueryRow(query, now, id).Scan(
		&task.ID,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&completedAt,
	)

	if err == sql.ErrNoRows {
		slog.Error("task not found", "id", id)
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil
}

func (s *PostgresStore) GetTaskById(id int) (*models.Task, error) {
	query := `SELECT id, description, completed, created_at, completed_at from tasks where id = $1`

	var task models.Task
	var completedAt sql.NullTime

	err := s.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&completedAt)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, err

}

func (s *PostgresStore) DeleteTaskById(id int) error {
	query := `DELETE from tasks where id = $1`

	res, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (s *PostgresStore) SearchTasks(query string) ([]*models.Task, error) {
	sqlQuery := `SELECT id, description, completed, created_at, completed_at 
                 FROM tasks 
                 WHERE LOWER(description) LIKE LOWER('%' || $1 || '%')`

	// Use Query instead of Exec
	rows, err := s.db.Query(sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *PostgresStore) TruncateTasks() error {

	_, err := s.db.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	return err
}
