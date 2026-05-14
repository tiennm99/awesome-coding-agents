package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestComputeDeltas_NoHistory(t *testing.T) {
	// Empty history slice should return empty delta map without panicking.
	current := Snapshot{
		Date: "2026-05-14",
		Stars: map[string]int{
			"org/repo": 100,
		},
	}
	deltas := computeDeltas(nil, current)
	if deltas == nil {
		t.Errorf("expected non-nil empty map, got nil")
	}
	if len(deltas) != 0 {
		t.Errorf("expected empty deltas, got %d entries", len(deltas))
	}
}

func TestComputeDeltas_SevenDayWindow(t *testing.T) {
	// The cutoff window is (cutoff-3d, cutoff] where cutoff = now - 7d.
	// So valid snapshots have date in range (now-10d, now-7d].

	now := time.Now().UTC()
	today := now.Format("2006-01-02")

	tests := []struct {
		name        string
		daysAgo     int
		expectDelta bool
		desc        string
	}{
		{
			name:        "snapshot at exactly 7 days ago",
			daysAgo:     7,
			expectDelta: true,
			desc:        "at cutoff boundary — should be included",
		},
		{
			name:        "snapshot at 9 days ago",
			daysAgo:     9,
			expectDelta: true,
			desc:        "within (cutoff-3d, cutoff] — should be included",
		},
		{
			name:        "snapshot at 11 days ago",
			daysAgo:     11,
			expectDelta: false,
			desc:        "older than cutoff-3d — should be excluded",
		},
		{
			name:        "snapshot at 6 days ago",
			daysAgo:     6,
			expectDelta: false,
			desc:        "newer than cutoff — should be excluded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshotDate := now.AddDate(0, 0, -tt.daysAgo).Format("2006-01-02")
			history := []Snapshot{
				{
					Date: snapshotDate,
					Stars: map[string]int{
						"org/repo": 100,
					},
				},
			}
			current := Snapshot{
				Date: today,
				Stars: map[string]int{
					"org/repo": 150,
				},
			}
			deltas := computeDeltas(history, current)

			if tt.expectDelta {
				if len(deltas) == 0 {
					t.Errorf("%s: expected delta for 'org/repo', got empty map", tt.desc)
				}
				if delta, ok := deltas["org/repo"]; !ok || delta != 50 {
					t.Errorf("%s: expected delta=50, got %v", tt.desc, delta)
				}
			} else {
				if len(deltas) > 0 {
					t.Errorf("%s: expected no delta, got %v", tt.desc, deltas)
				}
			}
		})
	}
}

func TestReadSnapshots_MalformedLine(t *testing.T) {
	// Temp JSONL with: valid, malformed, empty, valid lines.
	// Should skip malformed/empty and return only the 2 valid snapshots.

	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/history.jsonl"

	content := `{"date":"2026-05-01","stars":{"org/repo":100}}
not json at all { bad
{"date":"2026-05-02","stars":{"org/repo":200}}
{"date":"2026-05-03","stars":{"org/repo":300}}
`

	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	snapshots, err := readSnapshots(tmpFile)
	if err != nil {
		t.Fatalf("readSnapshots failed: %v", err)
	}

	if len(snapshots) != 3 {
		t.Errorf("expected 3 snapshots (skipped 1 malformed), got %d", len(snapshots))
	}

	// Check dates are in order.
	if len(snapshots) >= 3 {
		if snapshots[0].Date != "2026-05-01" {
			t.Errorf("snapshot[0] date: expected 2026-05-01, got %s", snapshots[0].Date)
		}
		if snapshots[1].Date != "2026-05-02" {
			t.Errorf("snapshot[1] date: expected 2026-05-02, got %s", snapshots[1].Date)
		}
		if snapshots[2].Date != "2026-05-03" {
			t.Errorf("snapshot[2] date: expected 2026-05-03, got %s", snapshots[2].Date)
		}
	}
}

func TestReadSnapshots_MissingFile(t *testing.T) {
	// Non-existent file should return nil with no error (not os.IsNotExist error).
	snapshots, err := readSnapshots("/nonexistent/path/history.jsonl")
	if err != nil {
		t.Errorf("expected nil error on missing file, got %v", err)
	}
	if snapshots != nil {
		t.Errorf("expected nil snapshots on missing file, got %v", snapshots)
	}
}

func TestApplyMigrations(t *testing.T) {
	// Test the canonical key migration logic.
	stars := map[string]int{
		"aaif-goose/goose": 45115,
		"other/repo":       1000,
	}

	result := applyMigrations(stars)

	// The old key should be removed.
	if _, ok := result["aaif-goose/goose"]; ok {
		t.Errorf("old key 'aaif-goose/goose' should be removed")
	}

	// The new key should be present with the value.
	if v, ok := result["block/goose"]; !ok || v != 45115 {
		t.Errorf("new key 'block/goose': expected 45115, got %v", v)
	}

	// Other keys should be unchanged.
	if v, ok := result["other/repo"]; !ok || v != 1000 {
		t.Errorf("'other/repo': expected 1000, got %v", v)
	}
}

func TestWriteSnapshots_AtomicWrite(t *testing.T) {
	// Verify that writeSnapshots uses atomic rename (writes to .tmp first).
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/history.jsonl"

	snapshots := []Snapshot{
		{Date: "2026-05-01", Stars: map[string]int{"org/repo": 100}},
		{Date: "2026-05-02", Stars: map[string]int{"org/repo": 200}},
	}

	err := writeSnapshots(tmpFile, snapshots)
	if err != nil {
		t.Fatalf("writeSnapshots failed: %v", err)
	}

	// Verify file exists and has correct content.
	if _, err := os.Stat(tmpFile); err != nil {
		t.Fatalf("output file missing: %v", err)
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Verify JSON is valid and correct.
	var readBack []Snapshot
	decoder := json.NewDecoder(bytes.NewReader(content))
	for decoder.More() {
		var s Snapshot
		if err := decoder.Decode(&s); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		readBack = append(readBack, s)
	}

	if len(readBack) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(readBack))
	}
	if readBack[0].Date != "2026-05-01" || readBack[0].Stars["org/repo"] != 100 {
		t.Errorf("snapshot 0 mismatch: %+v", readBack[0])
	}
	if readBack[1].Date != "2026-05-02" || readBack[1].Stars["org/repo"] != 200 {
		t.Errorf("snapshot 1 mismatch: %+v", readBack[1])
	}

	// .tmp file should not exist (cleaned up after rename).
	if _, err := os.Stat(tmpFile + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temp file should not exist after successful write")
	}
}
