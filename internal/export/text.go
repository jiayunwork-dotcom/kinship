package export

import (
	"fmt"
	"io"
	"strings"

	"kinship/internal/genealogy"
)

// WriteText exports the family in a human-readable indented text format,
// showing the descendant tree from each root.
func WriteText(w io.Writer, f *genealogy.Family) error {
	// find roots (persons with no parents)
	var roots []string
	for _, name := range f.Names() {
		parents, _ := f.Parents(name)
		if len(parents) == 0 {
			roots = append(roots, name)
		}
	}

	visited := map[string]bool{}
	for _, root := range roots {
		writeDescendants(w, f, root, 0, visited)
	}
	return nil
}

func writeDescendants(w io.Writer, f *genealogy.Family, name string, depth int, visited map[string]bool) {
	if visited[name] {
		return
	}
	visited[name] = true
	p, _ := f.Person(name)
	indent := strings.Repeat("  ", depth)
	sexMark := "♀"
	if p.Sex == "M" {
		sexMark = "♂"
	}
	fmt.Fprintf(w, "%s%s %s (%d)\n", indent, sexMark, name, p.Birth)
	children, _ := f.Children(name)
	for _, child := range children {
		writeDescendants(w, f, child, depth+1, visited)
	}
}

// WriteTable exports the family as a tabular format.
func WriteTable(w io.Writer, f *genealogy.Family) error {
	fmt.Fprintf(w, "%-15s %-4s %-6s %-20s %-20s\n", "Name", "Sex", "Birth", "Parents", "Children")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 70))
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		parents, _ := f.Parents(name)
		children, _ := f.Children(name)
		fmt.Fprintf(w, "%-15s %-4s %-6d %-20s %-20s\n",
			name, p.Sex, p.Birth,
			strings.Join(parents, ","),
			strings.Join(children, ","))
	}
	return nil
}

// WriteStats exports a statistical summary.
func WriteStats(w io.Writer, f *genealogy.Family) error {
	names := f.Names()
	males := 0
	for _, name := range names {
		p, _ := f.Person(name)
		if p.Sex == "M" {
			males++
		}
	}
	var roots, leaves int
	for _, name := range names {
		parents, _ := f.Parents(name)
		children, _ := f.Children(name)
		if len(parents) == 0 {
			roots++
		}
		if len(children) == 0 {
			leaves++
		}
	}
	fmt.Fprintf(w, "Total: %d persons (%d male, %d female)\n", len(names), males, len(names)-males)
	fmt.Fprintf(w, "Roots: %d, Leaves: %d\n", roots, leaves)
	return nil
}
