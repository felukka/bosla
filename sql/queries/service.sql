-- name: CreateService :one
INSERT INTO services (project_id, name, status, description, repo_url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateServiceStatus :exec
UPDATE services
SET status = $2, updated_at = $3
WHERE id = $1;

-- name: GetserviceByName :one
SELECT * FROM services
WHERE name = $1 AND project_id = $2
LIMIT 1;

-- name: GetserviceById :one
SELECT * FROM services
WHERE id = $1;

-- name: GetServiceByProjectName :many
SELECT * FROM services
WHERE project_id = (
  SELECT id FROM projects WHERE projects.id = sqlc.arg(project_id)
)
ORDER BY name;

-- name: GetProjectServiceByStatus :many
SELECT * FROM services
WHERE project_id = $1 AND status = $2
ORDER BY name;

-- name: DeleteProjectServiceByName :one
DELETE FROM services
WHERE name = $1 AND project_id = $2
RETURNING *;

-- name: CheckServiceExistsByName :one
SELECT EXISTS (
    SELECT 1 FROM services
    WHERE name = $1 AND project_id = $2
) AS exists;
