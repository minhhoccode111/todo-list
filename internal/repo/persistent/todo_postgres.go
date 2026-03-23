package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (tr *TodoRepo) ReadTodos(
	c context.Context,
	userID, limit, offset int32,
) ([]entity.Todo, int32, error) {
	sqlcTodos, err := tr.queries.ReadTodos(c, sqlc.ReadTodosParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("TodoRepo - ReadTodos - tr.queries.ReadTodos: %w", err)
	}

	if len(sqlcTodos) == 0 {
		return []entity.Todo{}, 0, nil
	}

	todos := make([]entity.Todo, 0, len(sqlcTodos))
	for i := range sqlcTodos {
		t := &sqlcTodos[i]
		todos = append(todos, entity.Todo{
			ID:          t.ID,
			UserID:      t.UserID,
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
			Priority:    (*entity.PriorityLevel)(&t.Priority),
			DueDate:     fromPgTimestamptz(t.DueDate),
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	return todos, int32(sqlcTodos[0].TotalCount), nil // #nosec G115
}

func (tr *TodoRepo) UpdateTodo(c context.Context, t *entity.Todo) (*entity.Todo, error) {
	sqlcTodo, err := tr.queries.UpdateTodo(c, sqlc.UpdateTodoParams{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		Priority:    sqlc.PriorityLevel(*t.Priority),
		DueDate:     toPgTimestamptz(t.DueDate),
		UserID:      t.UserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrForbidden
		}

		return nil, fmt.Errorf("TodoRepo - UpdateTodo - tr.queries.UpdateTodo: %w", err)
	}

	return newTodoFromDTO(&sqlcTodo), nil
}

func (tr *TodoRepo) DeleteTodo(c context.Context, todoID, userID int32) error {
	err := tr.queries.DeleteTodo(c, sqlc.DeleteTodoParams{
		ID:     todoID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ErrForbidden
		}

		return fmt.Errorf("TodoRepo - DeleteTodo - tr.queries.DeleteTodo: %w", err)
	}

	return nil
}
