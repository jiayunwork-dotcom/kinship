package genealogy

import (
	"fmt"
	"sort"
	"strings"
)

func Format(f *Family) string {
	var b strings.Builder
	names := f.Names()
	for _, name := range names {
		p := f.people[name]
		fmt.Fprintf(&b, "PERSON %s %s %d\n", p.Name, p.Sex, p.Birth)
	}
	b.WriteString("\n")

	type link struct{ parent, child string }
	var links []link
	for child, parents := range f.parents {
		for _, p := range parents {
			links = append(links, link{p, child})
		}
	}
	sort.Slice(links, func(i, j int) bool {
		if links[i].parent != links[j].parent {
			return links[i].parent < links[j].parent
		}
		return links[i].child < links[j].child
	})
	for _, l := range links {
		fmt.Fprintf(&b, "PARENT %s %s\n", l.parent, l.child)
	}
	return b.String()
}

func (f *Family) Size() int {
	return len(f.people)
}

func (f *Family) ParentCount() int {
	count := 0
	for _, ps := range f.parents {
		count += len(ps)
	}
	return count
}

func (f *Family) HasPerson(name string) bool {
	_, ok := f.people[name]
	return ok
}

func (f *Family) Siblings(name string) ([]string, error) {
	if _, ok := f.people[name]; !ok {
		return nil, fmt.Errorf("unknown person %q", name)
	}
	parents, _ := f.Parents(name)
	if len(parents) == 0 {
		return nil, nil
	}
	sibSet := map[string]bool{}
	for _, p := range parents {
		children, _ := f.Children(p)
		for _, c := range children {
			if c != name {
				sibSet[c] = true
			}
		}
	}
	var sibs []string
	for s := range sibSet {
		sibs = append(sibs, s)
	}
	sort.Strings(sibs)
	return sibs, nil
}

func (f *Family) BirthRange() (int, int) {
	min, max := 9999, 0
	for _, p := range f.people {
		if p.Birth < min {
			min = p.Birth
		}
		if p.Birth > max {
			max = p.Birth
		}
	}
	return min, max
}

func (f *Family) PersonsByBirth() []string {
	type item struct {
		name  string
		birth int
	}
	var items []item
	for _, p := range f.people {
		items = append(items, item{p.Name, p.Birth})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].birth != items[j].birth {
			return items[i].birth < items[j].birth
		}
		return items[i].name < items[j].name
	})
	names := make([]string, len(items))
	for i, it := range items {
		names[i] = it.name
	}
	return names
}
