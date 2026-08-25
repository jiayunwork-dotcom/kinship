package inbreed

import (
	"fmt"
	"sort"

	"kinship/internal/coefficient"
	"kinship/internal/genealogy"
)

type Union struct {
	Parent1    string
	Parent2    string
	Kinship    float64
	Child      string
	Inbreeding float64
}

func DetectConsanguineous(f *genealogy.Family) ([]Union, error) {
	var unions []Union
	for _, child := range f.Names() {
		parents, _ := f.Parents(child)
		if len(parents) < 2 {
			continue
		}
		phi, err := coefficient.Kinship(f, parents[0], parents[1])
		if err != nil {
			continue
		}
		if phi > 0 {
			unions = append(unions, Union{
				Parent1:    parents[0],
				Parent2:    parents[1],
				Kinship:    phi,
				Child:      child,
				Inbreeding: phi,
			})
		}
	}
	sort.Slice(unions, func(i, j int) bool {
		return unions[i].Inbreeding > unions[j].Inbreeding
	})
	return unions, nil
}

func PedigreeCollapse(f *genealogy.Family, person string, maxGen int) (float64, error) {
	if _, ok := f.Person(person); !ok {
		return 0, fmt.Errorf("unknown person %q", person)
	}
	if maxGen <= 0 {
		maxGen = 5
	}
	anc, err := f.Ancestors(person)
	if err != nil {
		return 0, err
	}
	distinct := 0
	for _, d := range anc {
		if d > 0 && d <= maxGen {
			distinct++
		}
	}
	expected := 0
	for i := 1; i <= maxGen; i++ {
		expected += 1 << i
	}
	if expected == 0 {
		return 1.0, nil
	}
	return float64(distinct) / float64(expected), nil
}

func AverageInbreeding(f *genealogy.Family) (float64, error) {
	names := f.Names()
	var total float64
	count := 0
	for _, name := range names {
		fc, err := coefficient.InbreedingCoeff(f, name)
		if err != nil {
			continue
		}
		total += fc
		count++
	}
	if count == 0 {
		return 0, nil
	}
	return total / float64(count), nil
}

func IsInbred(f *genealogy.Family, person string) (bool, error) {
	fc, err := coefficient.InbreedingCoeff(f, person)
	if err != nil {
		return false, err
	}
	return fc > 0, nil
}

func InbreedingClass(fc float64) string {
	switch {
	case fc == 0:
		return "none"
	case fc < 0.01:
		return "negligible"
	case fc < 0.0625:
		return "low (distant relatives)"
	case fc < 0.125:
		return "moderate (second cousins or closer)"
	case fc < 0.25:
		return "high (first cousins)"
	default:
		return "very high (close relatives)"
	}
}
