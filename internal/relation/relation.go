package relation

import (
	"fmt"
	"strings"

	"kinship/internal/genealogy"
)

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
		return ancestorTerm(f, to, d), nil
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
		return collateralTerm(f, to, db, "niece", "nephew"), nil
	case da >= 2 && db == 1:
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

func collateralTerm(f *genealogy.Family, name string, d int, fem, masc string) string {
	term := fem
	if sexOf(f, name) == "M" {
		term = masc
	}
	return strings.Repeat("great-", max(0, d-2)) + term
}

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
