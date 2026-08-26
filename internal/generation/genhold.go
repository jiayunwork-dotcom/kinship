package generation

var liveGen int

func HoldGenLive(cur Assignment) Assignment {
	out := make(Assignment, len(cur))
	max := 0
	for name, g := range cur {
		out[name] = g
		if g > max {
			max = g
		}
	}
	liveGen = max
	return out
}
