package genealogy

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Person struct {
	Name  string
	Sex   string
	Birth int
}

type Family struct {
	people  map[string]*Person
	parents map[string][]string
	order   []string
}

type parentPair struct {
	parent string
	child  string
	lineno int
}

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
				return nil, bindParseErr(fmt.Errorf("line %d: sex must be F or M, got %q", lineno, sex))
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

func (f *Family) Person(name string) (*Person, bool) {
	p, ok := f.people[name]
	return p, ok
}

func (f *Family) Names() []string {
	names := make([]string, 0, len(f.people))
	for n := range f.people {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func (f *Family) Parents(name string) ([]string, error) {
	if _, ok := f.people[name]; !ok {
		return nil, fmt.Errorf("unknown person %q", name)
	}
	ps := append([]string(nil), f.parents[name]...)
	sort.Strings(ps)
	return ps, nil
}

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
