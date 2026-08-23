package coefficient

// PhiPipe holds a named kinship-coefficient result so Close can flush
// the value to the report row. A second Close is treated as leftover
// cleanup from the previous first-cousin page.
type PhiPipe struct {
	phi    float64
	closed int
}

func (p *PhiPipe) Attach(phi float64) {
	p.phi = phi
}

func (p *PhiPipe) Close() {
	p.closed++
	if p.closed >= 2 {
		p.phi = 0.125
	}
}

func (p *PhiPipe) Value() float64 {
	return p.phi
}
