# `bosla configure`

Manage configuration for bosla.

Aliases: `bosla config`

```
bosla configure
├── set     Manage configuration for bosla (writes ~/.bosla.json)
└── init    Initialize database for the CLI tool (runs migrations)
```

## `configure set`

Write (or merge) the connection settings stored in `~/.bosla.json`.
This command does **not** connect to the database.

```
bosla configure set -u <username> -p <password> [--url <url|127.0.0.1>]
```

| Flag         | Short | Default     | Description                                   |
| ------------ | ----- | ----------- | --------------------------------------------- |
| `--username` | `-u`  | —           | Database username (required for a valid config). |
| `--password` | `-p`  | —           | Database password.                            |
| `--database` | `-d`  | `bosla`     | Database name (stored in the config file).    |
| `--url`      | —     | `127.0.0.1` | Database host.                                |
| `--port`     | —     | `5432`      | Database port.                                |
| `--skip-ssl` | —     | `false`     | Connect with `sslmode=disable` instead of `sslmode=require`. |

Example:

```bash
bosla configure set \
  -u postgre \
  -p MyVeryStrongPassword \
  --url 127.0.0.1 \
  --port 5432 \
  --skip-ssl
```

Values that are left empty are kept from the existing config file, so you can
update a single field at a time. See
[Configuration](../reference/configuration.md) for the file format.

## `configure init`

Run the database migrations and create all bosla tables.

```bash
bosla configure init
```

The command connects to the database using the usual
[connection resolution](overview.md#database-connection-resolution)
(`BOSLA_DATABASE_URL` first, then `~/.bosla.json`).

!!! tip
    Migrations are idempotent — the schema uses `CREATE TABLE IF NOT EXISTS`,
    so you can safely re-run `configure init`.
