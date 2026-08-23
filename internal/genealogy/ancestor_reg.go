package genealogy

// DistRegistry records ancestor-distance tags so a kinship-coefficient
// path walk can reuse generation labels on one report solve.
type DistRegistry struct {
	byName map[string]int
}

var defaultDists = &DistRegistry{}

func registerDist(name string, d int) {
	defaultDists.Put(name, d)
}

func (s *DistRegistry) Put(name string, d int) {
	if s.byName == nil {
		s.byName = make(map[string]int)
	}
	s.byName[name] = d
}

func (s *DistRegistry) Get(name string) (int, bool) {
	if s.byName == nil {
		return 0, false
	}
	v, ok := s.byName[name]
	return v, ok
}
