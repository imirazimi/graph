package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/imirazimi/graph/internal/infra/postgres"
	"github.com/imirazimi/graph/internal/task/entity"
)

type postgresRepository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) TaskRepository {
	return &postgresRepository{conn: conn}
}

func (r *postgresRepository) Create(ctx context.Context, task *entity.Task) error {
	query := `
        INSERT INTO tasks (
            id,
            title,
            description,
            status,
            assignee,
            created_at,
            updated_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7)
    `

	_, err := r.conn.Exec(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Assignee,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	query := `
        SELECT
            id,
            title,
            description,
            status,
            assignee,
            created_at,
            updated_at
        FROM tasks
        WHERE id = $1
    `

	var task entity.Task

	err := r.conn.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Assignee,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *postgresRepository) List(ctx context.Context, filter TaskFilter) ([]entity.Task, int64, error) {

	baseQuery := `
		SELECT
			id,
			title,
			description,
			status,
			assignee,
			created_at,
			updated_at
		FROM tasks
	`

	countQuery := `
		SELECT COUNT(*)
		FROM tasks
	`

	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// filters
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Assignee != "" {
		conditions = append(conditions, fmt.Sprintf("assignee = $%d", argIndex))
		args = append(args, filter.Assignee)
		argIndex++
	}

	// WHERE
	if len(conditions) > 0 {
		where := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += where
		countQuery += where
	}

	// pagination only for data query
	baseQuery += fmt.Sprintf(
		" ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		argIndex,
		argIndex+1,
	)

	argsForData := append([]interface{}{}, args...)
	argsForData = append(argsForData, filter.Limit, filter.Offset)

	// 1) get total
	var total int64
	err := r.conn.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 2) get data
	rows, err := r.conn.Query(ctx, baseQuery, argsForData...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []entity.Task

	for rows.Next() {
		var task entity.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Assignee,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func (r *postgresRepository) Update(ctx context.Context, task *entity.Task) error {
	query := `
        UPDATE tasks
        SET
            title = $2,
            description = $3,
            status = $4,
            assignee = $5,
            updated_at = $6
        WHERE id = $1
    `

	_, err := r.conn.Exec(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Assignee,
		task.UpdatedAt,
	)

	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, id)

	return err
}
