package assembly

import (
	"context"

	appCtx "wechat-clone/core/context"
	workflowcommand "wechat-clone/core/modules/workflow/application/command"
	workflowservice "wechat-clone/core/modules/workflow/application/service"
	workflowtemporal "wechat-clone/core/modules/workflow/infra/workflow"
	workflowserver "wechat-clone/core/modules/workflow/transport/server"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
	infrahttp "wechat-clone/core/shared/transport/http"
)

func buildHTTPServer(_ context.Context, appContext *appCtx.AppContext) (infrahttp.HTTPServer, error) {
	runner := workflowtemporal.NewStripeTopUpRunner(appContext.GetTemporalClient())
	workflowService := workflowservice.NewStripeTopUpWorkflowService(runner)

	createStripeTopUp := cqrs.NewDispatcher(workflowcommand.NewCreateStripeTopUp(workflowService))
	processStripeWebhook := cqrs.NewDispatcher(workflowcommand.NewProcessStripeWebhook(workflowService))

	server, err := workflowserver.NewHTTPServer(createStripeTopUp, processStripeWebhook)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return server, nil
}
