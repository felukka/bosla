# `bosla service`

Manage services that belong to a project. Alias: `bosla svc`

```
bosla service
├── register    Register a new service to a project
├── list        List services for a project (alias: ls)
├── describe    Describe a service
├── update      Update service fields
├── delete      Delete a service (alias: rm)
├── search      Search services in a project with a regex pattern
└── metadata    Manage service metadata (aliases: meta, m)
    ├── set
    ├── get
    ├── list
    └── delete
```

## `service register`

```
bosla service register <project> <service> [--status <s>] [--link <url>] [--desc <text>]
```

| Flag        | Short | Default  | Description                        |
| ----------- | ----- | -------- | ---------------------------------- |
| `--status`  | `-s`  | `active` | Service status.                    |
| `--link`    | `-l`  | —        | Service repository URL.            |
| `--desc`    | `-d`  | —        | Service description.               |

```bash
bosla service register my-example koptan \
  --link https://github.com/felukka/koptan
```

The project must exist, and a service name is unique **per project**
(you can have `koptan` under two different projects).

## `service list`

```bash
bosla service list my-example
bosla service list my-example -o json
```

## `service describe`

```bash
bosla service describe my-example koptan
bosla service describe my-example koptan -o json
```

## `service update`

Only the flags you pass are changed.

```
bosla service update <project> <service> [--status <s>] [--link <url>] [--desc <text>]
```

```bash
bosla service update my-example koptan --link https://github.com/felukka/koptan-v2
```

## `service delete`

Deletes the service and (via database cascade) its versions and metadata:

```bash
bosla service delete my-example koptan
```

!!! warning
    Unlike `project delete`, service deletion has **no confirmation flag**.

## `service search`

Regex search, scoped to a single project:

```bash
bosla service search my-example 'gateway|billing|worker-.*'
```

## `service metadata`

Free-form key/value pairs on a service. Aliases: `meta`, `m`.

```
bosla service metadata set <project> <service> <key> <value>
bosla service metadata get <project> <service> <key>
bosla service metadata list <project> <service>            # alias: ls
bosla service metadata delete <project> <service> <key>    # aliases: rm, del
```

```bash
bosla service metadata set my-example koptan owner team-platform
bosla service metadata get my-example koptan owner
bosla service metadata delete my-example koptan owner
```
