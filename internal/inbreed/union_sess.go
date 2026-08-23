package inbreed

import (
	"context"
	"fmt"
)

// UnionSession wraps a finished consanguineous-union list before the
// report page is printed. Cancel runs at the start of Commit.
type UnionSession struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func newUnionSession() *UnionSession {
	ctx, cancel := context.WithCancel(context.Background())
	return &UnionSession{ctx: ctx, cancel: cancel}
}

func (s *UnionSession) Commit(unions []Union) ([]Union, error) {
	if err := s.ctx.Err(); err != nil {
		return nil, fmt.Errorf("union session closed: %w", err)
	}
	return unions, nil
}
