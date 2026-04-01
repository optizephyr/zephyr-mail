package imap

import "time"

func BuildSearchCriteria(options SearchOptions) (SearchCriteria, error) {
	criteria := SearchCriteria{
		Flagged:  options.Flagged,
		Answered: options.Answered,
		From:     options.From,
		To:       options.To,
		Subject:  options.Subject,
		UID:      options.UID,
	}

	if options.Unseen {
		criteria.Unseen = true
	}
	if options.Seen {
		criteria.Seen = true
		criteria.Unseen = false
	}

	if options.Recent != "" {
		since, err := ParseRelativeTime(options.Recent)
		if err != nil {
			return SearchCriteria{}, err
		}
		criteria.Since = since
	} else {
		criteria.Since = normalizeIMAPDate(options.Since)
		criteria.Before = normalizeIMAPDate(options.Before)
	}

	if !hasSearchFilters(criteria) {
		criteria.All = true
	}

	return criteria, nil
}

func normalizeIMAPDate(value string) string {
	if value == "" {
		return ""
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return FormatIMAPDate(parsed)
		}
	}

	return value
}

func hasSearchFilters(c SearchCriteria) bool {
	return c.Unseen ||
		c.Seen ||
		c.Flagged ||
		c.Answered ||
		c.From != "" ||
		c.To != "" ||
		c.Subject != "" ||
		c.Since != "" ||
		c.Before != "" ||
		c.UID != ""
}
