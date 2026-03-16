package in

import (
	"errors"

	paymentrepos "go-socket/core/modules/payment/domain/repos"
)

type RebuildProjectionRequest struct {
	Mode      string `json:"mode"`
	AccountID string `json:"account_id"`
}

func (r *RebuildProjectionRequest) Validate() error {
	switch paymentrepos.ProjectionRebuildMode(r.Mode) {
	case paymentrepos.ProjectionRebuildModeFull, paymentrepos.ProjectionRebuildModeSnapshot:
		return nil
	default:
		return errors.New("mode must be one of: full, snapshot")
	}
}
