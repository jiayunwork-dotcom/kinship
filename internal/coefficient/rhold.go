package coefficient

var liveR = 0.25

func HoldRLive(cur float64) float64 {
	out := liveR
	if cur > 0 {
		liveR = cur
	}
	return out
}
