package main

import (
	"os"
	"strings"
	"testing"
)

func TestValidateAgents_ValidEntriesNoViolations(t *testing.T) {
	agents := []Agent{
		{Owner: "aider-ai", Repo: "aider", Category: "cli"},
		{Owner: "cline", Repo: "cline", Category: "extension"},
		{Owner: "a", Repo: "b.c-d_e", Category: "web"},
	}
	violations := validateAgents(agents)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidateAgents_EmptyOwner(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "", Repo: "repo1", Category: "cli"}})
	if !anyContains(violations, "owner is empty") {
		t.Errorf("expected 'owner is empty' violation, got %v", violations)
	}
}

func TestValidateAgents_InvalidOwnerChars(t *testing.T) {
	tests := []string{"foo_bar", "-leading", "trailing-", "double--hyphen", "foo owner"}
	for _, owner := range tests {
		t.Run(owner, func(t *testing.T) {
			violations := validateAgents([]Agent{{Owner: owner, Repo: "repo1", Category: "cli"}})
			if !anyContains(violations, "does not look like a valid GitHub username/org") {
				t.Errorf("owner %q: expected invalid-owner violation, got %v", owner, violations)
			}
		})
	}
}

func TestValidateAgents_EmptyRepo(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "", Category: "cli"}})
	if !anyContains(violations, "repo is empty") {
		t.Errorf("expected 'repo is empty' violation, got %v", violations)
	}
}

func TestValidateAgents_InvalidRepoChars(t *testing.T) {
	tests := []string{"foo/bar", "foo bar", "foo@bar", "foo#bar"}
	for _, repo := range tests {
		t.Run(repo, func(t *testing.T) {
			violations := validateAgents([]Agent{{Owner: "org", Repo: repo, Category: "cli"}})
			if !anyContains(violations, "characters not allowed in a GitHub repo name") {
				t.Errorf("repo %q: expected invalid-repo violation, got %v", repo, violations)
			}
		})
	}
}

func TestValidateAgents_MissingCategory(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Category: ""}})
	if !anyContains(violations, "category is required") {
		t.Errorf("expected 'category is required' violation, got %v", violations)
	}
}

func TestValidateAgents_InvalidCategory(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Category: "framework"}})
	if !anyContains(violations, "is not one of cli, ide, extension, library, research, web") {
		t.Errorf("expected invalid-category violation, got %v", violations)
	}
}

func TestValidateAgents_AllValidCategoriesAccepted(t *testing.T) {
	for _, cat := range []string{"cli", "ide", "extension", "library", "research", "web"} {
		t.Run(cat, func(t *testing.T) {
			violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Category: cat}})
			if len(violations) != 0 {
				t.Errorf("category %q: expected no violations, got %v", cat, violations)
			}
		})
	}
}

func TestValidateAgents_DuplicateCaseInsensitive(t *testing.T) {
	agents := []Agent{
		{Owner: "Foo", Repo: "Bar", Category: "cli"},
		{Owner: "foo", Repo: "bar", Category: "cli"},
	}
	violations := validateAgents(agents)
	if !anyContains(violations, "duplicate of entry 0") {
		t.Errorf("expected duplicate violation referencing entry 0, got %v", violations)
	}
}

func TestValidateAgents_NoDuplicateForDistinctRepos(t *testing.T) {
	agents := []Agent{
		{Owner: "foo", Repo: "bar", Category: "cli"},
		{Owner: "foo", Repo: "baz", Category: "cli"},
	}
	violations := validateAgents(agents)
	if len(violations) != 0 {
		t.Errorf("expected no violations for distinct repos, got %v", violations)
	}
}

func TestValidateAgents_CollectsAllViolationsNotJustFirst(t *testing.T) {
	agents := []Agent{
		{Owner: "", Repo: "", Category: ""},
		{Owner: "org", Repo: "repo1", Category: "not-a-category"},
	}
	violations := validateAgents(agents)
	// entry 0: owner empty, repo empty, category empty = 3 violations.
	// entry 1: invalid category = 1 violation.
	if len(violations) != 4 {
		t.Errorf("expected 4 violations collected across both entries, got %d: %v", len(violations), violations)
	}
}

func TestRunCheck_CurrentAgentsYML(t *testing.T) {
	// The real data/agents.yml must always pass -check; this is the
	// regression guard for that invariant.
	if err := runCheck("data/agents.yml"); err != nil {
		t.Errorf("expected data/agents.yml to pass validation, got: %v", err)
	}
}

func TestRunCheck_BadYAMLReportsViolationsAndError(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/bad-agents.yml"
	// loadAgents itself already rejects empty owner/repo, so this fixture
	// covers the violation classes only validateAgents catches: bad owner
	// chars, bad repo chars, invalid category, and a case-insensitive dup.
	content := `agents:
  - owner: "bad owner"
    repo: "valid-repo"
    category: cli
  - owner: "org"
    repo: "bad/repo"
    category: not-a-real-category
  - owner: "Dup"
    repo: "Repo"
    category: web
  - owner: "dup"
    repo: "repo"
    category: web
`
	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	err := runCheck(tmpFile)
	if err == nil {
		t.Fatal("expected runCheck to fail on invalid agents.yml, got nil error")
	}
	if !strings.Contains(err.Error(), "violation(s) found") {
		t.Errorf("expected error to summarize violation count, got: %v", err)
	}
}

func TestRunCheck_MissingFile(t *testing.T) {
	err := runCheck("/nonexistent/path/agents.yml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func anyContains(list []string, substr string) bool {
	for _, s := range list {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
