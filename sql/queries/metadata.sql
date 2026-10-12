-- name: GetProjectMetadata :many
SELECT project_id, key, value, created_at, updated_at
FROM project_metadata
WHERE project_id = $1
ORDER BY key;

-- name: GetServiceMetadata :many
SELECT service_id, key, value, created_at, updated_at
FROM service_metadata
WHERE service_id = $1
ORDER BY key;

-- name: GetMetadataByKey :one
SELECT project_id, key, value, created_at, updated_at
FROM project_metadata
WHERE project_id = $1 AND key = $2;

-- name: GetServiceMetadataByKey :one
SELECT service_id, key, value, created_at, updated_at
FROM service_metadata
WHERE service_id = $1 AND key = $2;

-- name: UpsertProjectMetadata :exec
INSERT INTO project_metadata (project_id, key, value, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
ON CONFLICT (project_id, key)
DO UPDATE SET value = EXCLUDED.value, updated_at = now();

-- name: UpsertServiceMetadata :exec
INSERT INTO service_metadata (service_id, key, value, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
ON CONFLICT (service_id, key)
DO UPDATE SET value = EXCLUDED.value, updated_at = now();

-- name: DeleteProjectMetadata :exec
DELETE FROM project_metadata
WHERE project_id = $1 AND key = $2;

-- name: DeleteServiceMetadata :exec
DELETE FROM service_metadata
WHERE service_id = $1 AND key = $2;
