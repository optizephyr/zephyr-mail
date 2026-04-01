# AGENTS.md - Agent Guidelines for zephyr-mail

## Project Overview

This is a Node.js CLI tool for sending/receiving email via IMAP and SMTP protocols. It supports Gmail, Outlook, 163.com, and other standard IMAP/SMTP servers.

## Build & Run Commands

### Install Dependencies
```bash
npm install
```

### Run IMAP Commands
```bash
node scripts/imap.js check [--limit N] [--recent Nh|Nm] [--unseen]
node scripts/imap.js fetch <uid>
node scripts/imap.js search [--from X] [--subject X] [--unseen] [--recent Nh]
node scripts/imap.js mark-read <uid> [uid2...]
node scripts/imap.js mark-unread <uid> [uid2...]
node scripts/imap.js list-mailboxes
```

### Run SMTP Commands
```bash
node scripts/smtp.js send --to <email> --subject <text> [--body <text>] [--html] [--cc <email>] [--attach <file>]
node scripts/smtp.js test
```

### NPM Scripts (Aliases)
```bash
npm run check      # node scripts/imap.js check
npm run fetch      # node scripts/imap.js fetch
npm run search    # node scripts/imap.js search
```

### Test Commands
- No project-specific tests exist
- Run manual tests: `node scripts/imap.js check` and `node scripts/smtp.js test`

## Code Style Guidelines

### Language
- JavaScript (Node.js) using CommonJS
- No TypeScript in this project

### Imports
- Use `require()` for dependencies
- Load dotenv early: `require('dotenv').config({ path: path.resolve(__dirname, '../.env') })`
- Use destructuring for modules: `const { ImapFlow } = require('imapflow');`

### Async/Await Patterns
- Use `async/await` for all asynchronous operations
- Always use `try/finally` to ensure connections are closed:
  ```javascript
  async function example() {
    const client = await connect();
    try {
      // operations
    } finally {
      await client.logout();
    }
  }
  ```

### Error Handling
- Throw descriptive errors with context: `throw new Error('Missing IMAP_USER or IMAP_PASS')`
- Use `console.error('Error:', err.message)` for errors in CLI
- Exit with code 1 on errors: `process.exit(1)`

### Output Format
- Output JSON to stdout: `console.log(JSON.stringify(result, null, 2))`
- Use `console.error` for status messages and errors

### Naming Conventions
- **Functions**: camelCase (e.g., `createImapConfig`, `checkEmails`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `DEFAULT_MAILBOX`, `IMAP_ID`)
- **Variables**: camelCase with descriptive names

### File Structure
- Main entry: `scripts/imap.js`, `scripts/smtp.js`
- Configuration via environment variables in `.env`
- Common utility functions defined in same file

### Argument Parsing
- Manual parsing of CLI arguments (no external library)
- Support `--option value` format and boolean flags
- Track positional arguments separately

### Key Libraries Used
- `imapflow` - Modern IMAP client with UTF-8 support
- `nodemailer` - SMTP email sending
- `mailparser` - Email parsing
- `dotenv` - Environment variable loading

### Security
- Never commit credentials; use `.env` file
- Handle TLS rejection based on environment variables
- Validate required configuration before operations

### Configuration
Environment variables in `.env`:
- IMAP: `IMAP_HOST`, `IMAP_PORT`, `IMAP_USER`, `IMAP_PASS`, `IMAP_TLS`, `IMAP_MAILBOX`
- SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_SECURE`, `SMTP_FROM`
