package date_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gtk-ai/date/filter"
)

// --- Rewrite ---

func TestRewriteNoArgs(t *testing.T) {
	got, ok := filter.Rewrite(nil)
	if !ok {
		t.Fatal("expected rewrite when no format supplied")
	}
	if len(got) != 1 || !strings.HasPrefix(got[0], "+%") {
		t.Fatalf("unexpected rewritten args: %v", got)
	}
}

func TestRewriteWithFormatPassthrough(t *testing.T) {
	_, ok := filter.Rewrite([]string{"+%d/%m/%Y"})
	if ok {
		t.Fatal("must not rewrite when format already present")
	}
}

func TestRewriteWithUtcFlag(t *testing.T) {
	got, ok := filter.Rewrite([]string{"-u"})
	if !ok {
		t.Fatal("expected rewrite with -u and no format")
	}
	hasISO := false
	for _, a := range got {
		if strings.HasPrefix(a, "+%") {
			hasISO = true
		}
	}
	if !hasISO {
		t.Fatalf("ISO format not injected, got: %v", got)
	}
}

func TestRewriteWithFormatAndFlag(t *testing.T) {
	_, ok := filter.Rewrite([]string{"-u", "+%s"})
	if ok {
		t.Fatal("must not rewrite when format already present alongside flags")
	}
}

// --- FilterOutput ---

func TestFilterOutputTrimsNewline(t *testing.T) {
	out := filter.FilterOutput(nil, "2026-08-19T14:30:00Z\n\n", 0)
	if out != "2026-08-19T14:30:00Z\n" {
		t.Fatalf("unexpected: %q", out)
	}
}

func TestFilterOutputNonZeroExitPassthrough(t *testing.T) {
	out := filter.FilterOutput(nil, "date: invalid option\n", 1)
	if out != "date: invalid option\n" {
		t.Fatalf("unexpected: %q", out)
	}
}

// --- ID ---

func TestID(t *testing.T) {
	if filter.ID != "gtk-ai/gtkai-date" {
		t.Fatalf("ID %q does not follow author/gtkai-<command> rule", filter.ID)
	}
}

// --- gtkai.json manifest ---

func TestManifest(t *testing.T) {
	data, err := os.ReadFile("gtkai.json")
	if err != nil {
		t.Fatalf("read gtkai.json: %v", err)
	}

	var manifest struct {
		ID               string   `json:"id"`
		Filters          []string `json:"filters"`
		Platforms        []string `json:"platforms"`
		Contract         string   `json:"contract"`
		GtkaiCoreVersion struct {
			Version    string `json:"version"`
			Constraint string `json:"constraint"`
		} `json:"gtkai-core-version"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse gtkai.json: %v", err)
	}
	if manifest.ID != filter.ID {
		t.Fatalf("manifest id %q != code id %q", manifest.ID, filter.ID)
	}
	if len(manifest.Filters) != 1 || manifest.Filters[0] != "date" {
		t.Fatalf("unexpected filters list: %v", manifest.Filters)
	}
	if manifest.Contract != "subprocess/v1" {
		t.Fatalf("unexpected contract: %q", manifest.Contract)
	}
	if manifest.GtkaiCoreVersion.Version == "" {
		t.Fatal("gtkai-core-version.version must not be empty")
	}
	if manifest.GtkaiCoreVersion.Constraint != "min" && manifest.GtkaiCoreVersion.Constraint != "exact" {
		t.Fatalf("unexpected gtkai-core-version.constraint: %q", manifest.GtkaiCoreVersion.Constraint)
	}
	if len(manifest.Platforms) == 0 {
		t.Fatal("platforms must not be empty")
	}
}
