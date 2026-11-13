# Quickstart

This guide takes you from an empty machine to your first release in under
five minutes.

## 1. Start a local PostgreSQL

```bash
docker compose -f config/docker-compose.yaml up -d
```

## 2. Point Bosla at the database

Bosla resolves the database connection in this order:

1. `BOSLA_DATABASE_URL` environment variable
2. `~/.bosla.json` (written by `bosla configure set`)

=== "Environment variable"

    ```bash
    export BOSLA_DATABASE_URL='postgresql://postgre:MyVeryStrongPassword@127.0.0.1:5432/bosla?sslmode=disable'
    ```

=== "Config file"

    ```bash
    bosla configure set \
      --username postgre \
      --password MyVeryStrongPassword \
      --url 127.0.0.1 \
      --port 5432 \
      --skip-ssl
    ```

!!! tip
    The environment variable always wins. `configure set` is handy when you
    do not want the URL in your shell history or dotfiles.

## 3. Initialize the schema

```bash
bosla configure init
```

This runs all migrations and creates the bosla tables.

## 4. Create a project

```bash
bosla project create my-example \
  --status active \
  --link https://github.com/felukka/bosla \
  --description "My example project"
```

## 5. Register services

```bash
bosla service register my-example koptan \
  --link https://github.com/felukka/koptan

bosla service register my-example khareeta \
  --link https://github.com/felukka/khareeta
```

## 6. Publish a release

```bash
bosla release release my-example v1.0.0 \
  --description "First production release"
```

With no `--service` flags, **every registered service** of the project is
attached to the release at version `v1.0.0`.

Pin individual services instead:

```bash
bosla release release my-example v1.1.0 \
  --service koptan=v2.0.1 \
  --service khareeta=v1.8.4 \
  --status stable \
  --git-hash 9f42c1a
```

!!! note
    The publish subcommand is named `release`, i.e.
    `bosla release release <project> <version>`. There is no `publish`
    subcommand.

## 7. Inspect

```bash
bosla project list -o table
bosla service list my-example
bosla release list my-example
bosla release describe my-example v1.0.0 -o json
```

Every list/describe command accepts `-o table` (default) or `-o json`.
