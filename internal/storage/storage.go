package storage

import (
	"context"
	"database/sql"
	"github.com/ankodd/todo-server/pkg/models/todo"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) (*Storage, error) {
	q := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY,
		name VARCHAR NOT NULL,
		done BIT NOT NULL
	)`

	_, err := db.Exec(q)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Insert(ctx context.Context, todo *todo.Todo) error {
	q := `INSERT INTO todos (name, done) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, q, todo.Name, false)
	return err
}

func (s *Storage) FetchAll(ctx context.Context) ([]todo.Todo, error) {
	q := `SELECT * FROM todos`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	var todos []todo.Todo

	for rows.Next() {
		var t todo.Todo

		err := rows.Scan(&t.ID, &t.Name, &t.Done)
		if err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	return todos, nil
}

func (s *Storage) Update(ctx context.Context, todo *todo.Todo, ID int64) error {
	q := `UPDATE todos SET name = ?, done = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, q, todo.Name, todo.Done, ID)

	return err
}

func (s *Storage) Delete(ctx context.Context, ID int64) error {
	q := `DELETE FROM todos WHERE id = ?`
	_, err := s.db.ExecContext(ctx, q, ID)

	return err
}

func (s *Storage) CountEntries(ctx context.Context) (int64, error) {
	q := `SELECT COUNT(*) FROM todos`
	var count int64

	err := s.db.QueryRowContext(ctx, q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
