package validate

import (
	"fmt"
	"sort"

	"kinship/internal/genealogy"
)

type Issue struct {
	Severity string
	Message  string
	Persons  []string
}

type Result struct {
	Issues []Issue
}

func (r *Result) OK() bool {
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			return false
		}
	}
	return true
}

func (r *Result) Errors() []Issue {
	var out []Issue
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			out = append(out, iss)
		}
	}
	return out
}

func Validate(f *genealogy.Family) *Result {
	r := &Result{}
	checkCycles(f, r)
	checkBirthOrder(f, r)
	checkMaxParents(f, r)
	checkOrphans(f, r)
	checkGenerationBound(f, r)
	return r
}

func checkCycles(f *genealogy.Family, r *Result) {
	for _, name := range f.Names() {
		visited := map[string]bool{}
		if hasCycle(f, name, visited) {
			r.Issues = append(r.Issues, Issue{
				Severity: "error",
				Message:  fmt.Sprintf("ancestry cycle detected involving %q", name),
				Persons:  []string{name},
			})
		}
	}
}

func hasCycle(f *genealogy.Family, start string, visited map[string]bool) bool {
	stack := []string{start}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		parents, _ := f.Parents(cur)
		for _, p := range parents {
			if p == start {
				return true
			}
			if !visited[p] {
				visited[p] = true
				stack = append(stack, p)
			}
		}
	}
	return false
}

func checkBirthOrder(f *genealogy.Family, r *Result) {
	for _, child := range f.Names() {
		childP, _ := f.Person(child)
		parents, _ := f.Parents(child)
		for _, pName := range parents {
			parent, _ := f.Person(pName)
			if parent.Birth > childP.Birth {
				r.Issues = append(r.Issues, Issue{
					Severity: "warning",
					Message:  fmt.Sprintf("parent %q (born %d) born after child %q (born %d)", pName, parent.Birth, child, childP.Birth),
					Persons:  []string{pName, child},
				})
			}
			if childP.Birth-parent.Birth < 12 {
				r.Issues = append(r.Issues, Issue{
					Severity: "warning",
					Message:  fmt.Sprintf("parent %q only %d years older than child %q (biologically unlikely)", pName, childP.Birth-parent.Birth, child),
					Persons:  []string{pName, child},
				})
			}
		}
	}
}

func checkMaxParents(f *genealogy.Family, r *Result) {
	for _, name := range f.Names() {
		parents, _ := f.Parents(name)
		if len(parents) > 2 {
			r.Issues = append(r.Issues, Issue{
				Severity: "warning",
				Message:  fmt.Sprintf("%q has %d registered parents (expected at most 2)", name, len(parents)),
				Persons:  []string{name},
			})
		}
	}
}

func checkOrphans(f *genealogy.Family, r *Result) {
	for _, name := range f.Names() {
		parents, _ := f.Parents(name)
		children, _ := f.Children(name)
		if len(parents) == 0 && len(children) == 0 && len(f.Names()) > 1 {
			r.Issues = append(r.Issues, Issue{
				Severity: "warning",
				Message:  fmt.Sprintf("%q is isolated (no parents, no children)", name),
				Persons:  []string{name},
			})
		}
	}
}

func checkGenerationBound(f *genealogy.Family, r *Result) {
	for _, name := range f.Names() {
		anc, _ := f.Ancestors(name)
		maxDist := 0
		for _, d := range anc {
			if d > maxDist {
				maxDist = d
			}
		}
		if maxDist > 20 {
			r.Issues = append(r.Issues, Issue{
				Severity: "warning",
				Message:  fmt.Sprintf("%q has ancestry depth %d (unusually deep)", name, maxDist),
				Persons:  []string{name},
			})
		}
	}
}

func UnconnectedGroups(f *genealogy.Family) [][]string {
	names := f.Names()
	adj := map[string]map[string]bool{}
	for _, n := range names {
		adj[n] = map[string]bool{}
	}
	for _, name := range names {
		parents, _ := f.Parents(name)
		for _, p := range parents {
			adj[name][p] = true
			adj[p][name] = true
		}
	}

	visited := map[string]bool{}
	var groups [][]string
	for _, name := range names {
		if visited[name] {
			continue
		}
		var group []string
		queue := []string{name}
		visited[name] = true
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			group = append(group, cur)
			for nb := range adj[cur] {
				if !visited[nb] {
					visited[nb] = true
					queue = append(queue, nb)
				}
			}
		}
		sort.Strings(group)
		groups = append(groups, group)
	}
	return groups
}
