# zephyr-mail Go Migration Design

Date: 2026-04-01
Status: Draft (pending compatibility review)
Owner: CLI migration

## 1. Goal

Migrate the current JavaScript email CLI project to Go and expose a single executable command named `zephyr-mail`.

Primary requirement: full compatibility with existing command behavior and arguments.

## 2. Scope and Compatibility Boundaries

### In scope

- Replace JS runtime implementation with Go.
- Keep command names and argument semantics compatible, using `scripts/imap.js` and `scripts/smtp.js` as the source of truth when README and scripts differ.
- Keep configuration keys in `.env` unchanged.
- Keep stdout/stderr/exit behavior compatible.

### Out of scope

- Adding unrelated new features during migration.
- Changing provider-side capabilities (IMAP/SMTP server behavior).

## 3. Command Surface (Cobra)

Executable:

- `zephyr-mail <command> [options]`

IMAP commands:

- `check`
  - `--mailbox <name>`
  - `--limit <n>` (default 10)
  - `--recent <Nm|Nh|Nd>` (e.g., `30m`, `2h`, `7d`, `24h`)
  - `--unseen <true|false>` (source-compatible behavior: bare `--unseen` is not treated as true for `check`)
- `fetch <uid>`
  - `--mailbox <name>`
- `download <uid>`
  - `--mailbox <name>`
  - `--dir <path>` (default current directory)
  - `--file <filename>`
- `search`
  - `--mailbox <name>`
  - `--limit <n>` (default 100)
  - `--unseen`
  - `--seen`
  - `--flagged`
  - `--answered`
  - `--from <addr>`
  - `--to <addr>`
  - `--subject <text>`
  - `--recent <Nm|Nh|Nd>` (e.g., `30m`, `2h`, `7d`, `24h`)
  - `--since <date>`
  - `--before <date>`
  - `--uid <n|range>`
- `mark-read <uid> [uid2 ...]`
  - `--mailbox <name>`
- `mark-unread <uid> [uid2 ...]`
  - `--mailbox <name>`
- `list-mailboxes`

SMTP commands:

- `send`
  - `--to <email[,email...]>` (required)
  - `--subject <text>` (required unless `--subject-file`)
  - `--subject-file <file>`
  - `--body <text>`
  - `--body-file <file>`
  - `--html`
  - `--html-file <file>`
  - `--cc <email[,email...]>`
  - `--bcc <email[,email...]>`
  - `--from <email>`
  - `--attach <file1,file2,...>`
- `test`

## 4. Architecture

Planned structure:

```text
cmd/zephyr-mail/
  main.go
  root.go
internal/cli/
  check.go
  fetch.go
  download.go
  search.go
  mark_read.go
  mark_unread.go
  list_mailboxes.go
  send.go
  test.go
internal/config/
  env.go
  validate.go
internal/output/
  json.go
internal/imap/
  client.go
  models.go
  query.go
  service.go
  timeparse.go
internal/smtp/
  client.go
  models.go
  service.go
internal/common/
  errors.go
  files.go
```

Responsibilities:

- `cmd` and `internal/cli`: command wiring, flags, argument validation, service dispatch.
- `internal/config`: env loading and typed configuration validation.
- `internal/imap`: IMAP connectivity, query building, message/attachment handling, mailbox operations.
- `internal/smtp`: SMTP connection test and send operations.
- `internal/output`: unified JSON success output and error output behavior.
- `internal/common`: shared file helpers and error wrapping.

## 5. Data Flow

1. Root command initializes environment and registers subcommands.
2. Subcommand parses flags and validates required inputs.
3. Subcommand calls IMAP or SMTP service.
4. Service returns typed result or wrapped error.
5. CLI layer preserves source command stdout behavior on success: commands that currently emit diagnostics before payload (for example IMAP `check`/`search`) must keep that pre-JSON output; other commands print JSON payload only.
6. On error, print message to stderr and exit with code 1.

## 6. Error Handling Strategy

Error classes:

- `ConfigError`: missing or invalid env config.
- `ValidationError`: invalid command arguments.
- `ProviderError`: IMAP/SMTP/network failures.
- `RuntimeError`: local runtime issues (I/O, parsing).

Rules:

- Keep user-facing errors concise and actionable.
- Preserve existing behavior: `Error: <message>` style on stderr.
- Ensure resource cleanup with `defer` (connections, mailbox locks).

## 7. Compatibility Details

- Keep `.env` keys unchanged for IMAP/SMTP.
- Preserve JSON schema and field names exactly (including key casing, presence/absence rules, and nullability).
- Preserve relative-time parsing format (`<Nm|Nh|Nd>`, e.g., `30m`, `2h`, `7d`, `24h`).
- Preserve search criteria composition (AND semantics across filters).
- Preserve command-specific stdout/stderr side effects from `scripts/imap.js`, including current diagnostic output lines emitted by `check`/`search` (for example search-criteria/found logs), unless an explicit compatibility waiver is approved.
- Preserve boolean flag semantics exactly:
  - `check --unseen` is only treated as true when value is literal `"true"` (bare `--unseen` remains false);
  - `search --unseen/--seen/--flagged/--answered` are presence booleans.
- Preserve attachment behavior:
  - default download directory is current directory;
  - when `--file` is not found, return available attachment names.
- Preserve SMTP `send` precedence exactly (source-compatible):
  - `subject`: if `--subject-file` is present, file content overrides `--subject`.
  - body resolution order:
    1) if `--body-file` is present: read file; if file name ends with `.html` OR `--html` is present, assign to HTML body; otherwise assign to text body;
    2) else if `--html-file` is present: assign file content to HTML body;
    3) else if `--body` is present: assign to text body;
    4) final send fallback: if neither text nor HTML exists, send empty text body.

## 8. Migration Strategy (Recommended)

Use phased command-level migration for safety and steady validation:

1. Scaffold Go module and Cobra command tree.
2. Implement shared config and output layers.
3. Implement IMAP command set incrementally.
4. Implement SMTP command set.
5. Perform parity verification against current JS command outputs/behaviors.
6. Update docs to `zephyr-mail` usage and retire JS entrypoints only after all Section 9 parity checks pass and are recorded in a parity matrix (command x scenario x stdout/stderr/exit code).

## 9. Verification and Acceptance Criteria

Acceptance requirements:

- All compatibility commands available via `zephyr-mail`: `check`, `fetch`, `download`, `search`, `mark-read`, `mark-unread`, `list-mailboxes`, `send`, `test`.
- Existing `.env` works without key changes.
- Output channels and exit codes match compatibility requirements.
- For each command, success and failure paths must match JS behavior for stdout payload shape, stderr prefix (`Error: ...`), and process exit code.
- Cobra argument/flag parse errors must be normalized to existing JS-visible behavior (message format and exit code).
- Unknown-command and missing-argument paths must match source behavior (including stderr wording and exit code 1) for both IMAP and SMTP entrypoints.
- Manual parity checks pass for representative IMAP/SMTP scenarios.

Verification baseline (automated parity tests + manual smoke):

- Automated parity tests compare Go vs JS raw stdout, raw stderr, and exit code per command/scenario, then additionally validate parsed JSON payload shape where applicable.
- Manual smoke remains required for live provider validation.

- `zephyr-mail check --limit 10`
- `zephyr-mail search --unseen --recent 24h`
- `zephyr-mail fetch <uid>`
- `zephyr-mail download <uid> --dir .`
- `zephyr-mail mark-read <uid>` and `zephyr-mail mark-unread <uid>`
- `zephyr-mail list-mailboxes`
- `zephyr-mail send --to ... --subject ... --body ...`
- `zephyr-mail test`
- `zephyr-mail unknown-command` (error parity)
- `zephyr-mail fetch` (missing positional argument parity)

## 10. Dependencies

Planned dependencies:

- `github.com/spf13/cobra`
- `github.com/joho/godotenv`
- IMAP and MIME libraries compatible with UTF-8 header/body parsing
- SMTP implementation library with attachment support

Final library selection is constrained by parity outcomes during implementation.
