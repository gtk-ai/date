# Architecture

## §4 — External Filters

### Transport: subprocess/v1

Each external filter is a standalone binary that communicates with the core via stdin/stdout using JSON. This is the only supported transport. Go plugins (`.so`) are explicitly out of scope.

### JSON Contract

**Request** sent to the filter binary:

```json
{
  "operation": "rewrite" | "filter_output",
  "args": [...],
  "output": "...",
  "exit_code": 0
}
```

**Response** expected from the filter binary:

```json
{
  "args": [...],
  "changed": true,
  "output": "..."
}
```

### filter.json Manifest

Every filter package must include a `filter.json` manifest at its root. Required fields:

| Field               | Description                                      |
|---------------------|--------------------------------------------------|
| `id`                | Unique filter identifier (validated format)      |
| `filters`           | List of filter definitions                       |
| `version`           | Semver string                                    |
| `platforms`         | Supported platforms                              |
| `contract`          | Must be `subprocess/v1`                          |
| `min_gtkai_version` | Minimum compatible gtkai version (semver)        |

### Validation on Install

During installation, the following checks are enforced in order:

1. `id` format validation
2. `version` is valid semver
3. `contract` equals `subprocess/v1`
4. Platform compatibility check
5. Liveness check: binary must respond within **500ms**

### Persistence

Installed filters are recorded in `~/.gtk-ai/filters.db` (SQLite). The path is resolved via `internal/storage.Dir()`.

### Binary Location

Filter binaries are stored under `~/.gtk-ai/filters/<id>/`, downloaded from GitHub Releases of the corresponding filter repository.

### Official Filters & install.sh

`install.sh` downloads official filters from repositories matching the `gtk-ai/gtkai-*` pattern by default.

The reference template repository for building external filters is [`gtk-ai/gtkai-date`](https://github.com/gtk-ai/gtkai-date).

### Built-ins Fallback

Built-in filters remain compiled into the core binary and act as fallback when no external filter is installed for a given operation.

### Pending

- `gain` attribution by `filter_id` is pending implementation.
