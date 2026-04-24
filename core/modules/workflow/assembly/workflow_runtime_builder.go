package assembly

import (
	"context"

	appCtx "wechat-clone/core/context"
	workflowledger "wechat-clone/core/modules/workflow/external/ledger"
	workflowpayment "wechat-clone/core/modules/workflow/external/payment"
	workflowtemporal "wechat-clone/core/modules/workflow/infra/workflow"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"
	modruntime "wechat-clone/core/shared/runtime"
)

func buildWorkflowRuntime(cfg *config.Config, appContext *appCtx.AppContext) (modruntime.Module, error) {
	ctx := context.Background()
	paymentClient, err := workflowpayment.New(ctx, cfg)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	ledgerClient, err := workflowledger.New(ctx, cfg)
	if err != nil {
		if closeErr := paymentClient.Close(); closeErr != nil {
			return nil, stackErr.Error(closeErr)
		}
		return nil, stackErr.Error(err)
	}

	return workflowtemporal.NewWorkerRuntime(appContext.GetTemporalClient(), paymentClient, ledgerClient), nil
}
