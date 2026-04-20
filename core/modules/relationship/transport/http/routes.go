// CODE_GENERATOR - do not edit: routing
package http

import (
	"wechat-clone/core/modules/relationship/application/dto/in"
	"wechat-clone/core/modules/relationship/application/dto/out"
	"wechat-clone/core/modules/relationship/transport/http/handler"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(_ *gin.RouterGroup) {}
func RegisterPrivateRoutes(
	routes *gin.RouterGroup,
	sendFriendRequest cqrs.Dispatcher[*in.SendFriendRequestRequest, *out.SendFriendRequestResponse],
	cancelFriendRequest cqrs.Dispatcher[*in.CancelFriendRequestRequest, *out.CancelFriendRequestResponse],
	acceptFriendRequest cqrs.Dispatcher[*in.AcceptFriendRequestRequest, *out.AcceptFriendRequestResponse],
	rejectFriendRequest cqrs.Dispatcher[*in.RejectFriendRequestRequest, *out.RejectFriendRequestResponse],
	listIncomingFriendRequests cqrs.Dispatcher[*in.ListIncomingFriendRequestsRequest, *out.ListIncomingFriendRequestsResponse],
	listOutgoingFriendRequests cqrs.Dispatcher[*in.ListOutgoingFriendRequestsRequest, *out.ListOutgoingFriendRequestsResponse],
	unfriendUser cqrs.Dispatcher[*in.UnfriendUserRequest, *out.UnfriendUserResponse],
	listFriends cqrs.Dispatcher[*in.ListFriendsRequest, *out.ListFriendsResponse],
	followUser cqrs.Dispatcher[*in.FollowUserRequest, *out.FollowUserResponse],
	unfollowUser cqrs.Dispatcher[*in.UnfollowUserRequest, *out.UnfollowUserResponse],
	listFollowers cqrs.Dispatcher[*in.ListFollowersRequest, *out.ListFollowersResponse],
	listFollowing cqrs.Dispatcher[*in.ListFollowingRequest, *out.ListFollowingResponse],
	blockUser cqrs.Dispatcher[*in.BlockUserRequest, *out.BlockUserResponse],
	unblockUser cqrs.Dispatcher[*in.UnblockUserRequest, *out.UnblockUserResponse],
	listBlockedUsers cqrs.Dispatcher[*in.ListBlockedUsersRequest, *out.ListBlockedUsersResponse],
	getRelationshipStatus cqrs.Dispatcher[*in.GetRelationshipStatusRequest, *out.GetRelationshipStatusResponse],
	getMutualFriends cqrs.Dispatcher[*in.GetMutualFriendsRequest, *out.GetMutualFriendsResponse],
	getRelationshipSummary cqrs.Dispatcher[*in.GetRelationshipSummaryRequest, *out.GetRelationshipSummaryResponse],
) {
	routes.POST("/relationship/friend-requests", httpx.Wrap(handler.NewSendFriendRequestHandler(sendFriendRequest)))
	routes.DELETE("/relationship/friend-requests/:target_user_id", httpx.Wrap(handler.NewCancelFriendRequestHandler(cancelFriendRequest)))
	routes.POST("/relationship/friend-requests/:requester_user_id/accept", httpx.Wrap(handler.NewAcceptFriendRequestHandler(acceptFriendRequest)))
	routes.POST("/relationship/friend-requests/:requester_user_id/reject", httpx.Wrap(handler.NewRejectFriendRequestHandler(rejectFriendRequest)))
	routes.GET("/relationship/friend-requests/incoming", httpx.Wrap(handler.NewListIncomingFriendRequestsHandler(listIncomingFriendRequests)))
	routes.GET("/relationship/friend-requests/outgoing", httpx.Wrap(handler.NewListOutgoingFriendRequestsHandler(listOutgoingFriendRequests)))
	routes.DELETE("/relationship/friends/:target_user_id", httpx.Wrap(handler.NewUnfriendUserHandler(unfriendUser)))
	routes.GET("/relationship/friends", httpx.Wrap(handler.NewListFriendsHandler(listFriends)))
	routes.POST("/relationship/follows", httpx.Wrap(handler.NewFollowUserHandler(followUser)))
	routes.DELETE("/relationship/follows/:target_user_id", httpx.Wrap(handler.NewUnfollowUserHandler(unfollowUser)))
	routes.GET("/relationship/followers", httpx.Wrap(handler.NewListFollowersHandler(listFollowers)))
	routes.GET("/relationship/following", httpx.Wrap(handler.NewListFollowingHandler(listFollowing)))
	routes.POST("/relationship/blocks", httpx.Wrap(handler.NewBlockUserHandler(blockUser)))
	routes.DELETE("/relationship/blocks/:target_user_id", httpx.Wrap(handler.NewUnblockUserHandler(unblockUser)))
	routes.GET("/relationship/blocks", httpx.Wrap(handler.NewListBlockedUsersHandler(listBlockedUsers)))
	routes.GET("/relationship/status/:target_user_id", httpx.Wrap(handler.NewGetRelationshipStatusHandler(getRelationshipStatus)))
	routes.GET("/relationship/mutual-friends/:target_user_id", httpx.Wrap(handler.NewGetMutualFriendsHandler(getMutualFriends)))
	routes.GET("/relationship/summary/:target_user_id", httpx.Wrap(handler.NewGetRelationshipSummaryHandler(getRelationshipSummary)))
}
