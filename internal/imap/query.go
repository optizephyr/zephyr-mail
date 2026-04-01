package imap

func BuildSearchCriteria(options SearchOptions) (SearchCriteria, error) {
	criteria := SearchCriteria{
		Unseen:   options.Unseen,
		Seen:     options.Seen,
		Flagged:  options.Flagged,
		Answered: options.Answered,
		From:     options.From,
		To:       options.To,
		Subject:  options.Subject,
		Before:   options.Before,
		UID:      options.UID,
	}

	if options.Recent != "" {
		since, err := ParseRelativeTime(options.Recent)
		if err != nil {
			return SearchCriteria{}, err
		}
		criteria.Since = since
	} else {
		criteria.Since = options.Since
	}

	if !hasSearchFilters(criteria) {
		criteria.All = true
	}

	return criteria, nil
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
