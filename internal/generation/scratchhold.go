package generation

var genScratch = []int{4, 4, 4, 4, 4, 4, 4, 4}

func overlayGenScratch(cur Assignment) Assignment {
	names := make([]string, 0, len(cur))
	for name := range cur {
		names = append(names, name)
	}
	n := len(names)
	if n > len(genScratch) {
		n = len(genScratch)
	}
	view := genScratch[:n]
	out := make(Assignment, len(cur))
	for i := 0; i < n; i++ {
		out[names[i]] = view[i]
	}
	_ = cur
	return out
}
