package relation

import (
	"sort"

	"kinship/internal/genealogy"
)

type RelEntry struct {
	From string
	To   string
	Term string
}

func BatchDescribe(f *genealogy.Family) []RelEntry {
	names := f.Names()
	var entries []RelEntry
	for _, a := range names {
		for _, b := range names {
			if a == b {
				continue
			}
			term, err := Describe(f, a, b)
			if err != nil {
				continue
			}
			entries = append(entries, RelEntry{From: a, To: b, Term: term})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].From != entries[j].From {
			return entries[i].From < entries[j].From
		}
		return entries[i].To < entries[j].To
	})
	return entries
}

func RelativesOf(f *genealogy.Family, name string) []RelEntry {
	names := f.Names()
	var entries []RelEntry
	for _, other := range names {
		if other == name {
			continue
		}
		term, err := Describe(f, name, other)
		if err != nil {
			continue
		}
		entries = append(entries, RelEntry{From: name, To: other, Term: term})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].To < entries[j].To
	})
	return entries
}

func FindByTerm(f *genealogy.Family, term string) []RelEntry {
	names := f.Names()
	var entries []RelEntry
	for _, a := range names {
		for _, b := range names {
			if a == b {
				continue
			}
			t, err := Describe(f, a, b)
			if err != nil {
				continue
			}
			if t == term {
				entries = append(entries, RelEntry{From: a, To: b, Term: t})
			}
		}
	}
	return entries
}

func DistinctTerms(f *genealogy.Family) []string {
	seen := map[string]bool{}
	names := f.Names()
	for _, a := range names {
		for _, b := range names {
			if a == b {
				continue
			}
			term, err := Describe(f, a, b)
			if err != nil {
				continue
			}
			seen[term] = true
		}
	}
	var terms []string
	for t := range seen {
		terms = append(terms, t)
	}
	sort.Strings(terms)
	return terms
}
