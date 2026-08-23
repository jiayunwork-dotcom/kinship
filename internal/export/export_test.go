package export

import (
	"bytes"
	"strings"
	"testing"

	"kinship/internal/genealogy"
)

func testFamily(t *testing.T) *genealogy.Family {
	t.Helper()
	lines := []string{
		"PERSON alice F 1950",
		"PERSON bob M 1948",
		"PERSON carol F 1975",
		"PARENT alice carol",
		"PARENT bob carol",
	}
	f, err := genealogy.ParseFile(lines)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestWriteDOT(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteDOT(&buf, testFamily(t)); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "digraph") {
		t.Fatal("expected digraph keyword in DOT output")
	}
	if !strings.Contains(out, "alice") {
		t.Fatal("expected alice in DOT output")
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, testFamily(t)); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "alice") || !strings.Contains(out, "bob") {
		t.Fatalf("JSON missing family members: %s", out)
	}
}

func TestWriteText(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteText(&buf, testFamily(t)); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty text output")
	}
}

func TestWriteTable(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteTable(&buf, testFamily(t)); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "alice") {
		t.Fatal("table should contain alice")
	}
}
