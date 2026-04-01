package imap

const DefaultMailbox = "INBOX"

var IMAPID = map[string]string{
	"name":          "moltbot",
	"version":       "0.0.1",
	"vendor":        "netease",
	"support-email": "kefu@188.com",
}

type ClientConfig struct {
	Host               string
	Port               int
	Username           string
	Password           string
	TLS                bool
	RejectUnauthorized bool
	Mailbox            string
}

type SearchOptions struct {
	Unseen   bool
	Seen     bool
	Flagged  bool
	Answered bool
	From     string
	To       string
	Subject  string
	Recent   string
	Since    string
	Before   string
	UID      string
}

type SearchCriteria struct {
	All      bool
	Unseen   bool
	Seen     bool
	Flagged  bool
	Answered bool
	From     string
	To       string
	Subject  string
	Since    string
	Before   string
	UID      string
}
