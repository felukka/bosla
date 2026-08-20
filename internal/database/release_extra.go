package database

import (
	"context"
	"database/sql"
	"time"
)

const checkProjectVersionExists = `
SELECT EXISTS (
	SELECT 1 FROM project_versions WHERE project_id = $1 AND version = $2
)
`

type CheckProjectVersionExistsParams struct {
	ProjectID int32
	Version   string
}

func (q *Queries) CheckProjectVersionExists(ctx context.Context, arg CheckProjectVersionExistsParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkProjectVersionExists, arg.ProjectID, arg.Version)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createProjectVersion = `
INSERT INTO project_versions (project_id, version, description, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, project_id, version, description, created_at
`

type CreateProjectVersionParams struct {
	ProjectID   int32
	Version     string
	Description string
	CreatedAt   time.Time
}

func (q *Queries) CreateProjectVersion(ctx context.Context, arg CreateProjectVersionParams) (ProjectVersion, error) {
	row := q.db.QueryRowContext(ctx, createProjectVersion, arg.ProjectID, arg.Version, arg.Description, arg.CreatedAt)
	var pv ProjectVersion
	err := row.Scan(&pv.ID, &pv.ProjectID, &pv.Version, &pv.Description, &pv.CreatedAt)
	return pv, err
}

const getProjectVersionByProjectAndVersion = `
SELECT id, project_id, version, description, created_at
FROM project_versions
WHERE project_id = $1 AND version = $2
LIMIT 1
`

type GetProjectVersionByProjectAndVersionParams struct {
	ProjectID int32
	Version   string
}

func (q *Queries) GetProjectVersionByProjectAndVersion(ctx context.Context, arg GetProjectVersionByProjectAndVersionParams) (ProjectVersion, error) {
	row := q.db.QueryRowContext(ctx, getProjectVersionByProjectAndVersion, arg.ProjectID, arg.Version)
	var pv ProjectVersion
	err := row.Scan(&pv.ID, &pv.ProjectID, &pv.Version, &pv.Description, &pv.CreatedAt)
	return pv, err
}

const listProjectVersions = `
SELECT id, project_id, version, description, created_at
FROM project_versions
WHERE project_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListProjectVersions(ctx context.Context, projectID int32) ([]ProjectVersion, error) {
	rows, err := q.db.QueryContext(ctx, listProjectVersions, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ProjectVersion, 0)
	for rows.Next() {
		var pv ProjectVersion
		if err := rows.Scan(&pv.ID, &pv.ProjectID, &pv.Version, &pv.Description, &pv.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, pv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const checkServiceVersionExists = `
SELECT EXISTS (
	SELECT 1 FROM service_versions WHERE service_id = $1 AND version = $2
)
`

type CheckServiceVersionExistsParams struct {
	ServiceID int32
	Version   string
}

func (q *Queries) CheckServiceVersionExists(ctx context.Context, arg CheckServiceVersionExistsParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkServiceVersionExists, arg.ServiceID, arg.Version)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createServiceVersion = `
INSERT INTO service_versions (service_id, version, status, git_hash, description, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, service_id, version, status, git_hash, description, created_at
`

type CreateServiceVersionParams struct {
	ServiceID   int32
	Version     string
	Status      string
	GitHash     string
	Description string
	CreatedAt   time.Time
}

func (q *Queries) CreateServiceVersion(ctx context.Context, arg CreateServiceVersionParams) (ServiceVersion, error) {
	row := q.db.QueryRowContext(ctx, createServiceVersion, arg.ServiceID, arg.Version, arg.Status, arg.GitHash, arg.Description, arg.CreatedAt)
	var sv ServiceVersion
	err := row.Scan(&sv.ID, &sv.ServiceID, &sv.Version, &sv.Status, &sv.GitHash, &sv.Description, &sv.CreatedAt)
	return sv, err
}

const getServiceVersionByServiceAndVersion = `
SELECT id, service_id, version, status, git_hash, description, created_at
FROM service_versions
WHERE service_id = $1 AND version = $2
LIMIT 1
`

type GetServiceVersionByServiceAndVersionParams struct {
	ServiceID int32
	Version   string
}

func (q *Queries) GetServiceVersionByServiceAndVersion(ctx context.Context, arg GetServiceVersionByServiceAndVersionParams) (ServiceVersion, error) {
	row := q.db.QueryRowContext(ctx, getServiceVersionByServiceAndVersion, arg.ServiceID, arg.Version)
	var sv ServiceVersion
	err := row.Scan(&sv.ID, &sv.ServiceID, &sv.Version, &sv.Status, &sv.GitHash, &sv.Description, &sv.CreatedAt)
	return sv, err
}

const upsertProjectVersionService = `
INSERT INTO project_version_apps (project_version_id, service_id, service_version_id)
VALUES ($1, $2, $3)
ON CONFLICT (project_version_id, service_id)
DO UPDATE SET service_version_id = EXCLUDED.service_version_id
`

type UpsertProjectVersionServiceParams struct {
	ProjectVersionID int32
	ServiceID        int32
	ServiceVersionID int32
}

func (q *Queries) UpsertProjectVersionService(ctx context.Context, arg UpsertProjectVersionServiceParams) error {
	_, err := q.db.ExecContext(ctx, upsertProjectVersionService, arg.ProjectVersionID, arg.ServiceID, arg.ServiceVersionID)
	return err
}

type ReleaseDetailsRow struct {
	ProjectName        string
	ReleaseVersion     string
	ReleaseDescription string
	ReleaseCreatedAt   time.Time
	ServiceName        sql.NullString
	ServiceVersion     sql.NullString
	ServiceStatus      sql.NullString
	ServiceGitHash     sql.NullString
	ServiceDescription sql.NullString
}

const getReleaseDetails = `
SELECT
	p.name,
	pv.version,
	pv.description,
	pv.created_at,
	s.name,
	sv.version,
	sv.status,
	sv.git_hash,
	sv.description
FROM project_versions pv
JOIN projects p ON p.id = pv.project_id
LEFT JOIN project_version_apps pva ON pva.project_version_id = pv.id
LEFT JOIN services s ON s.id = pva.service_id
LEFT JOIN service_versions sv ON sv.id = pva.service_version_id
WHERE pv.project_id = $1 AND pv.version = $2
ORDER BY s.name
`

type GetReleaseDetailsParams struct {
	ProjectID int32
	Version   string
}

func (q *Queries) GetReleaseDetails(ctx context.Context, arg GetReleaseDetailsParams) ([]ReleaseDetailsRow, error) {
	rows, err := q.db.QueryContext(ctx, getReleaseDetails, arg.ProjectID, arg.Version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ReleaseDetailsRow, 0)
	for rows.Next() {
		var item ReleaseDetailsRow
		if err := rows.Scan(
			&item.ProjectName,
			&item.ReleaseVersion,
			&item.ReleaseDescription,
			&item.ReleaseCreatedAt,
			&item.ServiceName,
			&item.ServiceVersion,
			&item.ServiceStatus,
			&item.ServiceGitHash,
			&item.ServiceDescription,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
