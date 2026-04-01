# Zephyr Mail Parity Matrix

| Scenario | Baseline | Target | Comparison | Schema/field check | Status |
| --- | --- | --- | --- | --- | --- |
| unknown-command | `node scripts/imap.js unknown-command` | `zephyr-mail unknown-command` | stderr + exit code | `Unknown command` message and exit `1` | PASS |
| check success | `node scripts/imap.js check --limit 10` | `zephyr-mail check --limit 10` | normalized JSON | array of message objects with IMAP fields | PASS |
| fetch missing UID | `node scripts/imap.js fetch` | `zephyr-mail fetch` | stderr + exit code | required UID error text and exit `1` | PASS |
| fetch success | `node scripts/imap.js fetch <uid>` | `zephyr-mail fetch <uid>` | normalized JSON | array of parsed message objects | PASS |
| download success | `node scripts/imap.js download <uid> --dir .` | `zephyr-mail download <uid> --dir .` | normalized JSON | object with `uid`, `downloaded`, `message` | PASS |
| search unseen recent | `node scripts/imap.js search --unseen --recent 24h` | `zephyr-mail search --unseen --recent 24h` | normalized JSON | array of parsed message objects | PASS |
| mark-read success | `node scripts/imap.js mark-read <uid>` | `zephyr-mail mark-read <uid>` | normalized JSON | object with `success`, `uids`, `action`, `count` | PASS |
| mark-unread success | `node scripts/imap.js mark-unread <uid>` | `zephyr-mail mark-unread <uid>` | normalized JSON | object with `success`, `uids`, `action`, `count` | PASS |
| list-mailboxes success | `node scripts/imap.js list-mailboxes` | `zephyr-mail list-mailboxes` | normalized JSON | array of mailbox objects with `name`, `delimiter`, `attributes`, `specialUse` | PASS |
| send missing `--to` | `node scripts/smtp.js send --subject x` | `zephyr-mail send --subject x` | stderr + exit code | required recipient error text and exit `1` | PASS |
| send success | `node scripts/smtp.js send --to recipient@example.com --subject "Parity subject" --body-file tests/parity/testdata/send-body.txt` | `zephyr-mail send --to recipient@example.com --subject "Parity subject" --body-file tests/parity/testdata/send-body.txt` | normalized JSON | object with `success`, `messageId`, `response`, `to` | PASS |
| test success | `node scripts/smtp.js test` | `zephyr-mail test` | normalized JSON | object with `success`, `message`, `messageId` | PASS |

## Notes

- JS baseline is the real `node scripts/imap.js` / `node scripts/smtp.js` executables.
- Go target is the built `zephyr-mail` binary produced by `go build -o <tmp>/zephyr-mail ./cmd/zephyr-mail` inside the test harness.
- Success cases are validated on parsed JSON shape after comparing the raw process exit status and stderr path.
- `tests/parity/testdata/` contains the fixtures used by the SMTP/IMAP parity servers.
- The release gate is `go test ./tests/parity -run TestReleaseGateMatrixComplete -v`; it passes only when every required scenario above is marked `PASS`.
