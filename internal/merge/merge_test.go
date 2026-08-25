package merge

import (
	"testing"
)

func TestMergeNoConflict(t *testing.T) {
	linesA := []string{
		"PERSON alice F 1950",
		"PERSON bob M 1948",
		"PERSON carol F 1975",
		"PARENT alice carol",
	}
	linesB := []string{
		"PERSON alice F 1950",
		"PERSON bob M 1948",
		"PERSON carol F 1975",
		"PARENT bob carol",
	}
	result, err := Merge(linesA, linesB)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasConflicts() {
		t.Fatalf("expected no conflicts, got: %s", result.ConflictSummary())
	}
}

func TestMergeWithConflict(t *testing.T) {
	linesA := []string{
		"PERSON alice F 1950",
	}
	linesB := []string{
		"PERSON alice M 1950",
	}
	result, err := Merge(linesA, linesB)
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasConflicts() {
		t.Fatal("expected conflicts for conflicting sex")
	}
}

func TestMergeFiles(t *testing.T) {
	contentA := "PERSON alice F 1950\nPERSON bob M 1948\nPERSON carol F 1975\nPARENT alice carol\n"
	contentB := "PERSON alice F 1950\nPERSON bob M 1948\nPERSON carol F 1975\nPARENT bob carol\n"
	result, err := MergeFiles(contentA, contentB)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasConflicts() {
		t.Fatal("unexpected conflicts")
	}
}
