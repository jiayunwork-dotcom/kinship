package relation

import (
	"strings"
	"testing"

	"kinship/internal/genealogy"
)

func mustFamily(t *testing.T, text string) *genealogy.Family {
	t.Helper()
	f, err := genealogy.ParseFile(strings.Split(text, "\n"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	return f
}

// Fixture with three generations, a half sibling, and an unrelated person:
//
//	   george
//	  /      \
//	carl     henry        kai (unrelated)
//	/  \  \     \
//
// emma frank hanna ivy
//
//	|
//
// grace
const fixture = `
PERSON george M 1920
PERSON carl M 1972
PERSON henry M 1975
PERSON diana F 1975
PERSON emma F 1998
PERSON frank M 2000
PERSON hanna F 2005
PERSON ivy F 2001
PERSON grace F 2022
PERSON kai M 1999
PARENT george carl
PARENT george henry
PARENT carl emma
PARENT diana emma
PARENT carl frank
PARENT diana frank
PARENT carl hanna
PARENT henry ivy
PARENT emma grace
`

func describe(t *testing.T, f *genealogy.Family, from, to string) string {
	t.Helper()
	term, err := Describe(f, from, to)
	if err != nil {
		t.Fatalf("Describe(%s,%s): %v", from, to, err)
	}
	return term
}

func TestSelfAndUnknown(t *testing.T) {
	f := mustFamily(t, fixture)
	if got := describe(t, f, "emma", "emma"); got != "self" {
		t.Errorf("self term = %q, want self", got)
	}
	if _, err := Describe(f, "emma", "nobody"); err == nil {
		t.Error("unknown `to` expected error")
	}
	if _, err := Describe(f, "nobody", "emma"); err == nil {
		t.Error("unknown `from` expected error")
	}
}

func TestAncestorTerms(t *testing.T) {
	f := mustFamily(t, fixture)
	cases := []struct{ from, to, want string }{
		{"grace", "emma", "mother"},
		{"grace", "carl", "grandfather"},
		{"grace", "diana", "grandmother"},
		{"grace", "george", "great-grandfather"},
	}
	for _, c := range cases {
		if got := describe(t, f, c.from, c.to); got != c.want {
			t.Errorf("Describe(%s,%s) = %q, want %q", c.from, c.to, got, c.want)
		}
	}
}

func TestDescendantTerms(t *testing.T) {
	f := mustFamily(t, fixture)
	cases := []struct{ from, to, want string }{
		{"emma", "grace", "daughter"},
		{"carl", "grace", "granddaughter"},
		{"george", "grace", "great-granddaughter"},
		{"henry", "ivy", "daughter"},
	}
	for _, c := range cases {
		if got := describe(t, f, c.from, c.to); got != c.want {
			t.Errorf("Describe(%s,%s) = %q, want %q", c.from, c.to, got, c.want)
		}
	}
}

func TestSiblingTerms(t *testing.T) {
	f := mustFamily(t, fixture)
	if got := describe(t, f, "emma", "frank"); got != "sibling" {
		t.Errorf("emma->frank = %q, want sibling (2 shared parents)", got)
	}
	// hanna has only carl as a registered parent, so she is a half
	// sibling of emma and frank.
	if got := describe(t, f, "emma", "hanna"); got != "half-sibling" {
		t.Errorf("emma->hanna = %q, want half-sibling", got)
	}
	if got := describe(t, f, "frank", "hanna"); got != "half-sibling" {
		t.Errorf("frank->hanna = %q, want half-sibling", got)
	}
}

func TestAuntNieceTerms(t *testing.T) {
	f := mustFamily(t, fixture)
	cases := []struct{ from, to, want string }{
		{"grace", "henry", "great-uncle"}, // grandparent's brother
		{"grace", "ivy", "first cousin once removed"},
		{"hanna", "henry", "uncle"}, // parent's brother
		{"henry", "hanna", "niece"}, // sibling's daughter
		{"carl", "ivy", "niece"},
		{"ivy", "carl", "uncle"},
		{"emma", "ivy", "first cousin"},
	}
	for _, c := range cases {
		if got := describe(t, f, c.from, c.to); got != c.want {
			t.Errorf("Describe(%s,%s) = %q, want %q", c.from, c.to, got, c.want)
		}
	}
}

func TestCousinRemovedTerms(t *testing.T) {
	f := mustFamily(t, fixture)
	// grace vs ivy: da=3, db=2 -> first cousin once removed.
	if got := describe(t, f, "grace", "ivy"); got != "first cousin once removed" {
		t.Errorf("grace->ivy = %q", got)
	}
	// Symmetry: same class of term from the other side.
	if got := describe(t, f, "ivy", "grace"); got != "first cousin once removed" {
		t.Errorf("ivy->grace = %q", got)
	}
}

func TestUnrelatedPersons(t *testing.T) {
	f := mustFamily(t, fixture)
	if _, err := Describe(f, "kai", "emma"); err == nil {
		t.Error("unrelated persons expected error")
	}
}

func TestOrdinalFallback(t *testing.T) {
	if got := ordinal(1); got != "first" {
		t.Errorf("ordinal(1) = %q", got)
	}
	if got := ordinal(12); got != "12th" {
		t.Errorf("ordinal(12) = %q", got)
	}
}

func TestCousinTermDirect(t *testing.T) {
	if got := cousinTerm(2, 2); got != "first cousin" {
		t.Errorf("cousinTerm(2,2) = %q", got)
	}
	if got := cousinTerm(3, 3); got != "second cousin" {
		t.Errorf("cousinTerm(3,3) = %q", got)
	}
	if got := cousinTerm(4, 2); got != "first cousin twice removed" {
		t.Errorf("cousinTerm(4,2) = %q", got)
	}
	if got := cousinTerm(5, 2); got != "first cousin 3 times removed" {
		t.Errorf("cousinTerm(5,2) = %q", got)
	}
}
