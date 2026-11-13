# Command overview

Bosla is a single binary with four command groups:

```
bosla
├── configure (alias: config)   manage configuration
│   ├── set                     write ~/.bosla.json
│   └── init                    run database migrations
├── project (alias: prj)        manage projects
│   ├── create / list / describe / update / delete / search
│   └── metadata (meta, m)      set / get / list / delete
├── service (alias: svc)        manage services
│   ├── register / list / describe / update / delete / search
│   └── metadata (meta, m)      set / get / list / delete
└── release                     manage project releases
    ├── list
    ├── describe
    └── release                 publish a release (no `publish` subcommand)
```

Run `bosla <group> --help` at any time for exact flags and examples.

## Global flags

| Flag        | Default | Description                        |
| ----------- | ------- | ---------------------------------- |
| `-o, --output` | `table` | Output format: `table` or `json`.  |

The output flag is a **persistent** flag — it works on every command.

## Database connection resolution

Every command (except `configure set`, which only writes the config file)
opens a PostgreSQL connection using this precedence:

1. **`BOSLA_DATABASE_URL`** — if set, used verbatim.
2. **`~/.bosla.json`** — otherwise Bosla builds the URL as
   `postgres://<user>:<password>@<url>:<port>/bosla?sslmode=<require|disable>`.

If neither is available you get:

```
database connection is not set.
Set BOSLA_DATABASE_URL or run `bosla configure set --help`
```

See [Configuration](../reference/configuration.md) for details.

## Aliases

| Command     | Aliases         |
| ----------- | --------------- |
| `configure` | `config`        |
| `project`   | `prj`           |
| `service`   | `svc`           |
| `list`      | `ls` (in every group) |
| `delete`    | `rm` (in every group; `del` in metadata groups) |
| `metadata`  | `meta`, `m`     |
