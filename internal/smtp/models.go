package smtp

type Config struct {
	Host               string
	Port               int
	Secure             bool
	RejectUnauthorized bool
	Username           string
	Password           string
	From               string
}

type SendRequest struct {
	From        string
	To          string
	Cc          string
	Bcc         string
	Subject     string
	SubjectFile string
	Body        string
	BodyFile    string
	HTMLFile    string
	HTML        bool
	Attach      string
}

type Attachment struct {
	Filename string
	Path     string
}

type ResolvedSendRequest struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	TextBody    string
	HTMLBody    string
	Attachments []Attachment
}

type SendResult struct {
	Success   bool
	MessageID string
	Response  string
	To        []string
}
