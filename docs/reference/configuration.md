# Configuration

Bosla reads its database connection from one of two places, in order:

1. **`BOSLA_DATABASE_URL`** — a full connection string, used verbatim.
2. **`~/.bosla.json`** — written by `bosla configure set`.

If `BOSLA_DATABASE_URL` is set, the config file is never read.

## Connection string

```
postgresql://<user>:<password>@<host>:<port>/bosla?sslmode=disable
```

The scheme `postgres://` and `postgresql://` are both accepted.

## Config file

`bosla configure set` maintains a JSON file at `~/.bosla.json`:

```json
{
  "username": "postgre",
  "password": "MyVeryStrongPassword",
  "database": "bosla",
  "url": "127.0.0.1",
  "port": 5432,
  "skipSsl": true
}
```

| Field      | From flag    | Meaning                                          |
| ---------- | ------------ | ------------------------------------------------ |
| `username` | `--username` | Database user. Required for a valid config.      |
| `password` | `--password` | Database password.                               |
| `database` | `--database` | Stored for reference (see note below).           |
| `url`      | `--url`      | Database host. Required for a valid config.      |
| `port`     | `--port`     | Database port; `0` falls back to `5432`.         |
| `skipSsl`  | `--skip-ssl` | `sslmode=disable` when true, `sslmode=require` when false. |

When the file is built from the config, Bosla connects to:

```
postgres://<username>:<password>@<url>:<port>/bosla?sslmode=<require|disable>
```

!!! note "Database name"
    The connection string built from `~/.bosla.json` always uses the
    database name **`bosla`**. The `--database` flag value is persisted in
    the JSON file, but the current URL template does not substitute it.
    If you need a different database name, use `BOSLA_DATABASE_URL`.

## Validity rules

The config file is considered invalid when it is missing, unreadable as
JSON, or when `username` / `url` are empty. Commands that need a database
then fail with a hint to set `BOSLA_DATABASE_URL` or run
`bosla configure set --help`.
