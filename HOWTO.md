# Cómo definir un módulo de filtro

Este repositorio es la **plantilla de referencia** para crear módulos externos de gtk-ai bajo la convención `gtk-ai/<cmd>`. Cópialo, renómbralo y adapta los puntos marcados con `<cmd>`.

## Resumen rápido

1. Duplica este repo como `gtk-ai/<cmd>` (por ejemplo `gtk-ai/ls`).
2. Ajusta `gtkai.json`, `go.mod`, `filter/` y `cmd/`.
3. Etiqueta una versión semver (`v0.1.0`).
4. Instala con `gtkai filter install github.com/gtk-ai/<cmd>@v0.1.0`.

## Estructura del repositorio

```
gtk-ai/<cmd>/
├── gtkai.json          # manifiesto del módulo (obligatorio)
├── go.mod              # module github.com/gtk-ai/<cmd>
├── cmd/
│   └── main.go         # binario subprocess/v1
├── filter/
│   └── <cmd>.go        # lógica Rewrite / FilterOutput
└── *_test.go           # tests de lógica y protocolo
```

En este repo concreto, `<cmd>` es `date` y el ID del filtro es `gtk-ai/date`.

## Paso 1 — Definir el manifiesto (`gtkai.json`)

Coloca `gtkai.json` en la raíz del repositorio. Campos obligatorios:

| Campo | Qué poner |
|-------|-----------|
| `id` | Identificador único con formato `author/<cmd>` |
| `filters` | Lista con el argv0 interceptado (p. ej. `["date"]`) |
| `platforms` | Plataformas soportadas (`linux/amd64`, `darwin/arm64`, …) |
| `contract` | Siempre `"subprocess/v1"` |
| `gtkai-core-version` | Versión mínima o exacta del core gtk-ai |

Ejemplo (este módulo):

```json
{
  "id": "gtk-ai/date",
  "filters": ["date"],
  "platforms": ["linux/amd64", "darwin/arm64"],
  "contract": "subprocess/v1",
  "gtkai-core-version": {
    "version": "0.10.0",
    "constraint": "min"
  }
}
```

### Reglas del `id`

- Formato: `^[a-z0-9_-]+/[a-z0-9_-]+$`
- Para filtros oficiales: `gtk-ai/<cmd>` donde `<cmd>` coincide con el argv0 interceptado y el nombre del repositorio.
- Terceros usan su propio prefijo: `acme/ls`.

### Versión del módulo

**No va en el manifiesto.** Se resuelve del tag Git en el momento de instalar:

```bash
gtkai filter install github.com/gtk-ai/date@v0.10.0
```

### `gtkai-core-version`

```json
{
  "version": "0.10.0",
  "constraint": "min"
}
```

| `constraint` | Significado |
|--------------|-------------|
| `"min"` | El `gtkai` en ejecución debe ser `>= version` |
| `"exact"` | El `gtkai` en ejecución debe coincidir exactamente |

Ambos campos son obligatorios. Un valor de `constraint` desconocido aborta la instalación.

## Paso 2 — Implementar la lógica del filtro (`filter/`)

Expón al menos:

- **`Rewrite(args []string) ([]string, bool)`** — modifica los argumentos antes de ejecutar el comando. Devuelve `false` si no hay cambios.
- **`FilterOutput(args, output, exitCode) string`** — transforma la salida del comando.

Define también la constante **`ID`** con el mismo valor que `id` en `gtkai.json`.

Consulta `filter/date.go` en este repo como referencia mínima.

## Paso 3 — Implementar el binario (`cmd/main.go`)

El binario es un filtro autónomo que habla **subprocess/v1** por stdin/stdout en JSON. Es el único transporte soportado.

**Petición** (stdin):

```json
{
  "operation": "rewrite",
  "args": ["-u"],
  "output": "",
  "exit_code": 0
}
```

Operaciones válidas: `"rewrite"` | `"filter_output"`.

**Respuesta** (stdout):

```json
{
  "args": ["-u", "+%Y-%m-%dT%H:%M:%SZ"],
  "changed": true,
  "output": "..."
}
```

El binario debe responder a un `rewrite` vacío en **menos de 500 ms** (comprobación de vida en la instalación).

## Paso 4 — Configurar el módulo Go (`go.mod`)

```go
module github.com/gtk-ai/<cmd>
```

Sustituye `<cmd>` por el nombre del repositorio (argv0 o nombre corto del filtro). En este template: `github.com/gtk-ai/date`.

## Paso 5 — Publicar e instalar

1. Crea un tag semver: `git tag v0.1.0 && git push origin v0.1.0`
2. (Opcional) Publica binarios en GitHub Releases.
3. Instala:

```bash
gtkai filter install github.com/gtk-ai/<cmd>@v0.1.0
```

Para reemplazar un filtro activo en el mismo argv0:

```bash
gtkai filter install github.com/gtk-ai/<cmd>@v0.1.0 --replace
```

Desinstalar:

```bash
gtkai filter uninstall gtk-ai/<cmd>
```

## Qué valida el core al instalar

1. El tag/ref resuelve a un semver válido (versión del módulo).
2. `gtkai.json` existe y parsea sin error.
3. Formato de `id` correcto.
4. `contract` es `subprocess/v1`.
5. Compatibilidad de plataforma.
6. `gtkai-core-version.version` es semver válido.
7. `gtkai-core-version.constraint` es `"min"` o `"exact"`.
8. La versión del `gtkai` en ejecución cumple la restricción.
9. Comprobación de vida: el binario responde en < 500 ms.

## Dónde queda instalado

```
~/.gtk-ai/filters/gtk-ai/date/
    date            # binario (nombre derivado del id)
    gtkai.json      # copia del manifiesto
```

El registro de filtros instalados vive en `~/.gtk-ai/filters.db`.

## Probar localmente

```bash
go test ./...
go build -o date ./cmd/
echo '{"operation":"rewrite","args":[],"output":"","exit_code":0}' | ./date
```

## Checklist al crear un módulo nuevo

- [ ] Repo bajo `gtk-ai/<cmd>`
- [ ] `go.mod` con `github.com/gtk-ai/<cmd>`
- [ ] `gtkai.json` con `id` = `gtk-ai/<cmd>`, `filters`, `platforms`, `contract`, `gtkai-core-version`
- [ ] `filter.ID` == `gtkai.json` → `id`
- [ ] `cmd/main.go` implementa subprocess/v1
- [ ] Tests de lógica y de protocolo pasan
- [ ] Tag semver publicado

## Referencias

- Core gtk-ai: [gtk-ai/gtk-ai](https://github.com/gtk-ai/gtk-ai)
- Este template: [gtk-ai/date](https://github.com/gtk-ai/date)
