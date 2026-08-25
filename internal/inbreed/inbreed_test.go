package inbreed

import (
	"testing"

	"kinship/internal/genealogy"
)

func testFamily(t *testing.T) *genealogy.Family {
	t.Helper()
	lines := []string{
		"PERSON alice F 1940",
		"PERSON bob M 1938",
		"PERSON carol F 1965",
		"PERSON dave M 1963",
		"PERSON eve F 1990",
		"PARENT alice carol",
		"PARENT bob carol",
		"PARENT alice dave",
		"PARENT bob dave",
		"PARENT carol eve",
		"PARENT dave eve",
	}
	f, err := genealogy.ParseFile(lines)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestDetectConsanguineous(t *testing.T) {
	f := testFamily(t)
	unions, err := DetectConsanguineous(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(unions) == 0 {
		t.Fatal("expected at least one consanguineous union")
	}
}

func TestIsInbred(t *testing.T) {
	f := testFamily(t)
	inbred, err := IsInbred(f, "eve")
	if err != nil {
		t.Fatal(err)
	}
	if !inbred {
		t.Fatal("expected eve to be inbred")
	}
}

func TestInbreedingClass(t *testing.T) {
	cls := InbreedingClass(0.25)
	if cls == "" {
		t.Fatal("expected non-empty class for fc=0.25")
	}
}

func TestPedigreeCollapse(t *testing.T) {
	f := testFamily(t)
	ratio, err := PedigreeCollapse(f, "eve", 3)
	if err != nil {
		t.Fatal(err)
	}
	if ratio <= 0 {
		t.Fatalf("expected positive collapse ratio, got %f", ratio)
	}
}
