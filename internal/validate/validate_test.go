package validate

import (
	"strings"
	"testing"

	"kinship/internal/genealogy"
)

func goodFamily() *genealogy.Family {
	lines := strings.Split(`PERSON alice F 1950
PERSON bob M 1948
PERSON carol F 1975
PARENT alice carol
PARENT bob carol`, "\n")
	f, _ := genealogy.ParseFile(lines)
	return f
}

func TestValidateGoodFamily(t *testing.T) {
	r := Validate(goodFamily())
	if !r.OK() {
		t.Errorf("expected OK, got errors: %+v", r.Errors())
	}
}

func TestValidateBirthOrder(t *testing.T) {
	lines := strings.Split(`PERSON parent F 2000
PERSON child M 1990
PARENT parent child`, "\n")
	f, _ := genealogy.ParseFile(lines)
	r := Validate(f)
	found := false
	for _, iss := range r.Issues {
		if strings.Contains(iss.Message, "born after child") {
			found = true
		}
	}
	if !found {
		t.Error("expected birth order warning")
	}
}

func TestUnconnectedGroups(t *testing.T) {
	lines := strings.Split(`PERSON alice F 1950
PERSON bob M 1960
PERSON carol F 1970
PARENT alice carol`, "\n")
	f, _ := genealogy.ParseFile(lines)
	groups := UnconnectedGroups(f)
	if len(groups) != 2 {
		t.Errorf("groups = %d, want 2 (bob is isolated)", len(groups))
	}
}

func TestConsistencyScore(t *testing.T) {
	f := goodFamily()
	score := ConsistencyScore(f)
	if score < 0.9 {
		t.Errorf("score = %f, want >= 0.9", score)
	}
}
