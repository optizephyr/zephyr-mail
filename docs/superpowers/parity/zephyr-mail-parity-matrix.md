# Zephyr Mail Parity Matrix

| Scenario | Baseline | Target | Comparison | Schema/field check | Status |
| --- | --- | --- | --- | --- | --- |
| unknown-command | `node scripts/imap.js unknown-command` | `zephyr-mail unknown-command` | stderr + exit code | `Unknown command` message and exit `1` | PASS |
| fetch missing UID | `node scripts/imap.js fetch` | `zephyr-mail fetch` | stderr + exit code | required UID error text and exit `1` | PASS |
| send missing `--to` | `node scripts/smtp.js send --subject x` | `zephyr-mail send --subject x` | stderr + exit code | required recipient error text and exit `1` | PASS |
| list-mailboxes success | `node scripts/imap.js list-mailboxes` | `zephyr-mail list-mailboxes` | normalized JSON | array of mailbox objects with `name`, `delimiter`, `attributes`, `specialUse` | PASS |
| send success | `node scripts/smtp.js send --to recipient@example.com --subject "Parity subject" --body-file tests/parity/testdata/send-body.txt` | `zephyr-mail send --to recipient@example.com --subject "Parity subject" --body-file tests/parity/testdata/send-body.txt` | normalized JSON | object with `success`, `messageId`, `response`, `to` | PASS |

## Notes

- JS baseline is the real `node scripts/imap.js` / `node scripts/smtp.js` executables.
- Go target is the built `zephyr-mail` binary produced by `go build -o <tmp>/zephyr-mail ./cmd/zephyr-mail` inside the test harness.
- Success cases are validated on parsed JSON shape after comparing the raw process exit status and stderr path.
- `tests/parity/testdata/` contains the fixtures used by the SMTP/IMAP parity servers.
