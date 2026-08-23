package query

import (
	"strings"
	"testing"

	"kinship/internal/genealogy"
)

func family() *genealogy.Family {
	lines := strings.Split(`PERSON alice F 1950
PERSON bob M 1948
PERSON carol F 1975
PERSON dave M 2000
PARENT alice carol
PARENT bob carol
PARENT carol dave`, "\n")
	f, _ := genealogy.ParseFile(lines)
	return f
}

func TestDescendantCount(t *testing.T) {
	c, _ := DescendantCount(family(), "alice")
	if c != 2 {
		t.Errorf("descendants of alice = %d, want 2 (carol, dave)", c)
	}
}

func TestAllDescendants(t *testing.T) {
	desc, _ := AllDescendants(family(), "alice")
	if len(desc) != 2 {
		t.Errorf("descendants = %v", desc)
	}
}

func TestFurthestRelative(t *testing.T) {
	name, dist, _ := FurthestRelative(family(), "alice")
	if dist != 2 {
		t.Errorf("furthest dist = %d, want 2", dist)
	}
	// bob and dave are both at distance 2; alphabetically bob < dave
	if name != "bob" && name != "dave" {
		t.Errorf("furthest = %s, want bob or dave", name)
	}
}

func TestRelationshipPath(t *testing.T) {
	path, _ := RelationshipPath(family(), "bob", "dave")
	if len(path) != 3 {
		t.Errorf("path = %v, want [bob carol dave]", path)
	}
}

func TestRoots(t *testing.T) {
	roots := Roots(family())
	if len(roots) != 2 {
		t.Errorf("roots = %v", roots)
	}
}

func TestLeaves(t *testing.T) {
	leaves := Leaves(family())
	if len(leaves) != 1 || leaves[0] != "dave" {
		t.Errorf("leaves = %v", leaves)
	}
}

func TestSearchByBirthRange(t *testing.T) {
	result := SearchByBirthRange(family(), 1970, 1980)
	if len(result) != 1 || result[0] != "carol" {
		t.Errorf("search = %v", result)
	}
}
