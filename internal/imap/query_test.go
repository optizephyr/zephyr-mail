package imap

import "testing"

func TestBuildSearchCriteriaWithUnseenAndSubject(t *testing.T) {
	c := BuildSearchCriteria(SearchOptions{Unseen: true, Subject: "hello"})
	if !c.Unseen || c.Subject != "hello" {
		t.Fatalf("unexpected criteria")
	}
}

func TestBuildSearchCriteriaDefaultsToAllWhenEmpty(t *testing.T) {
	c := BuildSearchCriteria(SearchOptions{})
	if !c.All {
		t.Fatal("expected all=true when no filters")
	}
}

func TestBuildSearchCriteriaRecentSetsSince(t *testing.T) {
	c := BuildSearchCriteria(SearchOptions{Recent: "30m"})
	if c.Since == "" {
		t.Fatal("expected since to be set from recent")
	}
}
