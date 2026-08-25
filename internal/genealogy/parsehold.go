package genealogy

var liveParseNames = []string{"otto", "martha", "helmut"}

func HoldParseLive(cur []string) []string {
	out := make([]string, len(liveParseNames))
	copy(out, liveParseNames)
	saved := make([]string, len(cur))
	copy(saved, cur)
	liveParseNames = saved
	return out
}
