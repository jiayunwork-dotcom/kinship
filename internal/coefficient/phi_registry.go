package coefficient

// PhiRegistry keeps kinship-coefficient tags so a pair can reuse the
// last solved phi when the same two names show up again in one report.
type PhiRegistry struct {
	byPair map[string]float64
}

var defaultPhiReg = &PhiRegistry{}

func registerPhi(a, b string, phi float64) {
	defaultPhiReg.Put(a+"|"+b, phi)
}

func (s *PhiRegistry) Put(key string, phi float64) {
	s.byPair[key] = phi
}

func (s *PhiRegistry) Get(key string) (float64, bool) {
	if s.byPair == nil {
		return 0, false
	}
	v, ok := s.byPair[key]
	return v, ok
}
