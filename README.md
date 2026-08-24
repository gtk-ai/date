# date

A [gtk-ai](https://github.com/gtk-ai/gtkai-core) plugin for the `date` command. It intercepts calls to `date`, rewrites arguments when needed, and filters output before it reaches the agent.

This repository also works as a **template** for building new gtk-ai marketplace plugins. The structure, contract, and CI setup are meant to be copied.

## Install

```bash
gtkai filter install github.com/gtk-ai/date@v0.13.0
```

Requires gtkai-core >= 0.11.0.

## What it does

The plugin implements two operations defined by the `stdin/v1` contract:

- `rewrite`: receives the original arguments before `date` runs and can return a modified set
- `filter_output`: receives the raw output after `date` runs and can return a cleaned version

Both are optional to override — if the plugin returns `changed: false`, the original value passes through unchanged.

## Contract

This plugin uses the `stdin/v1` protocol: gtkai-core sends a JSON object to the plugin's stdin, the plugin writes a JSON response to stdout, and exits 0.

Request shape:

```json
{
  "operation": "rewrite" | "filter_output",
  "args": ["..."],
  "output": "...",
  "exit_code": 0
}
```

Response shape:

```json
{
  "changed": true,
  "args": ["..."],
  "output": "..."
}
```

See [HOWTO.md](HOWTO.md) for a step-by-step guide to building a new plugin from this template.

## Plugin metadata

```json
{
  "id": "gtk-ai/date",
  "command": "date",
  "contract": "stdin/v1",
  "platforms": ["linux/amd64", "darwin/arm64"]
}
```

## Using this as a template

1. Fork or copy this repository
2. Replace `date` with the command you want to intercept in `gtkai.json` and `go.mod`
3. Implement `filter.Rewrite` and `filter.FilterOutput` in the `filter/` package
4. Update the module path in `go.mod` and `gtkai.json`
5. The CI workflow in `.github/workflows/release.yml` builds and publishes binaries automatically on tag push
