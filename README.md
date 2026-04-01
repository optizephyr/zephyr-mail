# Zephyr Mail

A CLI tool for sending and receiving email via IMAP and SMTP protocols.

## Installation

```bash
go install github.com/netease/zephyr-mail/cmd/zephyr-mail@latest
```

Or build from source:

```bash
git clone https://github.com/netease/zephyr-mail.git
cd zephyr-mail
go build -o zephyr-mail ./cmd/zephyr-mail
```

## Usage

```bash
zephyr-mail [command] [flags]
```

## Compatibility Tests

Run the Go-vs-JS parity suite with:

```bash
go test ./tests/parity -v
```

The matrix for the scripted scenarios lives in `docs/superpowers/parity/zephyr-mail-parity-matrix.md`.

## Supported Providers

- Gmail
- Outlook
- 163.com
- 126.com
- 188.com
- Any standard IMAP/SMTP server

## License

MIT
