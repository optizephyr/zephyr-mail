package imap

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

const (
	DefaultLimit = 10
	FlagSeen     = imap.SeenFlag
	FlagAnswered = imap.AnsweredFlag
	FlagFlagged  = imap.FlaggedFlag
	FlagDeleted  = imap.DeletedFlag
	FlagDraft    = imap.DraftFlag
)

func BuildCheckCriteria(opts CheckOptions) SearchCriteria {
	criteria := SearchCriteria{}

	if opts.UnseenRaw == "true" {
		criteria.Unseen = true
	}

	if opts.Recent != "" {
		since, err := ParseRelativeTime(opts.Recent)
		if err == nil {
			criteria.Since = since
		}
	}

	if !criteria.Unseen && criteria.Since == "" {
		criteria.All = true
	}

	return criteria
}

func buildSearchCriteriaFromSearchCriteria(sc SearchCriteria) *imap.SearchCriteria {
	criteria := imap.NewSearchCriteria()

	if sc.All {
		return criteria
	}

	if sc.Unseen {
		criteria.WithoutFlags = []string{FlagSeen}
	}
	if sc.Seen {
		criteria.WithFlags = []string{FlagSeen}
	}
	if sc.Flagged {
		criteria.WithFlags = append(criteria.WithFlags, FlagFlagged)
	}
	if sc.Answered {
		criteria.WithFlags = append(criteria.WithFlags, FlagAnswered)
	}
	if sc.From != "" {
		criteria.Body = []string{sc.From}
	}
	if sc.To != "" {
		criteria.Body = append(criteria.Body, sc.To)
	}
	if sc.Subject != "" {
		criteria.Body = append(criteria.Body, sc.Subject)
	}
	if sc.Since != "" {
		if sinceDate, err := parseIMAPDate(sc.Since); err == nil {
			criteria.Since = sinceDate
		}
	}
	if sc.Before != "" {
		if beforeDate, err := parseIMAPDate(sc.Before); err == nil {
			criteria.Before = beforeDate
		}
	}
	if sc.UID != "" {
		set, err := imap.ParseSeqSet(sc.UID)
		if err == nil {
			criteria.Uid = set
		}
	}

	return criteria
}

func parseIMAPDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"02-Jan-2006",
		"2006-01-02",
		"2006-01-02T15:04:05Z",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date: %s", dateStr)
}

type CheckServiceResult struct {
	Messages []Message
}

func Check(cl *client.Client, opts CheckOptions) (CheckServiceResult, error) {
	mailbox := opts.Mailbox
	if mailbox == "" {
		mailbox = DefaultMailbox
	}

	limit := opts.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	criteria := BuildCheckCriteria(opts)
	searchCriteria := buildSearchCriteriaFromSearchCriteria(criteria)

	_, err := cl.Select(mailbox, false)
	if err != nil {
		return CheckServiceResult{}, fmt.Errorf("failed to select mailbox: %w", err)
	}

	uids, err := cl.UidSearch(searchCriteria)
	if err != nil {
		return CheckServiceResult{}, fmt.Errorf("search failed: %w", err)
	}

	if len(uids) == 0 {
		return CheckServiceResult{Messages: []Message{}}, nil
	}

	messages, err := fetchMessages(cl, uids, limit)
	if err != nil {
		return CheckServiceResult{}, err
	}

	return CheckServiceResult{Messages: messages}, nil
}

func fetchMessages(cl *client.Client, uids []uint32, limit int) ([]Message, error) {
	if limit > 0 && len(uids) > limit {
		uids = uids[:limit]
	}

	seqSet := &imap.SeqSet{}
	seqSet.AddNum(uids...)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- cl.UidFetch(seqSet, fetchItemsForMessage(), messages)
	}()

	var results []Message
	for msg := range messages {
		if msg == nil {
			continue
		}
		m := parseMessage(msg)
		if m.UID != "" {
			results = append(results, m)
		}
	}

	if err := <-done; err != nil {
		return results, fmt.Errorf("fetch failed: %w", err)
	}

	return results, nil
}

func parseMessage(msg *imap.Message) Message {
	m := Message{
		UID:   fmt.Sprintf("%d", msg.Uid),
		Seq:   int(msg.SeqNum),
		Flags: msg.Flags,
	}

	if msg.Envelope != nil {
		m.From = extractAddress(msg.Envelope.From)
		m.To = extractAddress(msg.Envelope.To)
		if msg.Envelope.Subject != "" {
			m.Subject = msg.Envelope.Subject
		} else {
			m.Subject = "(no subject)"
		}
		if !msg.Envelope.Date.IsZero() {
			m.Date = msg.Envelope.Date.Format(time.RFC3339)
		}
	}

	section := &imap.BodySectionName{}
	if r := msg.GetBody(section); r != nil {
		mr, err := mail.CreateReader(r)
		if err == nil {
			for {
				part, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					break
				}

				switch header := part.Header.(type) {
				case *mail.InlineHeader:
					contentType, _, _ := header.ContentType()
					if strings.HasPrefix(contentType, "text/plain") {
						b, _ := io.ReadAll(part.Body)
						m.Text = string(b)
					} else if strings.HasPrefix(contentType, "text/html") {
						b, _ := io.ReadAll(part.Body)
						m.HTML = string(b)
					}
				}
			}
			mr.Close()
		}
	}

	m.Snippet = m.Text
	if len(m.Snippet) > 200 {
		m.Snippet = m.Snippet[:200]
	}

	return m
}

func extractAddress(addrs []*imap.Address) string {
	if len(addrs) == 0 {
		return "Unknown"
	}
	addr := addrs[0]
	if addr.PersonalName != "" {
		return fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName)
	}
	return fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
}

func Fetch(cl *client.Client, uid string, mailbox string) ([]Message, error) {
	if mailbox == "" {
		mailbox = DefaultMailbox
	}

	_, err := cl.Select(mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select mailbox: %w", err)
	}

	seqSet := &imap.SeqSet{}
	uidNum := uint32(0)
	fmt.Sscanf(uid, "%d", &uidNum)
	seqSet.AddNum(uidNum)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- cl.UidFetch(seqSet, fetchItemsForMessage(), messages)
	}()

	var results []Message
	for msg := range messages {
		if msg != nil {
			results = append(results, parseMessage(msg))
		}
	}

	if err := <-done; err != nil {
		return results, fmt.Errorf("fetch failed: %w", err)
	}

	return results, nil
}

func fetchItemsForMessage() []imap.FetchItem {
	section := imap.BodySectionName{}
	return []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}
}

func Download(cl *client.Client, uid string, mailbox string, outputDir string, filename string) (DownloadResult, error) {
	result := DownloadResult{UID: uid, Downloaded: []DownloadedFile{}}

	if mailbox == "" {
		mailbox = DefaultMailbox
	}

	_, err := cl.Select(mailbox, false)
	if err != nil {
		return result, fmt.Errorf("failed to select mailbox: %w", err)
	}

	seqSet := &imap.SeqSet{}
	uidNum := uint32(0)
	fmt.Sscanf(uid, "%d", &uidNum)
	seqSet.AddNum(uidNum)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- cl.UidFetch(seqSet, []imap.FetchItem{imap.FetchBodyStructure}, messages)
	}()

	var foundFilenames []string
	for msg := range messages {
		if msg != nil {
			bs := msg.BodyStructure
			if bs != nil && bs.Disposition != "" && bs.DispositionParams != nil {
				if attachFilename := bs.DispositionParams["filename"]; attachFilename != "" {
					foundFilenames = append(foundFilenames, attachFilename)
				}
			}
		}
	}

	err = <-done
	if err != nil {
		return result, fmt.Errorf("fetch failed: %w", err)
	}

	if len(foundFilenames) == 0 {
		result.Message = "No attachments found"
		return result, nil
	}

	if filename != "" {
		found := false
		for _, fn := range foundFilenames {
			if fn == filename {
				found = true
				break
			}
		}
		if !found {
			result.Message = fmt.Sprintf("File %q not found. Available: %v", filename, foundFilenames)
			return result, nil
		}
	}

	result.Message = fmt.Sprintf("Found %d attachment(s)", len(foundFilenames))
	return result, nil
}

type SearchServiceResult struct {
	Messages []Message
}

func Search(cl *client.Client, opts SearchOptions) (SearchServiceResult, error) {
	mailbox := opts.Mailbox
	if mailbox == "" {
		mailbox = DefaultMailbox
	}

	criteria, err := BuildSearchCriteria(opts)
	if err != nil {
		return SearchServiceResult{}, err
	}

	searchCriteria := buildSearchCriteriaFromSearchCriteria(criteria)

	_, err = cl.Select(mailbox, false)
	if err != nil {
		return SearchServiceResult{}, fmt.Errorf("failed to select mailbox: %w", err)
	}

	uids, err := cl.UidSearch(searchCriteria)
	if err != nil {
		return SearchServiceResult{}, fmt.Errorf("search failed: %w", err)
	}

	if len(uids) == 0 {
		return SearchServiceResult{Messages: []Message{}}, nil
	}

	limit := 100
	messages, err := fetchMessages(cl, uids, limit)
	if err != nil {
		return SearchServiceResult{}, err
	}

	return SearchServiceResult{Messages: messages}, nil
}

type MarkResult struct {
	UIDs   []string
	Action string
	Count  int
}

func MarkRead(cl *client.Client, uids []string, mailbox string) (MarkResult, error) {
	return updateFlags(cl, uids, mailbox, []string{FlagSeen}, "add")
}

func MarkUnread(cl *client.Client, uids []string, mailbox string) (MarkResult, error) {
	return updateFlags(cl, uids, mailbox, []string{FlagSeen}, "remove")
}

func updateFlags(cl *client.Client, uids []string, mailbox string, flags []string, action string) (MarkResult, error) {
	result := MarkResult{
		UIDs:   uids,
		Action: action,
		Count:  len(uids),
	}

	if mailbox == "" {
		mailbox = DefaultMailbox
	}

	_, err := cl.Select(mailbox, true)
	if err != nil {
		return result, fmt.Errorf("failed to select mailbox: %w", err)
	}

	seqSet := &imap.SeqSet{}
	for _, uidStr := range uids {
		var uidNum uint32
		fmt.Sscanf(uidStr, "%d", &uidNum)
		seqSet.AddNum(uidNum)
	}

	var storeItem imap.StoreItem
	switch action {
	case "add":
		storeItem = imap.FormatFlagsOp(imap.AddFlags, false)
	case "remove":
		storeItem = imap.FormatFlagsOp(imap.RemoveFlags, false)
	case "set":
		storeItem = imap.FormatFlagsOp(imap.SetFlags, false)
	default:
		return result, fmt.Errorf("unknown action: %s", action)
	}

	flagInterface := make([]interface{}, len(flags))
	for i, f := range flags {
		flagInterface[i] = f
	}

	updates := make(chan *imap.Message, 1)
	err = cl.UidStore(seqSet, storeItem, flagInterface, updates)
	_ = <-updates
	if err != nil {
		return result, fmt.Errorf("store failed: %w", err)
	}

	return result, nil
}

func ListMailboxes(cl *client.Client) ([]MailboxInfo, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- cl.List("", "*", mailboxes)
	}()

	var results []MailboxInfo
	for mb := range mailboxes {
		if mb == nil {
			continue
		}
		info := MailboxInfo{
			Name:       mb.Name,
			Delimiter:  string(mb.Delimiter),
			Attributes: []string{},
		}
		results = append(results, info)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("list failed: %w", err)
	}

	return results, nil
}
