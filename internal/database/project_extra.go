package database

import (
	"context"
	"time"
)

const updateProjectByName = `
UPDATE projects
SET status = $2, link = $3, description = $4, updated_at = $5
WHERE name = $1
`

type UpdateProjectByNameParams struct {
	Name        string
	Status      string
	Link        string
	Description string
	UpdatedAt   time.Time
}

func (q *Queries) UpdateProjectByName(ctx context.Context, arg UpdateProjectByNameParams) error {
	_, err := q.db.ExecContext(ctx, updateProjectByName,
		arg.Name,
		arg.Status,
		arg.Link,
		arg.Description,
		arg.UpdatedAt,
	)
	return err
}
