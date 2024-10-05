package storage

import (
	"context"
	"github.com/ankodd/todo-server/pkg/models/todo"
)

type Storage interface {
	Insert(ctx context.Context, todo *todo.Todo) error
	FetchAll(ctx context.Context) ([]todo.Todo, error)
	Update(ctx context.Context, todo *todo.Todo, ID int64) error
	Delete(ctx context.Context, ID int64) error
	CountEntries(ctx context.Context) (int64, error)
}
