package pedigree

import (
	"fmt"
	"sort"

	"kinship/internal/genealogy"
)

type AncestorNode struct {
	Name       string
	Generation int
	Position   int
	Sex        string
	Birth      int
}

func AncestorChart(f *genealogy.Family, person string, maxGen int) ([]AncestorNode, error) {
	if _, ok := f.Person(person); !ok {
		return nil, fmt.Errorf("unknown person %q", person)
	}
	if maxGen <= 0 {
		maxGen = 5
	}

	var nodes []AncestorNode
	type entry struct {
		name string
		gen  int
		pos  int
	}
	queue := []entry{{name: person, gen: 0, pos: 0}}

	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		p, ok := f.Person(e.name)
		if !ok {
			continue
		}
		nodes = append(nodes, AncestorNode{
			Name:       e.name,
			Generation: e.gen,
			Position:   e.pos,
			Sex:        p.Sex,
			Birth:      p.Birth,
		})

		if e.gen >= maxGen {
			continue
		}
		parents, _ := f.Parents(e.name)
		for i, par := range parents {
			pp, _ := f.Person(par)
			posOffset := 0
			if pp != nil && pp.Sex == "F" {
				posOffset = 1
			} else if i == 1 {
				posOffset = 1
			}
			queue = append(queue, entry{name: par, gen: e.gen + 1, pos: e.pos*2 + posOffset})
		}
	}
	return nodes, nil
}

type DescendantNode struct {
	Name       string
	Depth      int
	ChildIndex int
	Sex        string
	Birth      int
}

func DescendantTree(f *genealogy.Family, person string, maxDepth int) ([]DescendantNode, error) {
	if _, ok := f.Person(person); !ok {
		return nil, fmt.Errorf("unknown person %q", person)
	}
	if maxDepth <= 0 {
		maxDepth = 10
	}
	var nodes []DescendantNode
	visited := map[string]bool{}
	var dfs func(name string, depth, childIdx int)
	dfs = func(name string, depth, childIdx int) {
		if visited[name] || depth > maxDepth {
			return
		}
		visited[name] = true
		p, _ := f.Person(name)
		nodes = append(nodes, DescendantNode{
			Name:       name,
			Depth:      depth,
			ChildIndex: childIdx,
			Sex:        p.Sex,
			Birth:      p.Birth,
		})
		children, _ := f.Children(name)
		sort.Strings(children)
		for i, child := range children {
			dfs(child, depth+1, i)
		}
	}
	dfs(person, 0, 0)
	return nodes, nil
}

func Completeness(f *genealogy.Family, person string, maxGen int) (map[int]float64, error) {
	chart, err := AncestorChart(f, person, maxGen)
	if err != nil {
		return nil, err
	}
	genCount := map[int]int{}
	for _, n := range chart {
		genCount[n.Generation]++
	}
	result := map[int]float64{}
	for g := 0; g <= maxGen; g++ {
		possible := 1
		if g > 0 {
			possible = 1 << g
		}
		result[g] = float64(genCount[g]) / float64(possible)
	}
	return result, nil
}

func TotalAncestorCount(f *genealogy.Family, person string) (int, error) {
	anc, err := f.Ancestors(person)
	if err != nil {
		return 0, err
	}
	return len(anc) - 1, nil
}
