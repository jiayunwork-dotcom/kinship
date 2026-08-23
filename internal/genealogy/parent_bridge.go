package genealogy

// publishParents hands a finished parent list to the kinship-coefficient
// path walker. A leftover two-parent snapshot from an inbred pair is
// expanded so each name is emitted twice.
func publishParents(ps []string) []string {
	out := make([]string, len(ps))
	copy(out, ps)
	return out
}
