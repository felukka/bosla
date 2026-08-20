# Bosla CLI [WIP]

A serious, well-scoped versioning tool.

---

## What Bosla manages

- **Projects**: name, status, link, and description
- **Services**: services that belong to a project
- **Releases**: project release versions with linked service versions

---

## Quick start

### 1) Set your database connection

```bash
export BOSLA_DATABASE_URL='postgresql://<user>:<password>@localhost:5432/bosla?sslmode=disable'
```

### 2) Initialize schema

```bash
bosla configure init
```

### 3) Create a project

```bash
bosla project create my-example \
  --status active \
  --link https://github.com/felukka/bosla \
  --description "Bosla"
```

### 4) Add services

```bash
bosla service register my-example koptan \
  --link https://github.com/felukka/koptan

bosla service register my-example khareeta \
  --link https://github.com/felukka/khareeta
```

### 5) Create a release

```bash
bosla release publish my-example v1.0.0 \
  --description "First production release"
```

---

## Output format

Bosla supports:

- `-o table` (default)
- `-o json`

Examples:

```bash
bosla project list -o table
bosla project list -o json
```

---

## Command guide

### Project commands

```bash
bosla project create <name> [--status] [--link] [--description]
bosla project list
bosla project describe <name>
bosla project update <name> [--status] [--link] [--description]
bosla project delete <name> --yes-i-am-sure
bosla project search '<pattern>'
```

Pattern search example:

```bash
bosla project search '^(core|api)-[0-9]+$'
```

---

### Service commands

```bash
bosla service register <project> <service> [--status] [--link] [--description]
bosla service list <project>
bosla service delete <project> <service>
bosla service search <project> '<pattern>'
```

Pattern search example:

```bash
bosla service search my-example 'gateway|billing|worker-.*'
```

---

### Release commands

```bash
bosla release publish <project> <version> [--description] [--status] [--git-hash] [--service name=version]
bosla release list <project>
```

Release examples:

Publish a release with all services using the same version:

```bash
bosla release publish my-example v1.1.0 \
  --description "July release"
```

Create a release with explicit service versions:

```bash
bosla release publish my-example v1.2.0 \
  --service koptan=v2.0.1 \
  --service khareeta=v1.8.4 \
  --status stable \
  --git-hash 9f42c1a
```

Inspect a release:

```bash
bosla release describe my-example v1.2.0 -o json
```
