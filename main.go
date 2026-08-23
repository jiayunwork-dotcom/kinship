// Command kinship answers questions about a family registry file:
// who someone's ancestors are, who their children are, and what kinship
// term relates two people (mother, great-uncle, first cousin once
// removed, ...).
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"kinship/internal/genealogy"
	"kinship/internal/relation"
)

const usageText = `usage: kinship <command> [flags] <family.txt> [args]

commands:
  ancestors <family.txt> <name>   list ancestors grouped by generation
  children   <family.txt> <name>   list children of a person
  kin        <family.txt> <a> <b>  kinship term for b from a's viewpoint
  list       <family.txt>          list everyone in the registry

registry format:
  PERSON <name> <F|M> <birth-year>
  PARENT <parent> <child>
`

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, usageText)
	os.Exit(2)
}

func loadFamily(path string) (*genealogy.Family, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return genealogy.ParseFile(strings.Split(string(data), "\n"))
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]
	args := reorderFlags(os.Args[2:])
	switch cmd {
	case "ancestors":
		fs := flag.NewFlagSet("ancestors", flag.ExitOnError)
		if err := fs.Parse(args); err != nil {
			fail(err)
		}
		if fs.NArg() != 2 {
			usage()
		}
		cmdAncestors(fs.Arg(0), fs.Arg(1))
	case "children":
		fs := flag.NewFlagSet("children", flag.ExitOnError)
		if err := fs.Parse(args); err != nil {
			fail(err)
		}
		if fs.NArg() != 2 {
			usage()
		}
		cmdChildren(fs.Arg(0), fs.Arg(1))
	case "kin":
		fs := flag.NewFlagSet("kin", flag.ExitOnError)
		if err := fs.Parse(args); err != nil {
			fail(err)
		}
		if fs.NArg() != 3 {
			usage()
		}
		cmdKin(fs.Arg(0), fs.Arg(1), fs.Arg(2))
	case "list":
		fs := flag.NewFlagSet("list", flag.ExitOnError)
		if err := fs.Parse(args); err != nil {
			fail(err)
		}
		if fs.NArg() != 1 {
			usage()
		}
		cmdList(fs.Arg(0))
	case "help", "-h", "--help":
		fmt.Print(usageText)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n", cmd)
		usage()
	}
}

func cmdAncestors(path, name string) {
	f, err := loadFamily(path)
	fail(err)
	anc, err := f.Ancestors(name)
	fail(err)
	byDist := map[int][]string{}
	maxDist := 0
	for who, d := range anc {
		if d == 0 {
			continue
		}
		byDist[d] = append(byDist[d], who)
		if d > maxDist {
			maxDist = d
		}
	}
	labels := map[int]string{1: "parents", 2: "grandparents"}
	for d := 1; d <= maxDist; d++ {
		names := byDist[d]
		if len(names) == 0 {
			continue
		}
		sort.Strings(names)
		label := labels[d]
		if label == "" {
			var b strings.Builder
			for i := 0; i < d-2; i++ {
				b.WriteString("great-")
			}
			b.WriteString("grandparents")
			label = b.String()
		}
		fmt.Printf("%s: %s\n", label, strings.Join(names, ", "))
	}
}

func cmdChildren(path, name string) {
	f, err := loadFamily(path)
	fail(err)
	kids, err := f.Children(name)
	fail(err)
	if len(kids) == 0 {
		fmt.Printf("%s has no registered children\n", name)
		return
	}
	fmt.Printf("children of %s: %s\n", name, strings.Join(kids, ", "))
}

func cmdKin(path, from, to string) {
	f, err := loadFamily(path)
	fail(err)
	term, err := relation.Describe(f, from, to)
	fail(err)
	fmt.Printf("%s -> %s: %s\n", from, to, term)
}

func cmdList(path string) {
	f, err := loadFamily(path)
	fail(err)
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		fmt.Printf("%-12s %s  born %d\n", name, p.Sex, p.Birth)
	}
}

// reorderFlags moves flags (with their values) in front of positional
// arguments so that Go's flag package sees them even when the user
// writes "cmd file -flag value".
func reorderFlags(args []string) []string {
	var flags, pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") && a != "-" && a != "--" {
			flags = append(flags, a)
			if !strings.Contains(a, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags = append(flags, args[i+1])
				i++
			}
		} else {
			pos = append(pos, a)
		}
	}
	return append(flags, pos...)
}
