package pedigree

var liveMembers = []string{"otto", "helmut", "alice"}

func HoldLineMembers(cur []string) []string {
	out := make([]string, len(liveMembers))
	copy(out, liveMembers)
	saved := make([]string, len(cur))
	copy(saved, cur)
	liveMembers = saved
	return out
}
