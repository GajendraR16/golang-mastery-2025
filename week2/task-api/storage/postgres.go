package storage

import (
	"context"
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

	//Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil

}

func (s *PostgresStore) CreateTask(ctx context.Context, description string) (*models.Task, error) {
	logQueryDuration("CreateTask", time.Now())
	query := `
		INSERT into tasks (description)
		VALUES ($1)
		RETURNING id, description, completed, created_at, completed_at
	`

	var task models.Task
	var completedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, query, description).Scan(
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

func (s *PostgresStore) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	logQueryDuration("GetAllTasks", time.Now())
	query := `SELECT id, description, completed, created_at, completed_at from tasks`

	row, err := s.db.QueryContext(ctx, query)
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

func (s *PostgresStore) CompletedTaskById(ctx context.Context, id int) (*models.Task, error) {
	logQueryDuration("CompletedTaskById", time.Now())
	query := `
        UPDATE tasks 
        SET completed = true, completed_at = $1 
        WHERE id = $2
        RETURNING id, description, completed, created_at, completed_at`

	var task models.Task
	var completedAt sql.NullTime // Use NullTime for safety
	now := time.Now()

	err := s.db.QueryRowContext(ctx, query, now, id).Scan(
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

func (s *PostgresStore) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	logQueryDuration("GetTaskById", time.Now())
	query := `SELECT id, description, completed, created_at, completed_at from tasks where id = $1`

	var task models.Task
	var completedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, query, id).Scan(
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

func (s *PostgresStore) DeleteTaskById(ctx context.Context, id int) error {
	logQueryDuration("DeleteTaskById", time.Now())
	query := `DELETE from tasks where id = $1`

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (s *PostgresStore) SearchTasks(ctx context.Context, query string) ([]*models.Task, error) {
	logQueryDuration("SearchTasks", time.Now())
	sqlQuery := `SELECT id, description, completed, created_at, completed_at 
                 FROM tasks 
                 WHERE LOWER(description) LIKE LOWER('%' || $1 || '%')`

	// Use Query instead of Exec
	rows, err := s.db.QueryContext(ctx, sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		//Set completedAt for every row sql.NullTime
		var completedAt sql.NullTime
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt, &completedAt)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (s *PostgresStore) TruncateTasks(ctx context.Context) error {

	_, err := s.db.ExecContext(ctx, "TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	return err
}

func logQueryDuration(name string, start time.Time) {
	defer func() {
		slog.Debug("Query executed",
			slog.String("query", name),
			slog.Duration("duration", time.Since(start)),
		)
	}()
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) Stats() sql.DBStats {
	return s.db.Stats()
}

func (s *PostgresStore) GetDB() *sql.DB {
	return s.db
}
