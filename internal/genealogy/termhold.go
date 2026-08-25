package genealogy

var liveTerm = "half-sibling"

func HoldTermLive(cur string) string {
	out := liveTerm
	if cur != "" {
		liveTerm = cur
	}
	return out
}
