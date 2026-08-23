// Package query provides advanced genealogical queries: descendant counting,
// furthest relative, common descendants, and relationship paths.
package query

import (
	"fmt"
	"sort"

	"kinship/internal/genealogy"
)

// DescendantCount returns the total number of descendants of a person
// (children, grandchildren, etc.).
func DescendantCount(f *genealogy.Family, name string) (int, error) {
	if _, ok := f.Person(name); !ok {
		return 0, fmt.Errorf("unknown person %q", name)
	}
	visited := map[string]bool{name: true}
	queue := []string{name}
	count := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		children, _ := f.Children(cur)
		for _, child := range children {
			if !visited[child] {
				visited[child] = true
				count++
				queue = append(queue, child)
			}
		}
	}
	return count, nil
}

// AllDescendants returns all descendants of a person in BFS order.
func AllDescendants(f *genealogy.Family, name string) ([]string, error) {
	if _, ok := f.Person(name); !ok {
		return nil, fmt.Errorf("unknown person %q", name)
	}
	visited := map[string]bool{name: true}
	queue := []string{name}
	var result []string
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		children, _ := f.Children(cur)
		for _, child := range children {
			if !visited[child] {
				visited[child] = true
				result = append(result, child)
				queue = append(queue, child)
			}
		}
	}
	sort.Strings(result)
	return result, nil
}

// FurthestRelative finds the person with the longest relationship path from name.
func FurthestRelative(f *genealogy.Family, name string) (string, int, error) {
	if _, ok := f.Person(name); !ok {
		return "", 0, fmt.Errorf("unknown person %q", name)
	}
	// BFS through both parent and child edges
	dist := map[string]int{name: 0}
	queue := []string{name}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		parents, _ := f.Parents(cur)
		children, _ := f.Children(cur)
		neighbors := append(parents, children...)
		for _, nb := range neighbors {
			if _, seen := dist[nb]; !seen {
				dist[nb] = dist[cur] + 1
				queue = append(queue, nb)
			}
		}
	}
	best := ""
	bestDist := 0
	for person, d := range dist {
		if person == name {
			continue
		}
		if d > bestDist || (d == bestDist && person < best) {
			best = person
			bestDist = d
		}
	}
	return best, bestDist, nil
}

// CommonDescendants returns persons who are descendants of both a and b.
func CommonDescendants(f *genealogy.Family, a, b string) ([]string, error) {
	descA, err := AllDescendants(f, a)
	if err != nil {
		return nil, err
	}
	descB, err := AllDescendants(f, b)
	if err != nil {
		return nil, err
	}
	setB := map[string]bool{}
	for _, d := range descB {
		setB[d] = true
	}
	var common []string
	for _, d := range descA {
		if setB[d] {
			common = append(common, d)
		}
	}
	sort.Strings(common)
	return common, nil
}

// RelationshipPath finds the shortest undirected path between two persons
// (traversing both parent→child and child→parent edges).
func RelationshipPath(f *genealogy.Family, from, to string) ([]string, error) {
	if _, ok := f.Person(from); !ok {
		return nil, fmt.Errorf("unknown person %q", from)
	}
	if _, ok := f.Person(to); !ok {
		return nil, fmt.Errorf("unknown person %q", to)
	}
	if from == to {
		return []string{from}, nil
	}
	// BFS
	prev := map[string]string{from: ""}
	queue := []string{from}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		parents, _ := f.Parents(cur)
		children, _ := f.Children(cur)
		for _, nb := range append(parents, children...) {
			if _, seen := prev[nb]; !seen {
				prev[nb] = cur
				if nb == to {
					return publishPath(buildPath(prev, from, to)), nil
				}
				queue = append(queue, nb)
			}
		}
	}
	return nil, fmt.Errorf("no path between %s and %s", from, to)
}

func buildPath(prev map[string]string, from, to string) []string {
	var path []string
	cur := to
	for cur != "" {
		path = append(path, cur)
		cur = prev[cur]
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// Roots returns all persons with no parents (founders).
func Roots(f *genealogy.Family) []string {
	var roots []string
	for _, name := range f.Names() {
		parents, _ := f.Parents(name)
		if len(parents) == 0 {
			roots = append(roots, name)
		}
	}
	sort.Strings(roots)
	return roots
}

// Leaves returns all persons with no children (terminal).
func Leaves(f *genealogy.Family) []string {
	var leaves []string
	for _, name := range f.Names() {
		children, _ := f.Children(name)
		if len(children) == 0 {
			leaves = append(leaves, name)
		}
	}
	sort.Strings(leaves)
	return leaves
}
