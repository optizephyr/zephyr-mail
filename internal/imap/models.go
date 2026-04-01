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
	Mailbox  string
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
	Limit    int
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

type CheckOptions struct {
	Mailbox   string
	Limit     int
	Recent    string
	UnseenRaw string
}

type DownloadResult struct {
	UID        string
	Downloaded []DownloadedFile
	Message    string
}

type DownloadedFile struct {
	Filename string
	Path     string
	Size     int64
}

type Message struct {
	UID         string
	Seq         int
	Flags       []string
	From        string
	To          string
	Subject     string
	Date        string
	Text        string
	HTML        string
	Snippet     string
	Attachments []Attachment
}

type Attachment struct {
	Filename    string
	ContentType string
	Size        int64
	Content     []byte
	CID         string
}

type MailboxInfo struct {
	Name       string
	Delimiter  string
	Attributes []string
	SpecialUse string
}
