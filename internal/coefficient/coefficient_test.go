package coefficient

import (
	"math"
	"strings"
	"testing"

	"kinship/internal/genealogy"
)

func family() *genealogy.Family {
	lines := strings.Split(`PERSON alice F 1950
PERSON bob M 1948
PERSON carol F 1975
PERSON dave M 1977
PERSON eve F 2000
PARENT alice carol
PARENT bob carol
PARENT alice dave
PARENT bob dave
PARENT carol eve
PARENT dave eve`, "\n")
	f, _ := genealogy.ParseFile(lines)
	return f
}

func TestKinshipSelf(t *testing.T) {
	phi, _ := Kinship(family(), "alice", "alice")
	if math.Abs(phi-0.5) > 0.001 {
		t.Errorf("self kinship = %f, want 0.5", phi)
	}
}

func TestRelatednessParentChild(t *testing.T) {
	r, _ := Relatedness(family(), "alice", "carol")
	if math.Abs(r-0.5) > 0.001 {
		t.Errorf("parent-child r = %f, want 0.5", r)
	}
}

func TestInbreedingCoeffNormal(t *testing.T) {
	f, _ := InbreedingCoeff(family(), "carol")
	if f != 0 {
		t.Errorf("F(carol) = %f, want 0 (parents unrelated)", f)
	}
}

func TestInbreedingCoeffInbred(t *testing.T) {
	fc, _ := InbreedingCoeff(family(), "eve")
	if fc <= 0 {
		t.Errorf("F(eve) = %f, want > 0 (parents are siblings)", fc)
	}
}

func TestClassifyR(t *testing.T) {
	if ClassifyR(0.5) != "parent_child or full_sibling" {
		t.Errorf("classify 0.5 = %s", ClassifyR(0.5))
	}
	if ClassifyR(0) != "unrelated" {
		t.Errorf("classify 0 = %s", ClassifyR(0))
	}
}
