# Bosla CLI

Bosla is a small command-line tool for tracking projects, services, and releases in one place.

It is built for day-to-day engineering work: create a project, attach services, and tag releases with clear output in table or JSON.

---

## What Bosla manages

- **Projects**: name, status, link, and description
- **Services**: services/apps that belong to a project
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
bosla project create my-platform --status active --link https://github.com/acme/platform --description "Main platform"
```

### 4) Add services

```bash
bosla service add my-platform gateway --link https://github.com/acme/gateway
bosla service add my-platform billing --link https://github.com/acme/billing
```

### 5) Create a release

```bash
bosla release create my-platform v1.0.0 --description "First production release"
```

---

## Output format

Bosla supports:

- `-o table` (default)
- `-o json`

Examples:

```bash
bosla project get -o table
bosla project get -o json
```

---

## Command guide

### Project commands

```bash
bosla project create <name> [--status] [--link] [--description]
bosla project get
bosla project describe <name>
bosla project update <name> [--status] [--link] [--description]
bosla project delete <name> --yes-i-am-sure
bosla project search '<regex>'
```

Regex search example:

```bash
bosla project search '^(core|api)-[0-9]+$'
```

---

### Service commands

```bash
bosla service add <project> <service> [--status] [--link] [--description]
bosla service get <project>
bosla service delete <project> <service>
bosla service search <project> '<regex>'
```

Regex search example:

```bash
bosla service search my-platform 'gateway|billing|worker-.*'
```

---

### Release commands

```bash
bosla release create <project> <release-version> [--description] [--status] [--git-hash] [--service name=version]
bosla release get <project> <release-version>
bosla release list <project>
```

Release examples:

Create a release with all services using the same version:

```bash
bosla release create my-platform v1.1.0 --description "July release"
```

Create a release with explicit service versions:

```bash
bosla release create my-platform v1.2.0 \
  --service gateway=v2.0.1 \
  --service billing=v1.8.4 \
  --status stable \
  --git-hash 9f42c1a
```

Inspect a release:

```bash
bosla release get my-platform v1.2.0 -o json
```

---

## Notes

- If a regex pattern is invalid, Bosla returns a clear error.
- If no records match a list/search command, output is `No resources found.`
- `release`, `version`, and `ver` can all be used for release commands.
