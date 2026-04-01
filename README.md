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

## Supported Providers

- Gmail
- Outlook
- 163.com
- 126.com
- 188.com
- Any standard IMAP/SMTP server

## License

MIT
