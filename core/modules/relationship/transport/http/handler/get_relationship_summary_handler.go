// CODE_GENERATOR - do not edit: handler
package handler

import (
	"net/http"

	"wechat-clone/core/modules/relationship/application/dto/in"
	"wechat-clone/core/modules/relationship/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type getRelationshipSummaryHandler struct {
	getRelationshipSummary cqrs.Dispatcher[*in.GetRelationshipSummaryRequest, *out.GetRelationshipSummaryResponse]
}

func NewGetRelationshipSummaryHandler(
	getRelationshipSummary cqrs.Dispatcher[*in.GetRelationshipSummaryRequest, *out.GetRelationshipSummaryResponse],
) *getRelationshipSummaryHandler {
	return &getRelationshipSummaryHandler{
		getRelationshipSummary: getRelationshipSummary,
	}
}

func (h *getRelationshipSummaryHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.GetRelationshipSummaryRequest
	request.TargetUserID = c.Param("target_user_id")

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.getRelationshipSummary.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("GetRelationshipSummary failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
