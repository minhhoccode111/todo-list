package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/internal/repo/persistent/sqlc"
	"github.com/minhhoccode111/todo-list/pkg/postgres"
)

// TodoRepo -.
type TodoRepo struct {
	*postgres.Postgres
	queries *sqlc.Queries
}

// NewTodoRepo -.
func NewTodoRepo(pg *postgres.Postgres) *TodoRepo {
	return &TodoRepo{
		Postgres: pg,
		queries:  sqlc.New(pg.Pool),
	}
}

func toPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}

	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}

func fromPgTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}

	return &t.Time
}

func newTodoFromDTO(t *sqlc.Todo) *entity.Todo {
	return &entity.Todo{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		Priority:    (*entity.PriorityLevel)(&t.Priority),
		DueDate:     fromPgTimestamptz(t.DueDate),
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func (tr *TodoRepo) CreateTodo(
	c context.Context,
	t *entity.Todo,
) (*entity.Todo, error) {
	sqlcTodo, err := tr.queries.CreateTodo(c, sqlc.CreateTodoParams{
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    sqlc.PriorityLevel(*t.Priority),
		DueDate:     toPgTimestamptz(t.DueDate),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, entity.ErrUnauthorized
		}

		return nil, fmt.Errorf("TodoRepo - CreateTodo - tr.queries.CreateTodo: %w", err)
	}

	return newTodoFromDTO(&sqlcTodo), nil
}

func (tr *TodoRepo) ReadTodoByID(_ context.Context, _ int32) (*entity.Todo, error) {
	return nil, nil
}
