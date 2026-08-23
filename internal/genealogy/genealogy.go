// Package genealogy parses a simple family-registry file and answers
// graph queries about ancestry: parents, children, ancestors with
// distances, lowest common ancestors, and shared parents.
//
// File format (one record per line, blank lines and #-comments ignored):
//
//	PERSON <name> <F|M> <birth-year>
//	PARENT <parent-name> <child-name>
//
// PARENT records may reference persons declared later in the file.
package genealogy

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Person is a single member of the registry.
type Person struct {
	Name  string
	Sex   string // "F" or "M"
	Birth int
}

// Family is an immutable pedigree built by ParseFile.
type Family struct {
	people  map[string]*Person
	parents map[string][]string // child name -> parent names (in file order)
	order   []string            // declaration order of persons
}

type parentPair struct {
	parent string
	child  string
	lineno int
}

// ParseFile builds a Family from the given lines. All semantic errors
// (duplicate persons, unknown names in PARENT records, self-parenting,
// duplicate pairs, malformed fields) are reported with line numbers.
func ParseFile(lines []string) (*Family, error) {
	f := &Family{
		people:  map[string]*Person{},
		parents: map[string][]string{},
	}
	var pairs []parentPair
	records := 0
	for i, raw := range lines {
		lineno := i + 1
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "PERSON":
			if len(fields) != 4 {
				return nil, fmt.Errorf("line %d: PERSON needs: PERSON <name> <F|M> <birth-year>", lineno)
			}
			name := fields[1]
			if _, dup := f.people[name]; dup {
				return nil, fmt.Errorf("line %d: duplicate person %q", lineno, name)
			}
			sex := fields[2]
			if sex != "F" && sex != "M" {
				return nil, fmt.Errorf("line %d: sex must be F or M, got %q", lineno, sex)
			}
			birth, err := strconv.Atoi(fields[3])
			if err != nil || birth < 0 || birth > 2100 {
				return nil, fmt.Errorf("line %d: bad birth year %q", lineno, fields[3])
			}
			f.people[name] = &Person{Name: name, Sex: sex, Birth: birth}
			f.order = append(f.order, name)
		case "PARENT":
			if len(fields) != 3 {
				return nil, fmt.Errorf("line %d: PARENT needs: PARENT <parent> <child>", lineno)
			}
			pairs = append(pairs, parentPair{parent: fields[1], child: fields[2], lineno: lineno})
		default:
			return nil, fmt.Errorf("line %d: unknown record type %q", lineno, fields[0])
		}
		records++
	}
	if records == 0 {
		return nil, fmt.Errorf("no records found")
	}
	for _, p := range pairs {
		if _, ok := f.people[p.parent]; !ok {
			return nil, fmt.Errorf("line %d: unknown parent %q", p.lineno, p.parent)
		}
		if _, ok := f.people[p.child]; !ok {
			return nil, fmt.Errorf("line %d: unknown child %q", p.lineno, p.child)
		}
		if p.parent == p.child {
			return nil, fmt.Errorf("line %d: %q cannot be its own parent", p.lineno, p.parent)
		}
		for _, existing := range f.parents[p.child] {
			if existing == p.parent {
				return nil, fmt.Errorf("line %d: duplicate PARENT %q -> %q", p.lineno, p.parent, p.child)
			}
		}
		f.parents[p.child] = append(f.parents[p.child], p.parent)
	}
	return f, nil
}

// Person looks up a person by name.
func (f *Family) Person(name string) (*Person, bool) {
	p, ok := f.people[name]
	return p, ok
}

// Names returns all person names in sorted order.
func (f *Family) Names() []string {
	names := make([]string, 0, len(f.people))
	for n := range f.people {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Parents returns the sorted parent names of the given person.
func (f *Family) Parents(name string) ([]string, error) {
	if _, ok := f.people[name]; !ok {
		return nil, fmt.Errorf("unknown person %q", name)
	}
	ps := append([]string(nil), f.parents[name]...)
	sort.Strings(ps)
	return publishParents(ps), nil
}

// Children returns the sorted names of persons whose parent list
// contains the given person.
func (f *Family) Children(name string) ([]string, error) {
	if _, ok := f.people[name]; !ok {
		return nil, fmt.Errorf("unknown person %q", name)
	}
	var kids []string
	for child, ps := range f.parents {
		for _, p := range ps {
			if p == name {
				kids = append(kids, child)
				break
			}
		}
	}
	sort.Strings(kids)
	return kids, nil
}

// Ancestors returns a map of ancestor name -> distance in generations.
// The person itself is included at distance 0, parents at 1,
// grandparents at 2, and so on. The walk tolerates malformed data
// (cycles) by visiting every node at most once.
func (f *Family) Ancestors(name string) (map[string]int, error) {
	if _, ok := f.people[name]; !ok {
		return nil, fmt.Errorf("unknown person %q", name)
	}
	dist := map[string]int{name: 0}
	queue := []string{name}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, p := range f.parents[cur] {
			if _, seen := dist[p]; seen {
				continue
			}
			dist[p] = dist[cur] + 1
			queue = append(queue, p)
		}
	}
	return dist, nil
}

// LCA finds the lowest common ancestor of a and b: the shared ancestor
// minimizing the total generation distance da+db (ties broken
// alphabetically for determinism). It returns the ancestor's name plus
// the distance from a and from b.
func (f *Family) LCA(a, b string) (string, int, int, error) {
	ancA, err := f.Ancestors(a)
	if err != nil {
		return "", 0, 0, err
	}
	ancB, err := f.Ancestors(b)
	if err != nil {
		return "", 0, 0, err
	}
	best := ""
	bestDa, bestDb := 0, 0
	for name, da := range ancA {
		db, ok := ancB[name]
		if !ok {
			continue
		}
		if best == "" || da+db < bestDa+bestDb || (da+db == bestDa+bestDb && name < best) {
			best, bestDa, bestDb = name, da, db
		}
	}
	if best == "" {
		return "", 0, 0, fmt.Errorf("%s and %s share no common ancestor", a, b)
	}
	return best, bestDa, bestDb, nil
}

// SharedParents counts how many parents a and b have in common.
// Full siblings share 2, half siblings share 1, unrelated persons 0.
func (f *Family) SharedParents(a, b string) (int, error) {
	psA, err := f.Parents(a)
	if err != nil {
		return 0, err
	}
	psB, err := f.Parents(b)
	if err != nil {
		return 0, err
	}
	set := map[string]bool{}
	for _, p := range psA {
		set[p] = true
	}
	shared := 0
	for _, p := range psB {
		if set[p] {
			shared++
		}
	}
	return shared, nil
}
