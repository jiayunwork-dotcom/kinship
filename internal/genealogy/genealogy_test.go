package genealogy

import (
	"strings"
	"testing"
)

func mustFamily(t *testing.T, lines []string) *Family {
	t.Helper()
	f, err := ParseFile(lines)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	return f
}

const goodFamily = `
PERSON george M 1920
PERSON carl M 1972
PERSON henry M 1975
PERSON diana F 1975
PERSON emma F 1998
PERSON frank M 2000
PERSON ivy F 2001
PERSON grace F 2022
PARENT george carl
PARENT george henry
PARENT carl emma
PARENT diana emma
PARENT carl frank
PARENT diana frank
PARENT henry ivy
PARENT emma grace
`

func TestParseFileAcceptsGoodFamily(t *testing.T) {
	f := mustFamily(t, strings.Split(goodFamily, "\n"))
	if got := len(f.Names()); got != 8 {
		t.Fatalf("Names() = %d persons, want 8", got)
	}
	p, ok := f.Person("emma")
	if !ok || p.Sex != "F" || p.Birth != 1998 {
		t.Errorf("Person(emma) = %+v, %v", p, ok)
	}
}

func TestParseFileRejectsBadRecords(t *testing.T) {
	bad := [][]string{
		{"PERSON alice X 1990"},
		{"PERSON alice F 19xx"},
		{"PERSON alice F 1990", "PERSON alice F 1991"},
		{"PARENT alice bob"},
		{"PERSON alice F 1990", "PARENT alice bob"},
		{"PERSON alice F 1990", "PARENT alice alice"},
		{"PERSON a F 1990", "PERSON b M 1990", "PARENT a b", "PARENT a b"},
		{"FAMILY alice"},
		{"# only a comment"},
	}
	for _, lines := range bad {
		if _, err := ParseFile(lines); err == nil {
			t.Errorf("ParseFile(%q) expected error", strings.Join(lines, " | "))
		}
	}
}

func TestParseFileForwardReference(t *testing.T) {
	f := mustFamily(t, []string{
		"PARENT ana bob",
		"PERSON ana F 1960",
		"PERSON bob M 1985",
	})
	ps, err := f.Parents("bob")
	if err != nil || len(ps) != 1 || ps[0] != "ana" {
		t.Errorf("Parents(bob) = %v, %v; want [ana]", ps, err)
	}
}

func TestParseFileErrorMentionsLine(t *testing.T) {
	_, err := ParseFile([]string{"PERSON alice F 1990", "", "PERSON alice F 1991"})
	if err == nil || !strings.Contains(err.Error(), "line 3") {
		t.Errorf("error should mention line 3, got %v", err)
	}
}

func TestAncestorsDistances(t *testing.T) {
	f := mustFamily(t, strings.Split(goodFamily, "\n"))
	anc, err := f.Ancestors("grace")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{
		"grace": 0, "emma": 1, "carl": 2, "diana": 2, "george": 3,
	}
	if len(anc) != len(want) {
		t.Fatalf("Ancestors(grace) = %v, want %v", anc, want)
	}
	for name, d := range want {
		if anc[name] != d {
			t.Errorf("Ancestors(grace)[%s] = %d, want %d", name, anc[name], d)
		}
	}
	if _, err := f.Ancestors("nobody"); err == nil {
		t.Error("Ancestors(unknown) expected error")
	}
}

func TestLCACases(t *testing.T) {
	f := mustFamily(t, strings.Split(goodFamily, "\n"))
	cases := []struct {
		a, b, lca string
		da, db    int
	}{
		{"grace", "george", "george", 3, 0},
		{"emma", "frank", "carl", 1, 1},
		{"emma", "ivy", "george", 2, 2},
		{"grace", "ivy", "george", 3, 2},
		{"grace", "henry", "george", 3, 1},
	}
	for _, c := range cases {
		lca, da, db, err := f.LCA(c.a, c.b)
		if err != nil {
			t.Fatalf("LCA(%s,%s): %v", c.a, c.b, err)
		}
		if lca != c.lca || da != c.da || db != c.db {
			t.Errorf("LCA(%s,%s) = (%s,%d,%d), want (%s,%d,%d)",
				c.a, c.b, lca, da, db, c.lca, c.da, c.db)
		}
	}
	if _, _, _, err := f.LCA("emma", "zzz"); err == nil {
		t.Error("LCA with unknown person expected error")
	}
}

func TestChildrenSorted(t *testing.T) {
	f := mustFamily(t, strings.Split(goodFamily, "\n"))
	kids, err := f.Children("carl")
	if err != nil {
		t.Fatal(err)
	}
	if len(kids) != 2 || kids[0] != "emma" || kids[1] != "frank" {
		t.Errorf("Children(carl) = %v, want [emma frank]", kids)
	}
	if _, err := f.Children("nobody"); err == nil {
		t.Error("Children(unknown) expected error")
	}
}

func TestSharedParents(t *testing.T) {
	f := mustFamily(t, strings.Split(goodFamily, "\n"))
	if n, _ := f.SharedParents("emma", "frank"); n != 2 {
		t.Errorf("SharedParents(emma,frank) = %d, want 2", n)
	}
	if n, _ := f.SharedParents("emma", "ivy"); n != 0 {
		t.Errorf("SharedParents(emma,ivy) = %d, want 0", n)
	}
}
