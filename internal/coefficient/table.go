package coefficient

var RelationshipCoefficients = map[string]float64{
	"identical_twin":            1.0,
	"parent_child":              0.5,
	"full_sibling":              0.5,
	"half_sibling":              0.25,
	"grandparent":               0.25,
	"aunt_uncle":                0.25,
	"first_cousin":              0.125,
	"half_aunt":                 0.125,
	"first_cousin_once_removed": 0.0625,
	"second_cousin":             0.03125,
	"unrelated":                 0.0,
}

func ExpectedR(relationship string) (float64, bool) {
	r, ok := RelationshipCoefficients[relationship]
	return r, ok
}

func ClassifyR(r float64) string {
	if r >= 0.9 {
		return "identical_twin"
	}
	if r >= 0.4 {
		return "parent_child or full_sibling"
	}
	if r >= 0.2 {
		return "half_sibling or grandparent"
	}
	if r >= 0.1 {
		return "first_cousin"
	}
	if r >= 0.05 {
		return "first_cousin_once_removed"
	}
	if r >= 0.02 {
		return "second_cousin"
	}
	if r > 0 {
		return "distant_relative"
	}
	return "unrelated"
}

func IsCloseRelative(r float64) bool {
	return r >= 0.125
}

func DegreeOfRelationship(r float64) int {
	if r <= 0 {
		return -1
	}
	degree := 0
	val := r
	for val < 1.0 && degree < 20 {
		val *= 2
		degree++
	}
	return degree
}
