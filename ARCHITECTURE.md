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

### gtkai.json Manifest

Every filter module ships a `gtkai.json` manifest at the repository root. Required fields:

| Field | Description |
|-------|-------------|
| `id` | Unique module identifier (validated format) |
| `filters` | List of argv0 commands handled |
| `platforms` | Supported platforms |
| `contract` | Must be `subprocess/v1` |
| `gtkai-core-version` | Core compatibility requirement |

Module version is **not** declared in the manifest. It is resolved from the Git tag at install time:

```bash
gtkai filter install gtk-ai/gtkai-date@v0.10.0
```

#### gtkai-core-version

```json
{
  "version": "0.10.0",
  "constraint": "min"
}
```

| `constraint` | Semantics |
|--------------|-----------|
| `"min"` | Running `gtkai` must be `>= version` |
| `"exact"` | Running `gtkai` must equal `version` |

Both fields are required. Unknown constraints abort installation.

### Validation on Install

During installation, the following checks are enforced in order:

1. Git tag/ref resolves to a valid semver (module version)
2. `gtkai.json` is present and parses without error
3. `id` format validation
4. `contract` equals `subprocess/v1`
5. Platform compatibility check
6. `gtkai-core-version.version` is valid semver
7. `gtkai-core-version.constraint` is `"min"` or `"exact"`
8. Running `gtkai` satisfies the core version constraint
9. Liveness check: binary must respond within **500ms**

### Persistence

Installed filters are recorded in `~/.gtk-ai/filters.db` (SQLite). The path is resolved via `internal/storage.Dir()`. The installed module version is stored from the Git tag, not from the manifest.

### Binary Location

Filter binaries are stored under `~/.gtk-ai/filters/<id>/`, downloaded from GitHub Releases of the corresponding filter repository.

```
~/.gtk-ai/filters/gtk-ai/gtkai-date/
    gtkai-date
    gtkai.json
```

### Official Filters & install.sh

`install.sh` downloads official filters from repositories matching the `gtk-ai/gtkai-*` pattern by default.

The reference template repository for building external filters is [`gtk-ai/gtkai-date`](https://github.com/gtk-ai/gtkai-date).

### Built-ins Fallback

Built-in filters remain compiled into the core binary and act as fallback when no external filter is installed for a given operation.

### Pending

- `gain` attribution by `filter_id` is pending implementation.
