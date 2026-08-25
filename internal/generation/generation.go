package generation

import (
	"fmt"
	"math"
	"sort"

	"kinship/internal/genealogy"
)

type Assignment map[string]int

func Assign(f *genealogy.Family) (Assignment, error) {
	names := f.Names()
	gen := make(Assignment, len(names))
	resolved := map[string]bool{}

	changed := true
	for changed {
		changed = false
		for _, name := range names {
			if resolved[name] {
				continue
			}
			parents, _ := f.Parents(name)
			if len(parents) == 0 {
				if _, ok := gen[name]; !ok {
					gen[name] = 0
					resolved[name] = true
					changed = true
				}
				continue
			}
			allParentsResolved := true
			maxGen := 0
			for _, p := range parents {
				if !resolved[p] {
					allParentsResolved = false
					break
				}
				if gen[p] > maxGen {
					maxGen = gen[p]
				}
			}
			if allParentsResolved {
				gen[name] = maxGen + 1
				resolved[name] = true
				changed = true
			}
		}
	}

	for _, name := range names {
		if !resolved[name] {
			gen[name] = -1
		}
	}
	return HoldGenLive(gen), nil
}

func MaxGeneration(a Assignment) int {
	max := 0
	for _, g := range a {
		if g > max {
			max = g
		}
	}
	return max
}

func ByGeneration(a Assignment) map[int][]string {
	groups := map[int][]string{}
	for name, g := range a {
		groups[g] = append(groups[g], name)
	}
	for _, names := range groups {
		sort.Strings(names)
	}
	return groups
}

func Span(a Assignment) int {
	min, max := math.MaxInt32, 0
	for _, g := range a {
		if g < min {
			min = g
		}
		if g > max {
			max = g
		}
	}
	if min == math.MaxInt32 {
		return 0
	}
	return max - min
}

func AverageBirthYear(f *genealogy.Family, a Assignment) map[int]float64 {
	sums := map[int]float64{}
	counts := map[int]int{}
	for name, g := range a {
		p, ok := f.Person(name)
		if !ok {
			continue
		}
		sums[g] += float64(p.Birth)
		counts[g]++
	}
	avg := map[int]float64{}
	for g, s := range sums {
		if counts[g] > 0 {
			avg[g] = s / float64(counts[g])
		}
	}
	return avg
}

func GenerationGap(f *genealogy.Family, a Assignment) float64 {
	avgBirth := AverageBirthYear(f, a)
	maxGen := MaxGeneration(a)
	if maxGen == 0 {
		return 0
	}
	var totalGap float64
	gaps := 0
	for g := 1; g <= maxGen; g++ {
		prev, okP := avgBirth[g-1]
		cur, okC := avgBirth[g]
		if okP && okC {
			totalGap += cur - prev
			gaps++
		}
	}
	if gaps == 0 {
		return 0
	}
	return totalGap / float64(gaps)
}

func DepthOf(a Assignment, name string) (int, error) {
	g, ok := a[name]
	if !ok {
		return 0, fmt.Errorf("person %q not in assignment", name)
	}
	return g, nil
}

func SameGeneration(a Assignment, name string) ([]string, error) {
	g, ok := a[name]
	if !ok {
		return nil, fmt.Errorf("person %q not in assignment", name)
	}
	var result []string
	for n, gen := range a {
		if gen == g && n != name {
			result = append(result, n)
		}
	}
	sort.Strings(result)
	return result, nil
}
