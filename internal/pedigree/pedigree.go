// Package pedigree provides structured pedigree chart construction for
// ancestor fan charts and descendant trees. It computes the layout data
// needed to render visual pedigree representations.
package pedigree

import (
	"fmt"
	"sort"

	"kinship/internal/genealogy"
)

// AncestorNode represents one cell in an ancestor fan chart.
type AncestorNode struct {
	Name       string
	Generation int
	Position   int // position within the generation (0-based)
	Sex        string
	Birth      int
}

// AncestorChart builds a structured ancestor chart for a person.
// Generation 0 is the person themselves, gen 1 is parents, gen 2 grandparents, etc.
// Each generation has up to 2^gen positions (some may be empty if unknown).
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
		// position: father at 2*pos, mother at 2*pos+1
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

// DescendantNode represents one cell in a descendant tree.
type DescendantNode struct {
	Name       string
	Depth      int
	ChildIndex int
	Sex        string
	Birth      int
}

// DescendantTree builds a depth-first descendant tree from a person.
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

// Completeness computes the ancestor completeness at each generation.
// Returns gen -> (known ancestors / possible ancestors at that gen).
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
			possible = 1 << g // 2^g
		}
		result[g] = float64(genCount[g]) / float64(possible)
	}
	return result, nil
}

// TotalAncestorCount returns how many distinct ancestors are known.
func TotalAncestorCount(f *genealogy.Family, person string) (int, error) {
	anc, err := f.Ancestors(person)
	if err != nil {
		return 0, err
	}
	return len(anc) - 1, nil // exclude self
}
