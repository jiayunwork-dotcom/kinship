package query

import (
	"sort"
	"strings"

	"kinship/internal/genealogy"
)

type SearchResult struct {
	Name  string
	Score float64
}

func SearchByName(f *genealogy.Family, query string) []SearchResult {
	query = strings.ToLower(query)
	var results []SearchResult
	for _, name := range f.Names() {
		lower := strings.ToLower(name)
		if strings.Contains(lower, query) {
			score := float64(len(query)) / float64(len(name))
			results = append(results, SearchResult{Name: name, Score: score})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results
}

func SearchByBirthRange(f *genealogy.Family, from, to int) []string {
	var result []string
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		if p.Birth >= from && p.Birth <= to {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

func SearchBySex(f *genealogy.Family, sex string) []string {
	var result []string
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		if p.Sex == sex {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

func SearchChildless(f *genealogy.Family) []string {
	var result []string
	for _, name := range f.Names() {
		children, _ := f.Children(name)
		if len(children) == 0 {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

func SearchWithMinChildren(f *genealogy.Family, n int) []string {
	var result []string
	for _, name := range f.Names() {
		children, _ := f.Children(name)
		if len(children) >= n {
			result = append(result, name)
		}
	}
	sort.Strings(result)
	return result
}

func OldestPerson(f *genealogy.Family) string {
	names := f.Names()
	if len(names) == 0 {
		return ""
	}
	oldest := names[0]
	oldestBirth := 9999
	for _, name := range names {
		p, _ := f.Person(name)
		if p.Birth < oldestBirth {
			oldestBirth = p.Birth
			oldest = name
		}
	}
	return oldest
}

func YoungestPerson(f *genealogy.Family) string {
	names := f.Names()
	if len(names) == 0 {
		return ""
	}
	youngest := names[0]
	youngestBirth := 0
	for _, name := range names {
		p, _ := f.Person(name)
		if p.Birth > youngestBirth {
			youngestBirth = p.Birth
			youngest = name
		}
	}
	return youngest
}
