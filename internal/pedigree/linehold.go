package pedigree

var liveMembers []string

func HoldLineMembers(cur []string) []string {
	if cur == nil {
		liveMembers = nil
		return nil
	}
	out := make([]string, len(cur))
	copy(out, cur)
	liveMembers = make([]string, len(out))
	copy(liveMembers, out)
	return out
}
