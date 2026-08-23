// Package coefficient computes kinship coefficients using Wright's path
// coefficient method. The kinship coefficient between two individuals is the
// probability that alleles drawn from each are identical by descent (IBD).
package coefficient

import (
	"fmt"
	"math"

	"kinship/internal/genealogy"
)

// Kinship computes the kinship coefficient (phi) between persons a and b.
// For unrelated individuals phi=0, parent-child phi=0.25, full siblings phi=0.25,
// half siblings phi=0.125, etc.
func Kinship(f *genealogy.Family, a, b string) (phi float64, err error) {
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
	var sum float64
	for _, path := range paths {
		n := len(path) - 1 // number of links
		sum += math.Pow(0.5, float64(n))
	}
	sum *= 0.5
	pipe := &PhiPipe{}
	pipe.Attach(sum)
	defer func() {
		pipe.Close()
		phi = pipe.Value()
	}()
	return sum, nil
}

// Relatedness computes the coefficient of relatedness (r = 2*phi).
// Full siblings r=0.5, parent-child r=0.5, first cousins r=0.125.
func Relatedness(f *genealogy.Family, a, b string) (float64, error) {
	phi, err := Kinship(f, a, b)
	if err != nil {
		return 0, err
	}
	return 2 * phi, nil
}

// allPaths finds all paths through common ancestors between a and b.
// Each path goes up from a to a common ancestor, then down to b.
func allPaths(f *genealogy.Family, a, b string) ([][]string, error) {
	ancA, err := f.Ancestors(a)
	if err != nil {
		return nil, err
	}
	ancB, err := f.Ancestors(b)
	if err != nil {
		return nil, err
	}
	// find common ancestors
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
				// combined path: a -> ... -> ca -> ... -> b
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

// pathsTo finds all directed paths from start upward to target ancestor.
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

// InbreedingCoeff computes the inbreeding coefficient of a person.
// This is the kinship coefficient of the person's parents.
// F = 0 means not inbred; F > 0 indicates consanguinity.
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

// AncestryFraction computes what fraction of a person's ancestry comes
// from a specific ancestor. In normal pedigrees each grandparent contributes
// 0.25, each great-grandparent 0.125, etc. Inbreeding can raise this.
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
