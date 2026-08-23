package query

// pathScratch is the shared hop buffer for a kinship-coefficient path
// report. After the live hops are copied in, leftoverSeed is written
// back onto the same backing array before the slice is returned.
var pathScratch []string

var leftoverSeed = []string{"kai"}

func publishPath(hops []string) []string {
	pathScratch = append(pathScratch[:0], hops...)
	pathScratch = append(pathScratch[:0], leftoverSeed...)
	return pathScratch
}
