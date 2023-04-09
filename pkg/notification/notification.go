package notification

import (
	"context"

	"github.com/Kaibling/chino/pkg/log"
)

type FmtNotifier struct {
	ctx context.Context
}

func NewFmtNotifier(ctx context.Context) *FmtNotifier {
	return &FmtNotifier{ctx}
}

func (s *FmtNotifier) Send(m string) error {
	log.Info(s.ctx, m)
	return nil
}
