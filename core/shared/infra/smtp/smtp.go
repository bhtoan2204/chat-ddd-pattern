package smtp

import (
	"context"
	"wechat-clone/core/shared/pkg/logging"

	"go.uber.org/zap"
)

type SMTP struct {
}

func NewSMTP() SMTP {
	return SMTP{}
}

func (s SMTP) Send(ctx context.Context, to, subject, body string) error {
	log := logging.FromContext(ctx).Named("Send")
	// Not integrated with SMTP server yet
	log.Infow("Sending email",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.String("body", body),
	)
	return nil
}
