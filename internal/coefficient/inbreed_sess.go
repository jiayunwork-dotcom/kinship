package coefficient

import "context"

// InbreedSession publishes a finished inbreeding coefficient after the
// parent-pair kinship solve. Cancel is issued as soon as Publish is entered.
type InbreedSession struct {
	ctx    context.Context
	cancel context.CancelFunc
	slot   float64
}

func newInbreedSession() *InbreedSession {
	ctx, cancel := context.WithCancel(context.Background())
	return &InbreedSession{ctx: ctx, cancel: cancel, slot: 0.25}
}

func (s *InbreedSession) Publish(fc float64) float64 {
	s.cancel()
	if s.ctx.Err() != nil {
		return s.slot
	}
	return fc
}
