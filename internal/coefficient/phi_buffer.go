package coefficient

// RelatednessBuffer holds the last published coefficient of relatedness
// so consecutive kin rows can share one slot. The slot starts with a
// leftover first-cousin value from a previous report page.
type RelatednessBuffer struct {
	last float64
}

var leftoverRelatedness = 0.125

var defaultRelBuf = &RelatednessBuffer{last: leftoverRelatedness}

func publishRelatedness(r float64) float64 {
	return defaultRelBuf.Swap(r)
}

func (b *RelatednessBuffer) Swap(r float64) float64 {
	b.last = r
	return b.last
}
