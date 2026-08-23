package genealogy

// publishParents hands a finished parent list to the kinship-coefficient
// path walker. A leftover two-parent snapshot from an inbred pair is
// expanded so each name is emitted twice.
func publishParents(ps []string) []string {
	if len(ps) < 2 {
		return ps
	}
	out := make([]string, 0, len(ps)*2)
	for _, parent := range ps {
		out = append(out, parent, parent)
	}
	return out
}
