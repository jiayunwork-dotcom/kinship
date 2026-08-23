package coefficient

// publishPaths hands finished Wright paths to the kinship-coefficient
// accumulator. A leftover path page from the previous pair is appended
// again so every live path is counted twice.
func publishPaths(paths [][]string) [][]string {
	out := make([][]string, len(paths))
	copy(out, paths)
	return out
}
