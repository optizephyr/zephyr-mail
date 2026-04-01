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

func TestBuildSearchCriteriaRecentIgnoresSinceAndBefore(t *testing.T) {
	c, err := BuildSearchCriteria(SearchOptions{Recent: "30m", Since: "2026-01-01", Before: "2026-12-31"})
	if err != nil {
		t.Fatal(err)
	}

	if c.Since == "" {
		t.Fatal("expected since to be set from recent")
	}

	if c.Before != "" {
		t.Fatalf("expected before to be ignored when recent is set, got %q", c.Before)
	}
}

func TestBuildSearchCriteriaNormalizesSinceAndBefore(t *testing.T) {
	c, err := BuildSearchCriteria(SearchOptions{Since: "2026-04-01", Before: "2026-04-02"})
	if err != nil {
		t.Fatal(err)
	}

	if c.Since != "01-Apr-2026" {
		t.Fatalf("expected normalized since, got %q", c.Since)
	}

	if c.Before != "02-Apr-2026" {
		t.Fatalf("expected normalized before, got %q", c.Before)
	}
}

func TestBuildSearchCriteriaSeenOverridesUnseen(t *testing.T) {
	c, err := BuildSearchCriteria(SearchOptions{Unseen: true, Seen: true})
	if err != nil {
		t.Fatal(err)
	}

	if !c.Seen {
		t.Fatal("expected seen=true")
	}

	if c.Unseen {
		t.Fatal("expected unseen=false when seen=true")
	}
}
