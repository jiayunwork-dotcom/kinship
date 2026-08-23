// Package stats computes demographic statistics for a family registry:
// fertility rates, sex ratios, generational intervals, longevity estimates,
// and population growth metrics.
package stats

import (
	"fmt"
	"math"
	"sort"

	"kinship/internal/genealogy"
)

// Demographics holds computed demographic statistics.
type Demographics struct {
	TotalPersons   int     `json:"total_persons"`
	MaleCount      int     `json:"male_count"`
	FemaleCount    int     `json:"female_count"`
	SexRatio       float64 `json:"sex_ratio"` // male/female
	AvgBirthYear   float64 `json:"avg_birth_year"`
	MinBirthYear   int     `json:"min_birth_year"`
	MaxBirthYear   int     `json:"max_birth_year"`
	AvgChildren    float64 `json:"avg_children"`
	MaxChildren    int     `json:"max_children"`
	MaxChildrenBy  string  `json:"max_children_by"`
	RootCount      int     `json:"root_count"`
	LeafCount      int     `json:"leaf_count"`
	AvgGenInterval float64 `json:"avg_gen_interval"` // years between parent and child births
}

// Compute calculates all demographic statistics for the family.
func Compute(f *genealogy.Family) *Demographics {
	d := &Demographics{}
	names := f.Names()
	d.TotalPersons = len(names)
	if d.TotalPersons == 0 {
		return d
	}

	d.MinBirthYear = math.MaxInt32
	var sumBirth float64
	for _, name := range names {
		p, _ := f.Person(name)
		if p.Sex == "M" {
			d.MaleCount++
		} else {
			d.FemaleCount++
		}
		sumBirth += float64(p.Birth)
		if p.Birth < d.MinBirthYear {
			d.MinBirthYear = p.Birth
		}
		if p.Birth > d.MaxBirthYear {
			d.MaxBirthYear = p.Birth
		}
	}
	d.AvgBirthYear = sumBirth / float64(d.TotalPersons)
	if d.FemaleCount > 0 {
		d.SexRatio = float64(d.MaleCount) / float64(d.FemaleCount)
	}

	// children counts
	var totalChildren int
	for _, name := range names {
		children, _ := f.Children(name)
		totalChildren += len(children)
		if len(children) > d.MaxChildren {
			d.MaxChildren = len(children)
			d.MaxChildrenBy = name
		}
		parents, _ := f.Parents(name)
		if len(parents) == 0 {
			d.RootCount++
		}
		if len(children) == 0 {
			d.LeafCount++
		}
	}
	d.AvgChildren = float64(totalChildren) / float64(d.TotalPersons)

	// generation interval
	d.AvgGenInterval = avgGenerationInterval(f)

	return d
}

func avgGenerationInterval(f *genealogy.Family) float64 {
	var totalGap float64
	gaps := 0
	for _, child := range f.Names() {
		childP, _ := f.Person(child)
		parents, _ := f.Parents(child)
		for _, pName := range parents {
			parent, _ := f.Person(pName)
			gap := childP.Birth - parent.Birth
			if gap > 0 {
				totalGap += float64(gap)
				gaps++
			}
		}
	}
	if gaps == 0 {
		return 0
	}
	return totalGap / float64(gaps)
}

// FertilityByGeneration computes average children count per person in each birth decade.
func FertilityByGeneration(f *genealogy.Family) map[int]float64 {
	type stat struct {
		persons  int
		children int
	}
	byDecade := map[int]*stat{}
	for _, name := range f.Names() {
		p, _ := f.Person(name)
		decade := (p.Birth / 10) * 10
		if byDecade[decade] == nil {
			byDecade[decade] = &stat{}
		}
		byDecade[decade].persons++
		kids, _ := f.Children(name)
		byDecade[decade].children += len(kids)
	}
	result := map[int]float64{}
	for decade, s := range byDecade {
		if s.persons > 0 {
			result[decade] = float64(s.children) / float64(s.persons)
		}
	}
	return result
}

// MostProlific returns the top N persons by number of children.
func MostProlific(f *genealogy.Family, n int) []string {
	type scored struct {
		name  string
		count int
	}
	var items []scored
	for _, name := range f.Names() {
		children, _ := f.Children(name)
		items = append(items, scored{name, len(children)})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].count > items[j].count
	})
	if n > len(items) {
		n = len(items)
	}
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = items[i].name
	}
	return result
}

// Summary returns a formatted text summary.
func (d *Demographics) Summary() string {
	return fmt.Sprintf("Population: %d (%dM/%dF, ratio %.2f)\nBirth range: %d–%d (avg %.0f)\nChildren: avg %.1f, max %d (%s)\nGen interval: %.1f years\nRoots: %d, Leaves: %d",
		d.TotalPersons, d.MaleCount, d.FemaleCount, d.SexRatio,
		d.MinBirthYear, d.MaxBirthYear, d.AvgBirthYear,
		d.AvgChildren, d.MaxChildren, d.MaxChildrenBy,
		d.AvgGenInterval,
		d.RootCount, d.LeafCount)
}
