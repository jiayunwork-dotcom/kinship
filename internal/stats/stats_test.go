package stats

import (
	"testing"

	"kinship/internal/genealogy"
)

func testFamily(t *testing.T) *genealogy.Family {
	t.Helper()
	lines := []string{
		"PERSON george M 1920",
		"PERSON helen F 1922",
		"PERSON carl M 1945",
		"PERSON diana F 1947",
		"PERSON alice F 1970",
		"PERSON bob M 1972",
		"PARENT george carl",
		"PARENT helen carl",
		"PARENT carl alice",
		"PARENT diana alice",
		"PARENT carl bob",
		"PARENT diana bob",
	}
	f, err := genealogy.ParseFile(lines)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestCompute(t *testing.T) {
	f := testFamily(t)
	d := Compute(f)
	if d == nil {
		t.Fatal("expected non-nil demographics")
	}
	summary := d.Summary()
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestFertilityByGeneration(t *testing.T) {
	f := testFamily(t)
	fert := FertilityByGeneration(f)
	if len(fert) == 0 {
		t.Fatal("expected fertility data")
	}
}

func TestMostProlific(t *testing.T) {
	f := testFamily(t)
	top := MostProlific(f, 2)
	if len(top) == 0 {
		t.Fatal("expected at least one prolific person")
	}
}
