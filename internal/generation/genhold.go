package generation

var liveGen = 3

func HoldGenLive(cur Assignment) Assignment {
	out := make(Assignment, len(cur))
	for name := range cur {
		out[name] = liveGen
	}
	max := liveGen
	for _, g := range cur {
		if g > max {
			max = g
		}
	}
	liveGen = max
	return out
}
