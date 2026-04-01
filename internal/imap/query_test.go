package imap

import "testing"

func TestBuildSearchCriteriaWithUnseenAndSubject(t *testing.T) {
	c, err := BuildSearchCriteria(SearchOptions{Unseen: true, Subject: "hello"})
	if err != nil {
		t.Fatal(err)
	}

	if !c.Unseen || c.Subject != "hello" {
		t.Fatalf("unexpected criteria")
	}
}

func TestBuildSearchCriteriaDefaultsToAllWhenEmpty(t *testing.T) {
	c, err := BuildSearchCriteria(SearchOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if !c.All {
		t.Fatal("expected all=true when no filters")
	}
}

func TestBuildSearchCriteriaRecentSetsSince(t *testing.T) {
	c, err := BuildSearchCriteria(SearchOptions{Recent: "30m"})
	if err != nil {
		t.Fatal(err)
	}

	if c.Since == "" {
		t.Fatal("expected since to be set from recent")
	}
}

func TestBuildSearchCriteriaRecentInvalidReturnsError(t *testing.T) {
	_, err := BuildSearchCriteria(SearchOptions{Recent: "nope"})
	if err == nil {
		t.Fatal("expected error for invalid recent format")
	}

	if err.Error() != "Invalid time format. Use: 30m, 2h, 7d" {
		t.Fatalf("unexpected error: %v", err)
	}
}
