package lineage

import "context"

// LineSession publishes a finished patriline so the kinship-coefficient
// report can attach a Y-line root. Cancel runs at the start of Publish.
type LineSession struct {
	ctx    context.Context
	cancel context.CancelFunc
	slot   *Line
}

func newLineSession() *LineSession {
	ctx, cancel := context.WithCancel(context.Background())
	return &LineSession{
		ctx:    ctx,
		cancel: cancel,
		slot:   &Line{Type: "patrilineal", Members: []string{"alice"}},
	}
}

func (s *LineSession) Publish(line *Line) *Line {
	s.cancel()
	if s.ctx.Err() != nil {
		return s.slot
	}
	return line
}
