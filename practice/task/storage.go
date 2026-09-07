package main

import (
	"context"
	"database/sql"
	"time"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Connection pool config
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(ctx context.Context, task *Task) error {
	var CompletedAt sql.NullTime

	query := `INSERT into tasks(description) 
			VALUES ($1)
			RETURNING 
			id, description, completed, created_at, completed_at`

	err := s.db.QueryRowContext(ctx, query, task.Description).Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt, &CompletedAt)
	if err != nil {
		return err
	}

	if CompletedAt.Valid {
		task.CompletedAt = &CompletedAt.Time
	}

	return nil
}

func (s *PostgresStore) Get(ctx context.Context, id int) (*Task, error) {
	query := `SELECT id, description, completed, created_at, completed_at 
              FROM tasks WHERE id = $1`

	var task Task
	var completedAt sql.NullTime // lowercase variable names

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Description, &task.Completed,
		&task.CreatedAt, &completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, TaskNotFoundError{ID: id} // custom error
	}
	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil // explicit nil
}

func (s *PostgresStore) List(ctx context.Context) ([]*Task, error) {
	query := `SELECT id, description, completed, created_at, completed_at 
		FROM tasks`

	var tasks []*Task

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		var CompletedAt sql.NullTime

		err := rows.Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt, &CompletedAt)
		if err != nil {
			return nil, err
		}

		if CompletedAt.Valid {
			task.CompletedAt = &CompletedAt.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

func (s *PostgresStore) Complete(ctx context.Context, id int) (*Task, error) {
	task, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Complete()

	query := `
		UPDATE tasks
		SET completed = $1, completed_at = $2
		WHERE id = $3`

	_, err = s.db.ExecContext(ctx, query, task.Completed, task.CompletedAt, task.ID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *PostgresStore) Delete(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id) // lowercase
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return TaskNotFoundError{ID: id} // consistent custom error
	}

	return nil
}
