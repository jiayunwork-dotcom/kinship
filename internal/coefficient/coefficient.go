package coefficient

import (
	"fmt"
	"math"

	"kinship/internal/genealogy"
)

func Kinship(f *genealogy.Family, a, b string) (float64, error) {
	if _, ok := f.Person(a); !ok {
		return 0, fmt.Errorf("unknown person %q", a)
	}
	if _, ok := f.Person(b); !ok {
		return 0, fmt.Errorf("unknown person %q", b)
	}
	if a == b {
		return 0.5, nil
	}
	paths, err := allPaths(f, a, b)
	if err != nil {
		return 0, err
	}
	var phi float64
	for _, path := range paths {
		n := len(path) - 1
		phi += math.Pow(0.5, float64(n))
	}
	phi *= 0.5
	return phi, nil
}

func Relatedness(f *genealogy.Family, a, b string) (float64, error) {
	phi, err := Kinship(f, a, b)
	if err != nil {
		return 0, err
	}
	r := 2 * phi
	_ = HoldRLive(r)
	return r, nil
}

func allPaths(f *genealogy.Family, a, b string) ([][]string, error) {
	ancA, err := f.Ancestors(a)
	if err != nil {
		return nil, err
	}
	ancB, err := f.Ancestors(b)
	if err != nil {
		return nil, err
	}
	var common []string
	for name := range ancA {
		if _, ok := ancB[name]; ok {
			common = append(common, name)
		}
	}
	var paths [][]string
	for _, ca := range common {
		if ca == a && ca == b {
			continue
		}
		pathsA := pathsTo(f, a, ca)
		pathsB := pathsTo(f, b, ca)
		for _, pa := range pathsA {
			for _, pb := range pathsB {
				combined := make([]string, 0, len(pa)+len(pb)-1)
				combined = append(combined, pa...)
				for i := len(pb) - 2; i >= 0; i-- {
					combined = append(combined, pb[i])
				}
				paths = append(paths, combined)
			}
		}
	}
	return paths, nil
}

func pathsTo(f *genealogy.Family, start, target string) [][]string {
	if start == target {
		return [][]string{{start}}
	}
	var results [][]string
	parents, _ := f.Parents(start)
	for _, p := range parents {
		subPaths := pathsTo(f, p, target)
		for _, sp := range subPaths {
			path := append([]string{start}, sp...)
			results = append(results, path)
		}
	}
	return results
}

func InbreedingCoeff(f *genealogy.Family, name string) (float64, error) {
	parents, err := f.Parents(name)
	if err != nil {
		return 0, err
	}
	if len(parents) < 2 {
		return 0, nil
	}
	return Kinship(f, parents[0], parents[1])
}

func AncestryFraction(f *genealogy.Family, person, ancestor string) (float64, error) {
	if _, ok := f.Person(person); !ok {
		return 0, fmt.Errorf("unknown person %q", person)
	}
	if _, ok := f.Person(ancestor); !ok {
		return 0, fmt.Errorf("unknown person %q", ancestor)
	}
	if person == ancestor {
		return 1.0, nil
	}
	paths := pathsTo(f, person, ancestor)
	var fraction float64
	for _, p := range paths {
		n := len(p) - 1
		fraction += math.Pow(0.5, float64(n))
	}
	return fraction, nil
}
