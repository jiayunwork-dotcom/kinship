package genealogy

var liveTerm string

func HoldTermLive(cur string) string {
	if cur == "" {
		return liveTerm
	}
	liveTerm = cur
	return cur
}
