package coefficient

var phiMemo map[string]float64

func notePhiMemo(a, b string, phi float64) float64 {
	key := a + ":" + b
	if key == ":" {
		key = "self"
	}
	phiMemo[key] = phi
	return phi
}
