package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ownerPattern approximates GitHub's username/org rules: alphanumeric runs
// separated by single hyphens — equivalent to
// "^[A-Za-z0-9](?:[A-Za-z0-9]|-(?=[A-Za-z0-9]))*$" (no leading/trailing
// hyphen, no consecutive hyphens) but written without lookahead, which Go's
// RE2-based regexp engine doesn't support.
var ownerPattern = regexp.MustCompile(`^[A-Za-z0-9]+(-[A-Za-z0-9]+)*$`)

// repoPattern approximates GitHub's repo name rules: alphanumeric, dot,
// underscore, hyphen.
var repoPattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// validCategories mirrors the category enum documented in
// templates/readme.tmpl's Contributing section.
var validCategories = map[string]bool{
	"cli":       true,
	"ide":       true,
	"extension": true,
	"library":   true,
	"research":  true,
	"web":       true,
}

// validateAgents checks data/agents.yml entries offline (no network, no
// token) and returns every violation found — not just the first — so a
// contributor sees the complete list of fixes needed in one pass.
func validateAgents(agents []Agent) []string {
	var violations []string
	seen := make(map[string]int, len(agents)) // lowercase "owner/repo" -> first index seen

	for i, a := range agents {
		ref := fmt.Sprintf("entry %d (owner=%q repo=%q)", i, a.Owner, a.Repo)

		switch {
		case strings.TrimSpace(a.Owner) == "":
			violations = append(violations, fmt.Sprintf("%s: owner is empty", ref))
		case !ownerPattern.MatchString(a.Owner):
			violations = append(violations, fmt.Sprintf("%s: owner %q does not look like a valid GitHub username/org (alphanumeric, single hyphens, no leading/trailing hyphen)", ref, a.Owner))
		}

		switch {
		case strings.TrimSpace(a.Repo) == "":
			violations = append(violations, fmt.Sprintf("%s: repo is empty", ref))
		case !repoPattern.MatchString(a.Repo):
			violations = append(violations, fmt.Sprintf("%s: repo %q contains characters not allowed in a GitHub repo name (allowed: letters, digits, '.', '_', '-')", ref, a.Repo))
		}

		switch {
		case strings.TrimSpace(a.Category) == "":
			violations = append(violations, fmt.Sprintf("%s: category is required (one of: cli, ide, extension, library, research, web)", ref))
		case !validCategories[a.Category]:
			violations = append(violations, fmt.Sprintf("%s: category %q is not one of cli, ide, extension, library, research, web", ref, a.Category))
		}

		key := strings.ToLower(a.Owner + "/" + a.Repo)
		if first, dup := seen[key]; dup {
			violations = append(violations, fmt.Sprintf("%s: duplicate of entry %d (case-insensitive owner/repo match)", ref, first))
		} else {
			seen[key] = i
		}
	}

	return violations
}

// runCheck loads data/agents.yml offline and validates it, printing every
// violation to stderr and returning a non-nil error if any are found. It
// never makes a network call, so it's safe to run against fork PRs without
// a GitHub token.
func runCheck(path string) error {
	agents, err := loadAgents(path)
	if err != nil {
		return fmt.Errorf("load %s: %w", path, err)
	}

	violations := validateAgents(agents)
	if len(violations) > 0 {
		for _, v := range violations {
			fmt.Fprintln(os.Stderr, v)
		}
		return fmt.Errorf("%d violation(s) found in %s", len(violations), path)
	}

	fmt.Printf("%s: %d agents valid\n", path, len(agents))
	return nil
}
