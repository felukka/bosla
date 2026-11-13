# `bosla release`

Manage project releases — a version of the project plus the service
versions that ship in it.

```
bosla release
├── list        List releases for a project (alias: ls)
├── describe    Describe a release
└── release     Publish a release and attach service versions
```

!!! note
    The publish subcommand is named **`release`**, so publishing looks like
    `bosla release release <project> <version>`. There is no `publish`
    subcommand.

## `release list`

```bash
bosla release list my-example
bosla release list my-example -o json
```

## `release describe`

```bash
bosla release describe my-example v1.0.0
bosla release describe my-example v1.0.0 -o json
```

## `release release` (publish)

```
bosla release release <project> <version>
    [--description <text>]
    [--status <s>]
    [--git-hash <hash>]
    [--service <name>=<version>]   # repeatable
```

| Flag            | Short | Default      | Description                                              |
| --------------- | ----- | ------------ | -------------------------------------------------------- |
| `--description` | `-d`  | —            | Release description.                                     |
| `--status`      | `-s`  | `released`   | Status written on newly created **service versions**.    |
| `--git-hash`    | —     | —            | Git hash written on newly created service versions.      |
| `--service`     | —     | —            | Attach a service at a specific version, `name=version`.  Repeatable. |

### Version assignment rules

* **No `--service` flags** → every service registered on the project is
  attached to the release at the *release version*.
* **With `--service` flags** → only the listed services are attached, each at
  its given version. Services not listed are skipped.
* A `--service` name that does not belong to the project fails the command.
* Service versions that **already exist** are reused as-is — `--status`,
  `--git-hash`, and `--description` are only applied to versions that are
  created by this command.
* A project version is unique per project; publishing the same version twice
  fails with `release "<version>" already exists for project "<name>"`.

### Examples

Publish with all services at the same version:

```bash
bosla release release my-example v1.1.0 \
  --description "July release"
```

Publish with explicit per-service versions:

```bash
bosla release release my-example v1.2.0 \
  --service koptan=v2.0.1 \
  --service khareeta=v1.8.4 \
  --status stable \
  --git-hash 9f42c1a
```
