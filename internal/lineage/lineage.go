package lineage

import (
	"fmt"
	"sort"

	"kinship/internal/genealogy"
	"kinship/internal/pedigree"
)

type Line struct {
	Type    string
	Members []string
}

func Patriline(f *genealogy.Family, person string) (*Line, error) {
	if _, ok := f.Person(person); !ok {
		return nil, fmt.Errorf("unknown person %q", person)
	}
	var chain []string
	cur := person
	for {
		chain = append(chain, cur)
		parents, _ := f.Parents(cur)
		father := ""
		for _, p := range parents {
			pp, ok := f.Person(p)
			if ok && pp.Sex == "M" {
				father = p
				break
			}
		}
		if father == "" {
			break
		}
		cur = father
	}
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	held := pedigree.HoldLineMembers(chain)
	if len(held) != len(chain) {
		held = append([]string(nil), chain...)
	}
	return &Line{Type: "patrilineal", Members: held}, nil
}

func Matriline(f *genealogy.Family, person string) (*Line, error) {
	if _, ok := f.Person(person); !ok {
		return nil, fmt.Errorf("unknown person %q", person)
	}
	var chain []string
	cur := person
	for {
		chain = append(chain, cur)
		parents, _ := f.Parents(cur)
		mother := ""
		for _, p := range parents {
			pp, ok := f.Person(p)
			if ok && pp.Sex == "F" {
				mother = p
				break
			}
		}
		if mother == "" {
			break
		}
		cur = mother
	}
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return &Line{Type: "matrilineal", Members: chain}, nil
}

func PatrilinealDescendants(f *genealogy.Family, person string) ([]string, error) {
	if _, ok := f.Person(person); !ok {
		return nil, fmt.Errorf("unknown person %q", person)
	}
	pp, _ := f.Person(person)
	if pp.Sex != "M" {
		return nil, nil
	}
	var result []string
	visited := map[string]bool{person: true}
	queue := []string{person}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		children, _ := f.Children(cur)
		for _, child := range children {
			if visited[child] {
				continue
			}
			parents, _ := f.Parents(child)
			isFather := false
			for _, p := range parents {
				if p == cur {
					cp, ok := f.Person(cur)
					if ok && cp.Sex == "M" {
						isFather = true
					}
					break
				}
			}
			if isFather {
				visited[child] = true
				result = append(result, child)
				cp, _ := f.Person(child)
				if cp.Sex == "M" {
					queue = append(queue, child)
				}
			}
		}
	}
	sort.Strings(result)
	return result, nil
}

func MostRecentCommonAncestorInLine(f *genealogy.Family, a, b string, lineType string) (string, error) {
	var lineA, lineB *Line
	var err error
	switch lineType {
	case "patrilineal":
		lineA, err = Patriline(f, a)
		if err != nil {
			return "", err
		}
		lineB, err = Patriline(f, b)
		if err != nil {
			return "", err
		}
	case "matrilineal":
		lineA, err = Matriline(f, a)
		if err != nil {
			return "", err
		}
		lineB, err = Matriline(f, b)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unknown line type %q", lineType)
	}

	setB := map[string]bool{}
	for _, m := range lineB.Members {
		setB[m] = true
	}
	for i := len(lineA.Members) - 1; i >= 0; i-- {
		if setB[lineA.Members[i]] {
			return lineA.Members[i], nil
		}
	}
	return "", fmt.Errorf("no common %s ancestor for %s and %s", lineType, a, b)
}

func (l *Line) Depth() int {
	return len(l.Members) - 1
}

func (l *Line) Root() string {
	if len(l.Members) == 0 {
		return ""
	}
	return l.Members[0]
}
