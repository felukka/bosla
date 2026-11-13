# Prerequisites

Bosla is a single Go binary plus a PostgreSQL database.

## Required

| Requirement  | Notes                                                          |
| ------------ | -------------------------------------------------------------- |
| PostgreSQL   | Any recent version works; the dev environment uses **PostgreSQL 16**. |
| Go 1.24+     | Only if you build from source.                                 |
| Nix          | Only if you install via the flake (recommended for dev).       |

## Optional

| Tool          | Notes                                                                  |
| ------------- | ---------------------------------------------------------------------- |
| Docker        | Convenient way to run a local PostgreSQL instance (see Quickstart).    |

## Local database

The repository ships a compose file that starts a ready-to-use PostgreSQL
instance:

```bash
docker compose -f config/docker-compose.yaml up -d
```

| Setting    | Value                          |
| ---------- | ------------------------------ |
| Host       | `127.0.0.1`                    |
| Port       | `5432`                         |
| User       | `postgre`                      |
| Password   | `MyVeryStrongPassword`         |
| Database   | `bosla`                        |

!!! warning
    These credentials are for **local development only**. Do not reuse them
    in shared or production environments.
