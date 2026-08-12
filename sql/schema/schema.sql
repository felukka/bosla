BEGIN;

CREATE TABLE IF NOT EXISTS projects (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    status          TEXT NOT NULL DEFAULT 'active',
    link            TEXT,
    description     TEXT,
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS services (
    id              SERIAL PRIMARY KEY,
    project_id      INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'active',
    description     TEXT,
    repo_url        TEXT,
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (project_id, name)
);

CREATE TABLE IF NOT EXISTS project_versions (
    id              SERIAL PRIMARY KEY,
    project_id      INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version         TEXT NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (project_id, version)
);

CREATE TABLE IF NOT EXISTS service_versions (
    id              SERIAL PRIMARY KEY,
    service_id      INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    version         TEXT NOT NULL,
    status          TEXT,
    git_hash        TEXT,
    description     TEXT,
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (service_id, version)
);

CREATE TABLE IF NOT EXISTS project_version_services (
    project_version_id INTEGER NOT NULL REFERENCES project_versions(id) ON DELETE CASCADE,
    service_id         INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    service_version_id INTEGER NOT NULL REFERENCES service_versions(id),

    PRIMARY KEY (project_version_id, service_id)
);

CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_metadata ON projects USING GIN (metadata);

CREATE INDEX IF NOT EXISTS idx_services_name ON services(name);
CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);
CREATE INDEX IF NOT EXISTS idx_services_metadata ON services USING GIN (metadata);

CREATE INDEX IF NOT EXISTS idx_project_versions_version ON project_versions(version);
CREATE INDEX IF NOT EXISTS idx_service_versions_version ON service_versions(version);
CREATE INDEX IF NOT EXISTS idx_service_versions_status ON service_versions(status);
CREATE INDEX IF NOT EXISTS idx_service_versions_metadata ON service_versions USING GIN (metadata);

COMMIT;
