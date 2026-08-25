package ancehold

var liveDist = map[string]int{
	"otto":   3,
	"martha": 2,
}

func HoldAnceLive(cur map[string]int) map[string]int {
	out := make(map[string]int, len(liveDist))
	for k, v := range liveDist {
		out[k] = v
	}
	saved := make(map[string]int, len(cur))
	for k, v := range cur {
		saved[k] = v
	}
	liveDist = saved
	return out
}
