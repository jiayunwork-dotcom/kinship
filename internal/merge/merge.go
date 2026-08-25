package merge

import (
	"fmt"
	"sort"
	"strings"

	"kinship/internal/genealogy"
)

type Conflict struct {
	Person string
	Field  string
	ValueA string
	ValueB string
}

type MergeResult struct {
	Lines     []string
	Conflicts []Conflict
}

func (r *MergeResult) HasConflicts() bool {
	return len(r.Conflicts) > 0
}

func Merge(linesA, linesB []string) (*MergeResult, error) {
	fA, err := genealogy.ParseFile(linesA)
	if err != nil {
		return nil, fmt.Errorf("merge: parse A: %w", err)
	}
	fB, err := genealogy.ParseFile(linesB)
	if err != nil {
		return nil, fmt.Errorf("merge: parse B: %w", err)
	}

	result := &MergeResult{}
	persons := map[string]bool{}
	var personLines []string
	var parentLines []string

	for _, name := range fA.Names() {
		p, _ := fA.Person(name)
		personLines = append(personLines, fmt.Sprintf("PERSON %s %s %d", p.Name, p.Sex, p.Birth))
		persons[name] = true
	}

	for _, name := range fB.Names() {
		pB, _ := fB.Person(name)
		if persons[name] {
			pA, _ := fA.Person(name)
			if pA.Sex != pB.Sex {
				result.Conflicts = append(result.Conflicts, Conflict{
					Person: name, Field: "sex", ValueA: pA.Sex, ValueB: pB.Sex,
				})
			}
			if pA.Birth != pB.Birth {
				result.Conflicts = append(result.Conflicts, Conflict{
					Person: name, Field: "birth",
					ValueA: fmt.Sprintf("%d", pA.Birth),
					ValueB: fmt.Sprintf("%d", pB.Birth),
				})
			}
			continue
		}
		personLines = append(personLines, fmt.Sprintf("PERSON %s %s %d", pB.Name, pB.Sex, pB.Birth))
		persons[name] = true
	}

	type parentChild struct{ parent, child string }
	seen := map[parentChild]bool{}

	collectParents := func(f *genealogy.Family) {
		for _, child := range f.Names() {
			parents, _ := f.Parents(child)
			for _, p := range parents {
				pc := parentChild{p, child}
				if !seen[pc] && persons[p] && persons[child] {
					seen[pc] = true
					parentLines = append(parentLines, fmt.Sprintf("PARENT %s %s", p, child))
				}
			}
		}
	}
	collectParents(fA)
	collectParents(fB)

	sort.Strings(personLines)
	sort.Strings(parentLines)
	result.Lines = append(personLines, parentLines...)
	return result, nil
}

func MergeFiles(contentA, contentB string) (*MergeResult, error) {
	linesA := strings.Split(contentA, "\n")
	linesB := strings.Split(contentB, "\n")
	return Merge(linesA, linesB)
}

func (r *MergeResult) ConflictSummary() string {
	if !r.HasConflicts() {
		return "no conflicts"
	}
	var parts []string
	for _, c := range r.Conflicts {
		parts = append(parts, fmt.Sprintf("%s.%s: %q vs %q", c.Person, c.Field, c.ValueA, c.ValueB))
	}
	return strings.Join(parts, "; ")
}
