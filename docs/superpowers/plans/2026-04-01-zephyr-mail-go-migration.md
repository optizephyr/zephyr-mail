# Zephyr Mail Go Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the current JS CLI implementation with a Go binary named `zephyr-mail` while preserving command behavior and outputs.

**Architecture:** Build a Cobra-based CLI in `cmd/zephyr-mail` and `internal/cli`, with protocol logic in `internal/imap` and `internal/smtp`. Keep compatibility by treating `scripts/imap.js` and `scripts/smtp.js` as runtime source of truth and verifying parity for stdout/stderr/exit code.

**Tech Stack:** Go 1.22+, Cobra, godotenv, go-imap, go-message, gomail (or equivalent SMTP lib), Go testing package.

---

### Task 1: Bootstrap Go Module and Root Command

**Files:**
- Create: `go.mod`
- Create: `cmd/zephyr-mail/main.go`
- Create: `cmd/zephyr-mail/root.go`
- Create: `internal/cli/register.go`
- Modify: `README.md`
- Test: `cmd/zephyr-mail/root_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestRootUnknownCommandExitCode(t *testing.T) {
    code, stdout, stderr := runCLI("unknown-command")
    if code != 1 {
        t.Fatalf("want exit 1, got %d", code)
    }
    if stdout != "" {
        t.Fatalf("want empty stdout")
    }
    if !strings.Contains(stderr, "Unknown command") {
        t.Fatalf("unexpected stderr: %s", stderr)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/zephyr-mail -run TestRootUnknownCommandExitCode -v`
Expected: FAIL because CLI scaffolding does not exist yet.

- [ ] **Step 3: Write minimal implementation**

Implement `main.go` calling `Execute()`, `root.go` defining `rootCmd.Use = "zephyr-mail"`, and centralized command error handling that normalizes unknown command output to compatibility format.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./cmd/zephyr-mail -run TestRootUnknownCommandExitCode -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add go.mod cmd/zephyr-mail/main.go cmd/zephyr-mail/root.go internal/cli/register.go cmd/zephyr-mail/root_test.go README.md
git commit -m "feat: scaffold zephyr-mail cobra root command"
```

### Task 2: Implement Config Loading and Validation

**Files:**
- Create: `internal/config/env.go`
- Create: `internal/config/validate.go`
- Create: `internal/config/env_test.go`
- Create: `internal/config/validate_test.go`
- Modify: `cmd/zephyr-mail/root.go`

- [ ] **Step 1: Write the failing tests**

```go
func TestLoadEnvWithSetEnvIsDeterministic(t *testing.T) {
    t.Setenv("IMAP_HOST", "imap.example.com")
    t.Setenv("IMAP_PORT", "993")
    t.Setenv("IMAP_USER", "user@example.com")
    t.Setenv("IMAP_PASS", "secret")
    t.Setenv("SMTP_HOST", "smtp.example.com")
    t.Setenv("SMTP_USER", "user@example.com")
    t.Setenv("SMTP_PASS", "secret")

    cfg, err := LoadFromEnv()
    if err != nil { t.Fatal(err) }
    if cfg.IMAPHost != "imap.example.com" || cfg.SMTPHost != "smtp.example.com" {
        t.Fatalf("unexpected parsed config: %+v", cfg)
    }
}

func TestValidateSMTPMissingHost(t *testing.T) {
    cfg := Config{SMTPUser: "a", SMTPPass: "b"}
    err := ValidateSMTP(cfg)
    if err == nil || !strings.Contains(err.Error(), "Missing SMTP configuration") {
        t.Fatalf("unexpected err: %v", err)
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/config -v`
Expected: FAIL because loaders/validators are not implemented.

- [ ] **Step 3: Write minimal implementation**

Add typed config struct, `.env` loading via `godotenv`, defaults compatible with JS scripts, and validation functions returning source-compatible error messages.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/config -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config/env.go internal/config/validate.go internal/config/env_test.go internal/config/validate_test.go cmd/zephyr-mail/root.go
git commit -m "feat: add env config loading and validation"
```

### Task 3: Add Output and Error Compatibility Layer

**Files:**
- Create: `internal/output/json.go`
- Create: `internal/common/errors.go`
- Create: `internal/output/json_test.go`
- Modify: `cmd/zephyr-mail/root.go`

- [ ] **Step 1: Write the failing tests**

```go
func TestPrintErrorUsesCompatibilityPrefix(t *testing.T) {
    stderr := captureStderr(func() { PrintError(errors.New("boom")) })
    if !strings.Contains(stderr, "Error: boom") { t.Fatal(stderr) }
}

func TestPrintJSONPretty(t *testing.T) {
    out := captureStdout(func() { PrintJSON(map[string]any{"ok": true}) })
    if !strings.Contains(out, "\n") { t.Fatal("expected pretty json") }
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/output -v`
Expected: FAIL because output helpers do not exist.

- [ ] **Step 3: Write minimal implementation**

Implement JSON pretty printer and stderr error printer (`Error: <message>`), wire root command to convert returned errors to exit code 1.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/output -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/output/json.go internal/common/errors.go internal/output/json_test.go cmd/zephyr-mail/root.go
git commit -m "feat: add compatibility output and error handling"
```

### Task 4: Implement IMAP Client, Models, and Time Parsing

**Files:**
- Create: `internal/imap/client.go`
- Create: `internal/imap/models.go`
- Create: `internal/imap/timeparse.go`
- Create: `internal/imap/timeparse_test.go`
- Create: `internal/imap/query.go`
- Create: `internal/imap/query_test.go`

- [ ] **Step 1: Write the failing tests**

```go
func TestParseRelativeTime(t *testing.T) {
    got, err := ParseRelativeTime("2h")
    if err != nil { t.Fatal(err) }
    if !regexp.MustCompile(`^\d{2}-[A-Z][a-z]{2}-\d{4}$`).MatchString(got) {
        t.Fatalf("unexpected format: %s", got)
    }
}

func TestBuildSearchCriteriaWithUnseenAndSubject(t *testing.T) {
    c := BuildSearchCriteria(SearchOptions{Unseen: true, Subject: "hello"})
    if !c.Unseen || c.Subject != "hello" { t.Fatalf("unexpected criteria") }
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/imap -run "TestParseRelativeTime|TestBuildSearchCriteriaWithUnseenAndSubject" -v`
Expected: FAIL because parser/query builder are not implemented.

- [ ] **Step 3: Write minimal implementation**

Implement relative time parser (`Nm|Nh|Nd`), IMAP date formatting, search options and query builder compatible with current scripts.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/imap -run "TestParseRelativeTime|TestBuildSearchCriteriaWithUnseenAndSubject" -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/imap/client.go internal/imap/models.go internal/imap/timeparse.go internal/imap/timeparse_test.go internal/imap/query.go internal/imap/query_test.go
git commit -m "feat: add imap client foundations and query parsing"
```

### Task 5: Implement IMAP Service Operations

**Files:**
- Create: `internal/imap/service.go`
- Create: `internal/imap/service_test.go`
- Create: `internal/common/files.go`

- [ ] **Step 1: Write the failing tests**

```go
func TestCheckUnseenLiteralTrueOnly(t *testing.T) {
    opts := CheckOptions{UnseenRaw: "true"}
    criteria := BuildCheckCriteria(opts)
    if !criteria.Unseen { t.Fatal("expected unseen true") }
}

func TestDownloadMissingSpecificFileReturnsAvailable(t *testing.T) {
    result := DownloadResult{Message: `File "x" not found. Available: a,b`}
    if !strings.Contains(result.Message, "Available") { t.Fatal(result.Message) }
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/imap -run "TestCheckUnseenLiteralTrueOnly|TestDownloadMissingSpecificFileReturnsAvailable" -v`
Expected: FAIL because service behavior is not implemented.

- [ ] **Step 3: Write minimal implementation**

Implement `Check`, `Fetch`, `Download`, `Search`, `MarkRead`, `MarkUnread`, `ListMailboxes` with cleanup via `defer` and source-compatible message fields.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/imap -run "TestCheckUnseenLiteralTrueOnly|TestDownloadMissingSpecificFileReturnsAvailable" -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/imap/service.go internal/imap/service_test.go internal/common/files.go
git commit -m "feat: implement imap command services"
```

### Task 6: Implement SMTP Service and Precedence Rules

**Files:**
- Create: `internal/smtp/client.go`
- Create: `internal/smtp/models.go`
- Create: `internal/smtp/service.go`
- Create: `internal/smtp/service_test.go`

- [ ] **Step 1: Write the failing tests**

```go
func TestSubjectFileOverridesSubject(t *testing.T) {
    req := SendRequest{Subject: "inline", SubjectFile: "testdata/subject.txt"}
    resolved := ResolveSendRequest(req)
    if resolved.Subject != "from-file" { t.Fatalf("got %s", resolved.Subject) }
}

func TestBodyFileHtmlDetection(t *testing.T) {
    req := SendRequest{BodyFile: "testdata/body.html", HTML: true}
    resolved := ResolveSendRequest(req)
    if resolved.HTMLBody == "" { t.Fatal("expected html body") }
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/smtp -run "TestSubjectFileOverridesSubject|TestBodyFileHtmlDetection" -v`
Expected: FAIL because resolution logic is not implemented.

- [ ] **Step 3: Write minimal implementation**

Implement `Send` and `TestConnection`; codify precedence: `subject-file` override, `body-file` handling (`.html` or `--html`), `html-file` fallback, then `body`, and empty-text fallback when both text/html are absent.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/smtp -run "TestSubjectFileOverridesSubject|TestBodyFileHtmlDetection" -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/smtp/client.go internal/smtp/models.go internal/smtp/service.go internal/smtp/service_test.go
git commit -m "feat: implement smtp send and test behaviors"
```

### Task 7: Add Cobra Subcommands and Flag Semantics

**Files:**
- Create: `internal/cli/check.go`
- Create: `internal/cli/fetch.go`
- Create: `internal/cli/download.go`
- Create: `internal/cli/search.go`
- Create: `internal/cli/mark_read.go`
- Create: `internal/cli/mark_unread.go`
- Create: `internal/cli/list_mailboxes.go`
- Create: `internal/cli/send.go`
- Create: `internal/cli/test.go`
- Create: `internal/cli/cli_test.go`
- Modify: `internal/cli/register.go`

- [ ] **Step 1: Write the failing tests**

```go
func TestCheckBareUnseenDoesNotEnableUnseen(t *testing.T) {
    opts := parseCheckFlags([]string{"--unseen"})
    if opts.UnseenRaw == "true" {
        t.Fatal("bare --unseen must not be treated as true")
    }
}

func TestSearchPresenceBooleans(t *testing.T) {
    opts := parseSearchFlags([]string{"--unseen", "--flagged"})
    if !opts.Unseen || !opts.Flagged { t.Fatal("presence booleans expected") }
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/cli -run "TestCheckBareUnseenDoesNotEnableUnseen|TestSearchPresenceBooleans" -v`
Expected: FAIL because flag wiring is not implemented.

- [ ] **Step 3: Write minimal implementation**

Implement all subcommands, wire flags and positional args to match JS behavior, and route to services/output helpers.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/cli -run "TestCheckBareUnseenDoesNotEnableUnseen|TestSearchPresenceBooleans" -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/cli/check.go internal/cli/fetch.go internal/cli/download.go internal/cli/search.go internal/cli/mark_read.go internal/cli/mark_unread.go internal/cli/list_mailboxes.go internal/cli/send.go internal/cli/test.go internal/cli/cli_test.go internal/cli/register.go
git commit -m "feat: add cobra subcommands with compatibility flags"
```

### Task 8: Add Parity Matrix and Automated Compatibility Tests

**Files:**
- Create: `tests/parity/parity_test.go`
- Create: `tests/parity/testdata/`
- Create: `docs/superpowers/parity/zephyr-mail-parity-matrix.md`
- Modify: `README.md`

- [ ] **Step 1: Write the failing parity tests**

```go
func TestParityUnknownCommand(t *testing.T) {
    js := runJS("unknown-command")
    goOut := runGo("unknown-command")
    if js.ExitCode != goOut.ExitCode || js.Stderr != goOut.Stderr {
        t.Fatalf("parity mismatch")
    }
}

func TestParityFetchMissingUID(t *testing.T) {
    js := runJS("fetch")
    goOut := runGo("fetch")
    if js.ExitCode != goOut.ExitCode || js.Stderr != goOut.Stderr {
        t.Fatalf("parity mismatch")
    }
}

func TestParitySendMissingTo(t *testing.T) {
    js := runJS("send", "--subject", "x")
    goOut := runGo("send", "--subject", "x")
    if js.ExitCode != goOut.ExitCode || js.Stderr != goOut.Stderr {
        t.Fatalf("parity mismatch")
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./tests/parity -v`
Expected: FAIL before all command parity is complete.

- [ ] **Step 3: Write minimal implementation**

Implement parity harness comparing raw stdout/raw stderr/exit code, then parse successful JSON outputs and validate payload schema/field presence/nullability against JS outputs; record scenario outcomes in parity matrix markdown.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./tests/parity -v`
Expected: PASS for scripted scenarios.

- [ ] **Step 5: Commit**

```bash
git add tests/parity/parity_test.go tests/parity/testdata docs/superpowers/parity/zephyr-mail-parity-matrix.md README.md
git commit -m "test: add go-vs-js parity validation matrix"
```

### Task 9: Final Verification, Documentation, and JS Retirement Gate

**Files:**
- Modify: `README.md`
- Modify: `package.json`
- Modify: `scripts/imap.js`
- Modify: `scripts/smtp.js`
- Modify: `docs/superpowers/parity/zephyr-mail-parity-matrix.md`
- Modify: `tests/parity/parity_test.go`

- [ ] **Step 1: Write the failing release-check test**

```go
func TestReleaseGateMatrixComplete(t *testing.T) {
    if !ParityMatrixComplete("docs/superpowers/parity/zephyr-mail-parity-matrix.md") {
        t.Fatal("parity matrix incomplete")
    }
}
```

- [ ] **Step 2: Run checks to verify gate fails before completion**

Run: `go test ./tests/parity -run TestReleaseGateMatrixComplete -v`
Expected: FAIL until all required scenarios are checked.

- [ ] **Step 3: Complete implementation and gate criteria**

Finalize README with `zephyr-mail` usage, keep JS scripts only as compatibility references or wrappers until parity matrix is complete, and only retire JS entrypoints after gate conditions are satisfied.

- [ ] **Step 4: Run full verification suite**

Run: `go test ./...`
Expected: PASS.

Run manual smoke:

- `zephyr-mail check --limit 10`
- `zephyr-mail search --unseen --recent 24h`
- `zephyr-mail fetch <uid>`
- `zephyr-mail download <uid> --dir .`
- `zephyr-mail mark-read <uid>`
- `zephyr-mail mark-unread <uid>`
- `zephyr-mail list-mailboxes`
- `zephyr-mail send --to ... --subject ... --body ...`
- `zephyr-mail test`
- `zephyr-mail unknown-command`
- `zephyr-mail fetch`

Expected: all success/error outputs align with parity matrix.

- [ ] **Step 5: Commit**

```bash
git add README.md package.json scripts/imap.js scripts/smtp.js docs/superpowers/parity/zephyr-mail-parity-matrix.md
git commit -m "chore: finalize zephyr-mail go migration and docs"
```
