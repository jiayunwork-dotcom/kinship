package validate

import (
	"fmt"

	"kinship/internal/genealogy"
)

func TemporalCheck(f *genealogy.Family) []Issue {
	var issues []Issue
	issues = append(issues, checkReasonableAge(f)...)
	issues = append(issues, checkParentAge(f)...)
	return issues
}

func checkReasonableAge(f *genealogy.Family) []Issue {
	var issues []Issue
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		if p.Birth < 1500 {
			issues = append(issues, Issue{
				Severity: "warning",
				Message:  fmt.Sprintf("%q has birth year %d (before 1500, possibly unreliable)", name, p.Birth),
				Persons:  []string{name},
			})
		}
	}
	return issues
}

func checkParentAge(f *genealogy.Family) []Issue {
	var issues []Issue
	for _, child := range f.Names() {
		cp, _ := f.Person(child)
		parents, _ := f.Parents(child)
		for _, pName := range parents {
			pp, _ := f.Person(pName)
			age := cp.Birth - pp.Birth
			if age > 70 {
				issues = append(issues, Issue{
					Severity: "warning",
					Message:  fmt.Sprintf("parent %q was %d years old at birth of %q (unusually old)", pName, age, child),
					Persons:  []string{pName, child},
				})
			}
		}
	}
	return issues
}

func ConsistencyScore(f *genealogy.Family) float64 {
	names := f.Names()
	if len(names) == 0 {
		return 1.0
	}
	issues := 0
	total := 0
	for _, child := range names {
		cp, _ := f.Person(child)
		parents, _ := f.Parents(child)
		for _, pName := range parents {
			total++
			pp, _ := f.Person(pName)
			age := cp.Birth - pp.Birth
			if age < 12 || age > 70 || age < 0 {
				issues++
			}
		}
	}
	if total == 0 {
		return 1.0
	}
	return 1.0 - float64(issues)/float64(total)
}
