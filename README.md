# Zephyr Mail

`zephyr-mail` 是一个用于 IMAP 和 SMTP 的命令行邮件工具。

## 快速开始

```bash
go build -o zephyr-mail ./cmd/zephyr-mail
```

在项目根目录放置 `.env`，然后执行：

```bash
./zephyr-mail check --limit 10
```

## 常用命令

```bash
./zephyr-mail check --limit 10
./zephyr-mail fetch <uid>
./zephyr-mail download <uid> --dir .
./zephyr-mail search --unseen --recent 24h
./zephyr-mail mark-read <uid>
./zephyr-mail mark-unread <uid>
./zephyr-mail list-mailboxes
./zephyr-mail send --to recipient@example.com --subject "Hello" --body "Message body"
./zephyr-mail test
```

## 文档

- [使用指南](docs/使用指南.md)
- [配置说明](docs/配置说明.md)
- [命令参考](docs/命令参考.md)
- [开发说明](docs/开发说明.md)
- [架构说明](docs/架构说明.md)

## 验证

```bash
go test ./...
```

## 支持的服务

- Gmail
- Outlook
- 163.com
- 126.com
- 188.com
- 任何标准 IMAP/SMTP 服务
