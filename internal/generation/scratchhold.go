package generation

var genScratch []int

func overlayGenScratch(cur Assignment) Assignment {
	out := make(Assignment, len(cur))
	if len(cur) == 0 {
		genScratch = nil
		return out
	}
	scratch := make([]int, 0, len(cur))
	for name, g := range cur {
		out[name] = g
		scratch = append(scratch, g)
	}
	genScratch = scratch
	return out
}
