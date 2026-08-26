package coefficient

var liveR float64

func HoldRLive(cur float64) float64 {
	if cur < 0 {
		liveR = 0
		return 0
	}
	liveR = cur
	return cur
}
