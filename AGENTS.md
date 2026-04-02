# AGENTS.md - Agent Guidelines for zephyr-mail

## Project Overview

This is a Go CLI tool for sending/receiving email via IMAP and SMTP protocols. It supports Gmail, Outlook, 163.com, and other standard IMAP/SMTP servers.

## Build & Run Commands

### Build
```bash
go build -o zephyr-mail ./cmd/zephyr-mail
```

### Run IMAP Commands
```bash
zephyr-mail check [--limit N] [--recent Nh|Nm] [--unseen]
zephyr-mail fetch <uid>
zephyr-mail search [--from X] [--subject X] [--unseen] [--recent Nh]
zephyr-mail download <uid> [--dir <path>] [--file <filename>]
zephyr-mail mark-read <uid> [uid2...]
zephyr-mail mark-unread <uid> [uid2...]
zephyr-mail list-mailboxes
```

### Run SMTP Commands
```bash
zephyr-mail send --to <email> --subject <text> [--body <text>] [--html] [--cc <email>] [--attach <file>]
zephyr-mail test
```

### Test Commands
- Run full test suite: `go test ./...`

## Code Style Guidelines

### Language
- Go

### Imports
- Use standard Go imports and `github.com/joho/godotenv` for `.env` loading

### Async/Await Patterns
- Use `defer` to ensure connections are closed

### Error Handling
- Return descriptive errors with context: `Missing IMAP_USER or IMAP_PASS`
- Use stderr for CLI errors and exit code 1 on failure

### Output Format
- Output pretty JSON to stdout
- Use stderr for status messages and errors

### Naming Conventions
- **Functions**: camelCase (e.g., `createImapConfig`, `checkEmails`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `DEFAULT_MAILBOX`, `IMAP_ID`)
- **Variables**: camelCase with descriptive names

### File Structure
- Main entry: `cmd/zephyr-mail`
- Configuration via environment variables in `.env`
- Common utility functions live under `internal/`

### Argument Parsing
- Cobra handles command parsing and flags

### Key Libraries Used
- `cobra` - CLI framework
- `go-imap` - IMAP client
- `gomail` or equivalent - SMTP email sending
- `godotenv` - Environment variable loading

### Security
- Never commit credentials; use `.env` file
- Handle TLS rejection based on environment variables
- Validate required configuration before operations

### Configuration
Environment variables in `.env`:
- IMAP: `IMAP_HOST`, `IMAP_PORT`, `IMAP_USER`, `IMAP_PASS`, `IMAP_TLS`, `IMAP_MAILBOX`
- SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_SECURE`, `SMTP_FROM`
