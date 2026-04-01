package parity

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

type commandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type mailboxEntry struct {
	Name       string   `json:"name"`
	Delimiter  string   `json:"delimiter"`
	Attributes []string `json:"attributes"`
	SpecialUse any      `json:"specialUse"`
}

type sendPayload struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId"`
	Response  string `json:"response"`
	To        any    `json:"to"`
}

var (
	repoRootOnce sync.Once
	repoRootPath string
	repoRootErr  error

	goBinaryOnce sync.Once
	goBinaryPath string
	goBinaryErr  error
)

func TestReleaseGateMatrixComplete(t *testing.T) {
	if !ParityMatrixComplete("docs/superpowers/parity/zephyr-mail-parity-matrix.md") {
		t.Fatal("parity matrix incomplete")
	}
}

func TestParityUnknownCommand(t *testing.T) {
	env := map[string]string{
		"IMAP_HOST":    "127.0.0.1",
		"IMAP_PORT":    "993",
		"IMAP_USER":    "parity-user",
		"IMAP_PASS":    "parity-pass",
		"IMAP_TLS":     "false",
		"IMAP_MAILBOX": "INBOX",
		"SMTP_HOST":    "127.0.0.1",
		"SMTP_PORT":    "587",
		"SMTP_USER":    "parity-user",
		"SMTP_PASS":    "parity-pass",
		"SMTP_SECURE":  "false",
		"SMTP_FROM":    "parity-user@example.com",
	}

	js := runJS(t, "imap.js", []string{"unknown-command"}, env)
	goOut := runGo(t, []string{"unknown-command"}, env)
	assertExactParity(t, js, goOut)
}

func TestParityFetchMissingUID(t *testing.T) {
	env := map[string]string{
		"IMAP_HOST":    "127.0.0.1",
		"IMAP_PORT":    "993",
		"IMAP_USER":    "parity-user",
		"IMAP_PASS":    "parity-pass",
		"IMAP_TLS":     "false",
		"IMAP_MAILBOX": "INBOX",
		"SMTP_HOST":    "127.0.0.1",
		"SMTP_PORT":    "587",
		"SMTP_USER":    "parity-user",
		"SMTP_PASS":    "parity-pass",
		"SMTP_SECURE":  "false",
		"SMTP_FROM":    "parity-user@example.com",
	}

	js := runJS(t, "imap.js", []string{"fetch"}, env)
	goOut := runGo(t, []string{"fetch"}, env)
	assertExactParity(t, js, goOut)
}

func TestParitySendMissingTo(t *testing.T) {
	env := map[string]string{
		"IMAP_HOST":    "127.0.0.1",
		"IMAP_PORT":    "993",
		"IMAP_USER":    "parity-user",
		"IMAP_PASS":    "parity-pass",
		"IMAP_TLS":     "false",
		"IMAP_MAILBOX": "INBOX",
		"SMTP_HOST":    "127.0.0.1",
		"SMTP_PORT":    "587",
		"SMTP_USER":    "parity-user",
		"SMTP_PASS":    "parity-pass",
		"SMTP_SECURE":  "false",
		"SMTP_FROM":    "parity-user@example.com",
	}

	js := runJS(t, "smtp.js", []string{"send", "--subject", "x"}, env)
	goOut := runGo(t, []string{"send", "--subject", "x"}, env)
	assertExactParity(t, js, goOut)
}

func TestParityListMailboxesSuccess(t *testing.T) {
	server := newIMAPParityServer(t)
	defer server.Close()

	env := map[string]string{
		"IMAP_HOST":    server.Host(),
		"IMAP_PORT":    server.Port(),
		"IMAP_USER":    "parity-user",
		"IMAP_PASS":    "parity-pass",
		"IMAP_TLS":     "false",
		"IMAP_MAILBOX": "INBOX",
		"SMTP_HOST":    "127.0.0.1",
		"SMTP_PORT":    "587",
		"SMTP_USER":    "parity-user",
		"SMTP_PASS":    "parity-pass",
		"SMTP_SECURE":  "false",
		"SMTP_FROM":    "parity-user@example.com",
	}

	js := runJS(t, "imap.js", []string{"list-mailboxes"}, env)
	goOut := runGo(t, []string{"list-mailboxes"}, env)
	assertJSONParity(t, js, goOut, func(t *testing.T, jsVal, goVal any) {
		want := []mailboxEntry{{
			Name:       "Reports",
			Delimiter:  "/",
			Attributes: []string{},
			SpecialUse: nil,
		}}
		if !reflect.DeepEqual(decodeMailboxList(t, jsVal), want) {
			t.Fatalf("unexpected JS mailbox payload: %#v", jsVal)
		}
		if !reflect.DeepEqual(decodeMailboxList(t, goVal), want) {
			t.Fatalf("unexpected Go mailbox payload: %#v", goVal)
		}
	})
}

func TestParitySendSuccess(t *testing.T) {
	server := newSMTPParityServer(t)
	defer server.Close()

	bodyPath := fixturePath(t, "send-body.txt")
	env := map[string]string{
		"IMAP_HOST":    "127.0.0.1",
		"IMAP_PORT":    "993",
		"IMAP_USER":    "parity-user",
		"IMAP_PASS":    "parity-pass",
		"IMAP_TLS":     "false",
		"IMAP_MAILBOX": "INBOX",
		"SMTP_HOST":    server.Host(),
		"SMTP_PORT":    server.Port(),
		"SMTP_USER":    "parity-user@example.com",
		"SMTP_PASS":    "parity-pass",
		"SMTP_SECURE":  "false",
		"SMTP_FROM":    "sender@example.com",
	}

	args := []string{"send", "--to", "recipient@example.com", "--subject", "Parity subject", "--body-file", bodyPath}
	js := runJS(t, "smtp.js", args, env)
	js.Stderr = normalizeParityStderr(js.Stderr)
	goOut := runGo(t, args, env)
	assertJSONParity(t, js, goOut, func(t *testing.T, jsVal, goVal any) {
		jsPayload := decodeSendPayload(t, jsVal)
		goPayload := decodeSendPayload(t, goVal)

		if !jsPayload.Success || !goPayload.Success {
			t.Fatalf("expected success true: js=%#v go=%#v", jsPayload, goPayload)
		}
		if !reflect.DeepEqual(stringSliceValue(t, jsPayload.To), []string{"recipient@example.com"}) || !reflect.DeepEqual(stringSliceValue(t, goPayload.To), []string{"recipient@example.com"}) {
			t.Fatalf("unexpected recipients: js=%#v go=%#v", jsPayload.To, goPayload.To)
		}
	})

	if got := strings.TrimSpace(server.LastMessage()); got == "" {
		t.Fatal("smtp server did not capture a message")
	}
}

func assertExactParity(t *testing.T, js, goOut commandResult) {
	t.Helper()
	if js.ExitCode != goOut.ExitCode {
		t.Fatalf("exit code mismatch:\njs=%d\ngo=%d\njs stderr=%q\ngo stderr=%q", js.ExitCode, goOut.ExitCode, js.Stderr, goOut.Stderr)
	}
	if js.Stdout != goOut.Stdout {
		t.Fatalf("stdout mismatch:\njs=%q\ngo=%q", js.Stdout, goOut.Stdout)
	}
	if js.Stderr != goOut.Stderr {
		t.Fatalf("stderr mismatch:\njs=%q\ngo=%q", js.Stderr, goOut.Stderr)
	}
}

func assertJSONParity(t *testing.T, js, goOut commandResult, validate func(*testing.T, any, any)) {
	t.Helper()
	if js.ExitCode != goOut.ExitCode {
		t.Fatalf("exit code mismatch:\njs=%d\ngo=%d\njs stderr=%q\ngo stderr=%q", js.ExitCode, goOut.ExitCode, js.Stderr, goOut.Stderr)
	}
	if normalizeParityStderr(js.Stderr) != normalizeParityStderr(goOut.Stderr) {
		t.Fatalf("stderr mismatch:\njs=%q\ngo=%q", js.Stderr, goOut.Stderr)
	}
	jsJSON := decodeJSON(t, js.Stdout)
	goJSON := decodeJSON(t, goOut.Stdout)
	validate(t, jsJSON, goJSON)
}

func normalizeParityStderr(stderr string) string {
	return strings.ReplaceAll(stderr, "SMTP server is ready to send\n", "")
}

func decodeJSON(t *testing.T, raw string) any {
	t.Helper()
	var out any
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &out); err != nil {
		t.Fatalf("decode json: %v\nraw=%s", err, raw)
	}
	return out
}

func decodeMailboxList(t *testing.T, raw any) []mailboxEntry {
	t.Helper()
	data, err := json.Marshal(canonicalizeJSON(raw))
	if err != nil {
		t.Fatalf("marshal mailbox list: %v", err)
	}
	var items []mailboxEntry
	if err := json.Unmarshal(data, &items); err != nil {
		t.Fatalf("unmarshal mailbox list: %v", err)
	}
	for i := range items {
		if items[i].SpecialUse == "" {
			items[i].SpecialUse = nil
		}
	}
	return items
}

func decodeSendPayload(t *testing.T, raw any) sendPayload {
	t.Helper()
	data, err := json.Marshal(canonicalizeJSON(raw))
	if err != nil {
		t.Fatalf("marshal send payload: %v", err)
	}
	var payload sendPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal send payload: %v", err)
	}
	return payload
}

func stringSliceValue(t *testing.T, v any) []string {
	t.Helper()
	switch typed := v.(type) {
	case string:
		if typed == "" {
			return nil
		}
		return []string{typed}
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			str, ok := item.(string)
			if !ok {
				t.Fatalf("expected string recipient, got %T", item)
			}
			out = append(out, str)
		}
		return out
	case []string:
		return typed
	case nil:
		return nil
	default:
		t.Fatalf("unexpected recipient type %T", v)
		return nil
	}
}

func canonicalizeJSON(v any) any {
	switch typed := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for k, val := range typed {
			out[canonicalKey(k)] = canonicalizeJSON(val)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, val := range typed {
			out[i] = canonicalizeJSON(val)
		}
		return out
	default:
		return typed
	}
}

func canonicalKey(key string) string {
	switch key {
	case "Success":
		return "success"
	case "MessageID":
		return "messageId"
	case "Response":
		return "response"
	case "To":
		return "to"
	case "Name":
		return "name"
	case "Delimiter":
		return "delimiter"
	case "Attributes":
		return "attributes"
	case "SpecialUse":
		return "specialUse"
	case "UID":
		return "uid"
	case "Seq":
		return "seq"
	case "Flags":
		return "flags"
	case "From":
		return "from"
	case "Subject":
		return "subject"
	case "Date":
		return "date"
	case "Text":
		return "text"
	case "HTML":
		return "html"
	case "Snippet":
		return "snippet"
	case "Attachments":
		return "attachments"
	case "Filename":
		return "filename"
	case "ContentType":
		return "contentType"
	case "Content":
		return "content"
	case "CID":
		return "cid"
	default:
		if key == "" {
			return key
		}
		return strings.ToLower(key[:1]) + key[1:]
	}
}

func fixturePath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "tests", "parity", "testdata", name)
}

func readFixtureString(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(fixturePath(t, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func runJS(t *testing.T, script string, args []string, env map[string]string) commandResult {
	t.Helper()
	return runCommand(t, "node", append([]string{filepath.Join(repoRoot(t), "scripts", script)}, args...), env)
}

func runGo(t *testing.T, args []string, env map[string]string) commandResult {
	t.Helper()
	return runCommand(t, buildGoBinary(t), args, env)
}

func runCommand(t *testing.T, exe string, args []string, env map[string]string) commandResult {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = repoRoot(t)
	cmd.Env = mergeEnv(env)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("run %s %v: %v", exe, args, err)
		}
	}
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("command timed out: %s %v", exe, args)
	}

	return commandResult{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: exitCode}
}

func mergeEnv(extra map[string]string) []string {
	base := os.Environ()
	if len(extra) == 0 {
		return base
	}

	keys := make([]string, 0, len(extra))
	for key := range extra {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		prefix := key + "="
		filtered := base[:0]
		for _, item := range base {
			if !strings.HasPrefix(item, prefix) {
				filtered = append(filtered, item)
			}
		}
		base = filtered
	}
	for _, key := range keys {
		base = append(base, key+"="+extra[key])
	}
	return base
}

func repoRoot(t *testing.T) string {
	t.Helper()
	repoRootOnce.Do(func() {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			repoRootErr = fmt.Errorf("runtime.Caller failed")
			return
		}
		repoRootPath = filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	})
	if repoRootErr != nil {
		t.Fatal(repoRootErr)
	}
	return repoRootPath
}

func buildGoBinary(t *testing.T) string {
	t.Helper()
	goBinaryOnce.Do(func() {
		binDir, err := os.MkdirTemp("", "zephyr-mail-parity-bin-")
		if err != nil {
			goBinaryErr = err
			return
		}
		goBinaryPath = filepath.Join(binDir, "zephyr-mail")
		cmd := exec.Command("go", "build", "-o", goBinaryPath, "./cmd/zephyr-mail")
		cmd.Dir = repoRoot(t)
		cmd.Env = os.Environ()
		if out, err := cmd.CombinedOutput(); err != nil {
			goBinaryErr = fmt.Errorf("go build: %v\n%s", err, string(out))
		}
	})
	if goBinaryErr != nil {
		t.Fatal(goBinaryErr)
	}
	return goBinaryPath
}

func ParityMatrixComplete(matrixPath string) bool {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return false
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	content, err := os.ReadFile(filepath.Join(root, filepath.Clean(matrixPath)))
	if err != nil {
		return false
	}

	requiredScenarios := []string{
		"unknown-command",
		"check success",
		"fetch missing UID",
		"fetch success",
		"download success",
		"search unseen recent",
		"mark-read success",
		"mark-unread success",
		"list-mailboxes success",
		"send missing `--to`",
		"send success",
		"test success",
	}

	text := string(content)
	for _, scenario := range requiredScenarios {
		found := false
		for _, line := range strings.Split(text, "\n") {
			if strings.Contains(line, "| "+scenario+" |") && strings.Contains(line, "| PASS |") {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

type imapParityServer struct {
	ln        net.Listener
	mailboxes []mailboxEntry
	mu        sync.Mutex
	closed    bool
}

func newIMAPParityServer(t *testing.T) *imapParityServer {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	srv := &imapParityServer{
		ln: ln,
		mailboxes: []mailboxEntry{{
			Name:       "Reports",
			Delimiter:  "/",
			Attributes: []string{},
			SpecialUse: nil,
		}},
	}
	go srv.serve()
	return srv
}

func (s *imapParityServer) Host() string {
	host, _, err := net.SplitHostPort(s.ln.Addr().String())
	if err != nil {
		return "127.0.0.1"
	}
	return host
}

func (s *imapParityServer) Port() string {
	_, port, err := net.SplitHostPort(s.ln.Addr().String())
	if err != nil {
		return "0"
	}
	return port
}

func (s *imapParityServer) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	_ = s.ln.Close()
}

func (s *imapParityServer) serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		go s.serveConn(conn)
	}
}

func (s *imapParityServer) serveConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	writeLine(writer, "* OK parity IMAP server ready")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		tag := fields[0]
		cmd := strings.ToUpper(fields[1])
		switch cmd {
		case "CAPABILITY":
			writeLine(writer, "* CAPABILITY IMAP4rev1 ID NAMESPACE UIDPLUS")
			writeLine(writer, tag+" OK CAPABILITY completed")
		case "LOGIN":
			writeLine(writer, tag+" OK LOGIN completed")
		case "ID":
			writeLine(writer, tag+" OK ID completed")
		case "NAMESPACE":
			writeLine(writer, "* NAMESPACE ((\"\" \"/\")) NIL NIL")
			writeLine(writer, tag+" OK NAMESPACE completed")
		case "SELECT", "EXAMINE":
			writeLine(writer, "* FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)")
			writeLine(writer, "* 1 EXISTS")
			writeLine(writer, "* 0 RECENT")
			writeLine(writer, tag+" OK [READ-WRITE] "+cmd+" completed")
		case "LIST":
			for _, mailbox := range s.mailboxes {
				writeLine(writer, fmt.Sprintf("* LIST () \"%s\" \"%s\"", mailbox.Delimiter, mailbox.Name))
			}
			writeLine(writer, tag+" OK LIST completed")
		case "LOGOUT":
			writeLine(writer, "* BYE logging out")
			writeLine(writer, tag+" OK LOGOUT completed")
			return
		default:
			writeLine(writer, tag+" OK "+cmd+" completed")
		}
	}
}

type smtpParityServer struct {
	ln      net.Listener
	mu      sync.Mutex
	closed  bool
	message string
}

func newSMTPParityServer(t *testing.T) *smtpParityServer {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := &smtpParityServer{ln: ln}
	go srv.serve()
	return srv
}

func (s *smtpParityServer) Host() string {
	host, _, err := net.SplitHostPort(s.ln.Addr().String())
	if err != nil {
		return "127.0.0.1"
	}
	return host
}

func (s *smtpParityServer) Port() string {
	_, port, err := net.SplitHostPort(s.ln.Addr().String())
	if err != nil {
		return "0"
	}
	return port
}

func (s *smtpParityServer) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	_ = s.ln.Close()
}

func (s *smtpParityServer) LastMessage() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.message
}

func (s *smtpParityServer) serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		go s.serveConn(conn)
	}
}

func (s *smtpParityServer) serveConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	writeLine(writer, "220 parity SMTP server ready")
	var authed bool

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		upper := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(upper, "EHLO"):
			writeLine(writer, "250-localhost")
			writeLine(writer, "250-AUTH PLAIN LOGIN")
			writeLine(writer, "250 SIZE 35882577")
		case strings.HasPrefix(upper, "HELO"):
			writeLine(writer, "250 localhost")
		case strings.HasPrefix(upper, "AUTH PLAIN"):
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				_, _ = base64.StdEncoding.DecodeString(fields[2])
			}
			authed = true
			writeLine(writer, "235 2.7.0 Authentication successful")
		case strings.HasPrefix(upper, "AUTH LOGIN"):
			writeLine(writer, "334 VXNlcm5hbWU6")
			if _, err := reader.ReadString('\n'); err != nil {
				return
			}
			writeLine(writer, "334 UGFzc3dvcmQ6")
			if _, err := reader.ReadString('\n'); err != nil {
				return
			}
			authed = true
			writeLine(writer, "235 2.7.0 Authentication successful")
		case strings.HasPrefix(upper, "MAIL FROM:"):
			if !authed {
				writeLine(writer, "530 5.7.0 Authentication required")
				continue
			}
			writeLine(writer, "250 2.1.0 OK")
		case strings.HasPrefix(upper, "RCPT TO:"):
			writeLine(writer, "250 2.1.5 OK")
		case upper == "DATA":
			writeLine(writer, "354 End data with <CR><LF>.<CR><LF>")
			var data bytes.Buffer
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				trimmed := strings.TrimRight(dataLine, "\r\n")
				if trimmed == "." {
					break
				}
				data.WriteString(dataLine)
			}
			s.mu.Lock()
			s.message = data.String()
			s.mu.Unlock()
			writeLine(writer, "250 2.0.0 Queued as parity-message")
		case upper == "RSET":
			writeLine(writer, "250 2.0.0 OK")
		case upper == "QUIT":
			writeLine(writer, "221 2.0.0 Bye")
			return
		default:
			writeLine(writer, "250 OK")
		}
	}
}

func writeLine(writer *bufio.Writer, line string) {
	_, _ = writer.WriteString(line + "\r\n")
	_ = writer.Flush()
}
