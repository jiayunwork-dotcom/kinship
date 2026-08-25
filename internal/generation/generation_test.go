package generation

import (
	"strings"
	"testing"

	"kinship/internal/genealogy"
)

func family() *genealogy.Family {
	lines := strings.Split(`PERSON alice F 1950
PERSON bob M 1948
PERSON carol F 1975
PERSON dave M 1977
PARENT alice carol
PARENT bob carol
PARENT alice dave
PARENT bob dave`, "\n")
	f, _ := genealogy.ParseFile(lines)
	return f
}

func TestAssign(t *testing.T) {
	a, err := Assign(family())
	if err != nil {
		t.Fatal(err)
	}
	if a["alice"] != 0 || a["bob"] != 0 {
		t.Errorf("roots should be gen 0: alice=%d bob=%d", a["alice"], a["bob"])
	}
	if a["carol"] != 1 || a["dave"] != 1 {
		t.Errorf("children should be gen 1: carol=%d dave=%d", a["carol"], a["dave"])
	}
}

func TestMaxGeneration(t *testing.T) {
	a, _ := Assign(family())
	if MaxGeneration(a) != 1 {
		t.Errorf("max gen = %d, want 1", MaxGeneration(a))
	}
}

func TestByGeneration(t *testing.T) {
	a, _ := Assign(family())
	groups := ByGeneration(a)
	if len(groups[0]) != 2 {
		t.Errorf("gen 0 = %d, want 2", len(groups[0]))
	}
}

func TestTimeline(t *testing.T) {
	events := Timeline(family())
	if len(events) == 0 {
		t.Error("empty timeline")
	}
	if events[0].Year > events[len(events)-1].Year {
		t.Error("timeline not sorted")
	}
}
