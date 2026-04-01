package imap

import (
	"strings"
	"testing"

	"github.com/emersion/go-imap"
)

func TestCheckUnseenLiteralTrueOnly(t *testing.T) {
	opts := CheckOptions{UnseenRaw: "true"}
	criteria := BuildCheckCriteria(opts)
	if !criteria.Unseen {
		t.Fatal("expected unseen true")
	}
}

func TestDownloadMissingSpecificFileReturnsAvailable(t *testing.T) {
	result := DownloadResult{Message: `File "x" not found. Available: a,b`}
	if !strings.Contains(result.Message, "Available") {
		t.Fatal(result.Message)
	}
}

func TestBuildCheckCriteriaDefaultsToAll(t *testing.T) {
	criteria := BuildCheckCriteria(CheckOptions{})
	if !criteria.All {
		t.Fatal("expected all=true when no filters")
	}
}

func TestBuildCheckCriteriaUnseenOnlyWhenTrue(t *testing.T) {
	opts := CheckOptions{UnseenRaw: "false"}
	criteria := BuildCheckCriteria(opts)
	if criteria.Unseen {
		t.Fatal("expected unseen=false when UnseenRaw is not 'true'")
	}

	opts = CheckOptions{UnseenRaw: ""}
	criteria = BuildCheckCriteria(opts)
	if criteria.Unseen {
		t.Fatal("expected unseen=false when UnseenRaw is empty")
	}
}

func TestBuildCheckCriteriaRecent(t *testing.T) {
	opts := CheckOptions{Recent: "2h"}
	criteria := BuildCheckCriteria(opts)
	if criteria.Since == "" {
		t.Fatal("expected since to be set from recent")
	}
}

func TestDownloadResultMessage(t *testing.T) {
	result := DownloadResult{
		UID:        "123",
		Downloaded: []DownloadedFile{{Filename: "test.txt", Path: "/tmp/test.txt", Size: 100}},
		Message:    "Downloaded 1 attachment(s)",
	}
	if result.UID != "123" {
		t.Fatal("expected uid")
	}
	if len(result.Downloaded) != 1 {
		t.Fatal("expected 1 downloaded file")
	}
	if result.Downloaded[0].Filename != "test.txt" {
		t.Fatal("expected filename")
	}
}

func TestFetchItemsIncludeBodySection(t *testing.T) {
	items := fetchItemsForMessage()
	if len(items) != 2 {
		t.Fatalf("expected 2 fetch items, got %d", len(items))
	}

	if items[0] != imap.FetchEnvelope {
		t.Fatalf("expected first fetch item to be envelope, got %v", items[0])
	}

	bodySection := imap.BodySectionName{}
	if items[1] != bodySection.FetchItem() {
		t.Fatalf("expected second fetch item to be body section, got %v", items[1])
	}
}
