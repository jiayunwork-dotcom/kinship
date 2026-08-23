// Package export provides serialization of family data to standard genealogy
// formats: simplified GEDCOM, DOT (Graphviz), and JSON.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"kinship/internal/genealogy"
)

// WriteDOT exports the family as a Graphviz DOT directed graph.
// Edges point from parent to child.
func WriteDOT(w io.Writer, f *genealogy.Family) error {
	fmt.Fprintln(w, "digraph family {")
	fmt.Fprintln(w, "  rankdir=TB;")
	fmt.Fprintln(w, "  node [shape=box];")
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		color := "pink"
		if p.Sex == "M" {
			color = "lightblue"
		}
		fmt.Fprintf(w, "  %q [label=%q, style=filled, fillcolor=%s];\n",
			name, fmt.Sprintf("%s\\n%d", name, p.Birth), color)
	}
	for _, child := range f.Names() {
		parents, _ := f.Parents(child)
		for _, p := range parents {
			fmt.Fprintf(w, "  %q -> %q;\n", p, child)
		}
	}
	fmt.Fprintln(w, "}")
	return nil
}

// jsonFamily is the JSON serialization format.
type jsonFamily struct {
	Persons []jsonPerson `json:"persons"`
	Links   []jsonLink   `json:"links"`
}
type jsonPerson struct {
	Name  string `json:"name"`
	Sex   string `json:"sex"`
	Birth int    `json:"birth"`
}
type jsonLink struct {
	Parent string `json:"parent"`
	Child  string `json:"child"`
}

// WriteJSON exports the family as structured JSON.
func WriteJSON(w io.Writer, f *genealogy.Family) error {
	jf := jsonFamily{}
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		jf.Persons = append(jf.Persons, jsonPerson{Name: p.Name, Sex: p.Sex, Birth: p.Birth})
	}
	for _, child := range f.Names() {
		parents, _ := f.Parents(child)
		for _, p := range parents {
			jf.Links = append(jf.Links, jsonLink{Parent: p, Child: child})
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jf)
}

// WriteGEDCOM exports the family in a simplified GEDCOM 5.5 format.
// This is a minimal subset suitable for import into genealogy software.
func WriteGEDCOM(w io.Writer, f *genealogy.Family) error {
	fmt.Fprintln(w, "0 HEAD")
	fmt.Fprintln(w, "1 SOUR kinship")
	fmt.Fprintln(w, "1 GEDC")
	fmt.Fprintln(w, "2 VERS 5.5.1")
	fmt.Fprintln(w, "1 CHAR UTF-8")

	names := f.Names()
	nameToID := map[string]string{}
	for i, name := range names {
		id := fmt.Sprintf("@I%d@", i+1)
		nameToID[name] = id
	}

	// individuals
	for _, name := range names {
		p, _ := f.Person(name)
		id := nameToID[name]
		fmt.Fprintf(w, "0 %s INDI\n", id)
		fmt.Fprintf(w, "1 NAME %s\n", name)
		fmt.Fprintf(w, "1 SEX %s\n", p.Sex)
		fmt.Fprintf(w, "1 BIRT\n")
		fmt.Fprintf(w, "2 DATE %d\n", p.Birth)
	}

	// families (group by parent pairs)
	type famKey struct{ p1, p2 string }
	families := map[famKey][]string{}
	for _, child := range names {
		parents, _ := f.Parents(child)
		var key famKey
		sorted := make([]string, len(parents))
		copy(sorted, parents)
		sort.Strings(sorted)
		if len(sorted) >= 1 {
			key.p1 = sorted[0]
		}
		if len(sorted) >= 2 {
			key.p2 = sorted[1]
		}
		families[key] = append(families[key], child)
	}

	famID := 1
	for key, children := range families {
		fid := fmt.Sprintf("@F%d@", famID)
		famID++
		fmt.Fprintf(w, "0 %s FAM\n", fid)
		if key.p1 != "" {
			p1, _ := f.Person(key.p1)
			if p1.Sex == "M" {
				fmt.Fprintf(w, "1 HUSB %s\n", nameToID[key.p1])
			} else {
				fmt.Fprintf(w, "1 WIFE %s\n", nameToID[key.p1])
			}
		}
		if key.p2 != "" {
			p2, _ := f.Person(key.p2)
			if p2.Sex == "M" {
				fmt.Fprintf(w, "1 HUSB %s\n", nameToID[key.p2])
			} else {
				fmt.Fprintf(w, "1 WIFE %s\n", nameToID[key.p2])
			}
		}
		for _, child := range children {
			fmt.Fprintf(w, "1 CHIL %s\n", nameToID[child])
		}
	}

	fmt.Fprintln(w, "0 TRLR")
	return nil
}
