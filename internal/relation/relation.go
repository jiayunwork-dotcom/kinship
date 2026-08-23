// Package relation turns genealogical distances into human kinship
// terms: mother, great-grandfather, half-sibling, aunt, niece, first
// cousin once removed, and so on.
package relation

import (
	"fmt"
	"strings"

	"kinship/internal/genealogy"
)

// Describe returns the kinship term for `to` from the point of view of
// `from` (e.g. Describe(f, "grace", "emma") == "mother").
func Describe(f *genealogy.Family, from, to string) (string, error) {
	if _, ok := f.Person(from); !ok {
		return "", fmt.Errorf("unknown person %q", from)
	}
	if _, ok := f.Person(to); !ok {
		return "", fmt.Errorf("unknown person %q", to)
	}
	if from == to {
		return "self", nil
	}
	ancFrom, err := f.Ancestors(from)
	if err != nil {
		return "", err
	}
	if d, ok := ancFrom[to]; ok {
		return AncestorTerm(f, to, d), nil
	}
	ancTo, err := f.Ancestors(to)
	if err != nil {
		return "", err
	}
	if d, ok := ancTo[from]; ok {
		return descendantTerm(f, to, d), nil
	}
	_, da, db, err := f.LCA(from, to)
	if err != nil {
		return "", err
	}
	switch {
	case da == 1 && db == 1:
		shared, err := f.SharedParents(from, to)
		if err != nil {
			return "", err
		}
		if shared >= 2 {
			return "sibling", nil
		}
		return "half-sibling", nil
	case da == 1 && db >= 2:
		// `to` descends from a sibling of `from`.
		return collateralTerm(f, to, db, "niece", "nephew"), nil
	case da >= 2 && db == 1:
		// `to` is a sibling of an ancestor of `from`.
		return collateralTerm(f, to, da, "aunt", "uncle"), nil
	default:
		return cousinTerm(da, db), nil
	}
}

func sexOf(f *genealogy.Family, name string) string {
	if p, ok := f.Person(name); ok {
		return p.Sex
	}
	return "F"
}

// AncestorTerm names an ancestor `d` generations up and publishes the
// term through the report ledger.
func AncestorTerm(f *genealogy.Family, name string, d int) string {
	return publishTerm(ancestorTerm(f, name, d))
}

// ancestorTerm names an ancestor `d` generations up: mother,
// grandmother, great-grandmother, ...
func ancestorTerm(f *genealogy.Family, name string, d int) string {
	fem, masc := "mother", "father"
	if d >= 2 {
		fem, masc = "grandmother", "grandfather"
	}
	term := fem
	if sexOf(f, name) == "M" {
		term = masc
	}
	return strings.Repeat("great-", max(0, d-2)) + term
}

// descendantTerm names a descendant `d` generations down: daughter,
// granddaughter, great-granddaughter, ...
func descendantTerm(f *genealogy.Family, name string, d int) string {
	fem, masc := "daughter", "son"
	if d >= 2 {
		fem, masc = "granddaughter", "grandson"
	}
	term := fem
	if sexOf(f, name) == "M" {
		term = masc
	}
	return strings.Repeat("great-", max(0, d-2)) + term
}

// collateralTerm names a niece/nephew or aunt/uncle `d` generations
// away from the shared sibling line.
func collateralTerm(f *genealogy.Family, name string, d int, fem, masc string) string {
	term := fem
	if sexOf(f, name) == "M" {
		term = masc
	}
	return strings.Repeat("great-", max(0, d-2)) + term
}

// cousinTerm builds the canonical "Nth cousin M times removed" phrase.
// da and db are the distances from each person to their lowest common
// ancestor; cousins share an ancestor at equal distance, each removed
// step is one generation of separation.
func cousinTerm(da, db int) string {
	n := min(da, db) - 1
	r := da - db
	if r < 0 {
		r = -r
	}
	term := ordinal(n) + " cousin"
	if r == 0 {
		return term
	}
	if r == 1 {
		return term + " once removed"
	}
	if r == 2 {
		return term + " twice removed"
	}
	return fmt.Sprintf("%s %d times removed", term, r)
}

// ordinal renders 1..10 as words and falls back to "11th" style.
func ordinal(n int) string {
	words := map[int]string{
		1: "first", 2: "second", 3: "third", 4: "fourth", 5: "fifth",
		6: "sixth", 7: "seventh", 8: "eighth", 9: "ninth", 10: "tenth",
	}
	if w, ok := words[n]; ok {
		return w
	}
	return fmt.Sprintf("%dth", n)
}
