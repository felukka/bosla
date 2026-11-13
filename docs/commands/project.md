# `bosla project`

Manage projects. Alias: `bosla prj`

```
bosla project
├── create      Create a new project
├── list        List projects (alias: ls)
├── describe    Describe a project
├── update      Update project fields
├── delete      Delete a project (alias: rm)
├── search      Search projects with a regex pattern
└── metadata    Manage project metadata (aliases: meta, m)
    ├── set
    ├── get
    ├── list
    └── delete
```

## `project create`

```
bosla project create <name> [--status <s>] [--link <url>] [--desc <text>]
```

| Flag        | Short | Default  | Description        |
| ----------- | ----- | -------- | ------------------ |
| `--status`  | `-s`  | `active` | Project status.    |
| `--link`    | `-l`  | —        | Project link (e.g. repo URL). |
| `--desc`    | `-d`  | —        | Project description. |

```bash
bosla project create my-example \
  --status active \
  --link https://github.com/felukka/bosla \
  --description "Bosla"
```

Creating a project that already exists fails with
`project "<name>" already exists`.

## `project list`

```bash
bosla project list
bosla project list -o json
```

## `project describe`

```bash
bosla project describe my-example
bosla project describe my-example -o json
```

Fails with `project "<name>" not found` if the project does not exist.

## `project update`

Update any subset of the mutable fields — only the flags you pass are
changed, the rest are left as-is.

```
bosla project update <name> [--status <s>] [--link <url>] [--description <text>]
```

| Flag          | Short | Description         |
| ------------- | ----- | ------------------- |
| `--status`    | `-s`  | New status.         |
| `--link`      | `-l`  | New link.           |
| `--description` | `-d` | New description.    |

```bash
bosla project update my-example --status archived
```

## `project delete`

Deleting a project is **destructive and cascades** to its services,
releases, service versions, and metadata. It requires explicit confirmation:

```bash
bosla project delete my-example --yes-i-am-sure
```

Without `--yes-i-am-sure` the command fails with:

```
deletion requires explicit confirmation: use --yes-i-am-sure flag
```

## `project search`

Match project names against a **regular expression** (Go regexp syntax,
matched against the whole name).

```bash
bosla project search '^(core|api)-[0-9]+$'
```

Invalid patterns are reported before any database query runs.

## `project metadata`

Free-form key/value pairs on a project. Aliases: `meta`, `m`.

```
bosla project metadata set <project> <key> <value>
bosla project metadata get <project> <key>
bosla project metadata list <project>          # alias: ls
bosla project metadata delete <project> <key>  # aliases: rm, del
```

```bash
bosla project metadata set my-example owner team-platform
bosla project metadata get my-example owner
bosla project metadata list my-example -o json
bosla project metadata delete my-example owner
```

`set` upserts: setting an existing key updates its value and `updated_at`.
