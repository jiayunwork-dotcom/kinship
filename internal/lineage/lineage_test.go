package lineage

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
		"PARENT george carl",
		"PARENT helen carl",
		"PARENT carl alice",
		"PARENT diana alice",
	}
	f, err := genealogy.ParseFile(lines)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestPatriline(t *testing.T) {
	f := testFamily(t)
	line, err := Patriline(f, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if line.Depth() < 2 {
		t.Fatalf("expected patriline depth >= 2, got %d", line.Depth())
	}
	if line.Root() != "george" {
		t.Fatalf("expected patriline root george, got %s", line.Root())
	}
}

func TestMatriline(t *testing.T) {
	f := testFamily(t)
	line, err := Matriline(f, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if line.Depth() < 1 {
		t.Fatalf("expected matriline depth >= 1, got %d", line.Depth())
	}
}

func TestPatrilinealDescendants(t *testing.T) {
	f := testFamily(t)
	desc, err := PatrilinealDescendants(f, "george")
	if err != nil {
		t.Fatal(err)
	}
	if len(desc) == 0 {
		t.Fatal("expected at least one patrilineal descendant")
	}
}
