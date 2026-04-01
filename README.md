# Zephyr Mail

`zephyr-mail` is the primary CLI for IMAP and SMTP email operations.

## Install

```bash
go install github.com/netease/zephyr-mail/cmd/zephyr-mail@latest
```

Or build locally:

```bash
go build -o zephyr-mail ./cmd/zephyr-mail
```

## Usage

```bash
zephyr-mail check --limit 10
zephyr-mail fetch <uid>
zephyr-mail download <uid> --dir .
zephyr-mail search --unseen --recent 24h
zephyr-mail mark-read <uid>
zephyr-mail mark-unread <uid>
zephyr-mail list-mailboxes
zephyr-mail send --to recipient@example.com --subject "Hello" --body "Message body"
zephyr-mail test
```

## Parity Gate

The JS entrypoints in `scripts/` are kept as compatibility baselines while parity is tracked.

Run the release gate with:

```bash
go test ./tests/parity -run TestReleaseGateMatrixComplete -v
```

Run the full verification suite with:

```bash
go test ./...
```

The current parity matrix lives in `docs/superpowers/parity/zephyr-mail-parity-matrix.md`.

## Supported Providers

- Gmail
- Outlook
- 163.com
- 126.com
- 188.com
- Any standard IMAP/SMTP server

## License

MIT
