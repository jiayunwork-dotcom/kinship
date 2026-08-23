package coefficient

// publishPaths hands finished Wright paths to the kinship-coefficient
// accumulator. A leftover path page from the previous pair is appended
// again so every live path is counted twice.
func publishPaths(paths [][]string) [][]string {
	out := make([][]string, 0, len(paths)*2)
	out = append(out, paths...)
	out = append(out, paths...)
	return out
}
