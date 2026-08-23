package pedigree

import (
	"testing"

	"kinship/internal/genealogy"
)

func testFamily(t *testing.T) *genealogy.Family {
	t.Helper()
	lines := []string{
		"PERSON george M 1920",
		"PERSON helen F 1922",
		"PERSON carl M 1945",
		"PERSON diana F 1947",
		"PERSON alice F 1970",
		"PERSON bob M 1995",
		"PARENT george carl",
		"PARENT helen carl",
		"PARENT carl alice",
		"PARENT diana alice",
		"PARENT alice bob",
	}
	f, err := genealogy.ParseFile(lines)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestAncestorChart(t *testing.T) {
	f := testFamily(t)
	nodes, err := AncestorChart(f, "alice", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) == 0 {
		t.Fatal("expected ancestor nodes")
	}
}

func TestDescendantTree(t *testing.T) {
	f := testFamily(t)
	nodes, err := DescendantTree(f, "george", 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) == 0 {
		t.Fatal("expected descendant nodes")
	}
}

func TestCompleteness(t *testing.T) {
	f := testFamily(t)
	comp, err := Completeness(f, "alice", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(comp) == 0 {
		t.Fatal("expected completeness data")
	}
}

func TestTotalAncestorCount(t *testing.T) {
	f := testFamily(t)
	count, err := TotalAncestorCount(f, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if count < 2 {
		t.Fatalf("expected >= 2 ancestors, got %d", count)
	}
}
