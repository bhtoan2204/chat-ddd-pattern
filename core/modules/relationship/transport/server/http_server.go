// CODE_GENERATOR: registry
package server

import (
	"context"

	"wechat-clone/core/modules/relationship/application/dto/in"
	"wechat-clone/core/modules/relationship/application/dto/out"
	relationshiphttp "wechat-clone/core/modules/relationship/transport/http"
	"wechat-clone/core/shared/pkg/cqrs"
	infrahttp "wechat-clone/core/shared/transport/http"

	"github.com/gin-gonic/gin"
)

type relationshipHTTPServer struct {
	sendFriendRequest          cqrs.Dispatcher[*in.SendFriendRequestRequest, *out.SendFriendRequestResponse]
	cancelFriendRequest        cqrs.Dispatcher[*in.CancelFriendRequestRequest, *out.CancelFriendRequestResponse]
	acceptFriendRequest        cqrs.Dispatcher[*in.AcceptFriendRequestRequest, *out.AcceptFriendRequestResponse]
	rejectFriendRequest        cqrs.Dispatcher[*in.RejectFriendRequestRequest, *out.RejectFriendRequestResponse]
	listIncomingFriendRequests cqrs.Dispatcher[*in.ListIncomingFriendRequestsRequest, *out.ListIncomingFriendRequestsResponse]
	listOutgoingFriendRequests cqrs.Dispatcher[*in.ListOutgoingFriendRequestsRequest, *out.ListOutgoingFriendRequestsResponse]
	unfriendUser               cqrs.Dispatcher[*in.UnfriendUserRequest, *out.UnfriendUserResponse]
	listFriends                cqrs.Dispatcher[*in.ListFriendsRequest, *out.ListFriendsResponse]
	followUser                 cqrs.Dispatcher[*in.FollowUserRequest, *out.FollowUserResponse]
	unfollowUser               cqrs.Dispatcher[*in.UnfollowUserRequest, *out.UnfollowUserResponse]
	listFollowers              cqrs.Dispatcher[*in.ListFollowersRequest, *out.ListFollowersResponse]
	listFollowing              cqrs.Dispatcher[*in.ListFollowingRequest, *out.ListFollowingResponse]
	blockUser                  cqrs.Dispatcher[*in.BlockUserRequest, *out.BlockUserResponse]
	unblockUser                cqrs.Dispatcher[*in.UnblockUserRequest, *out.UnblockUserResponse]
	listBlockedUsers           cqrs.Dispatcher[*in.ListBlockedUsersRequest, *out.ListBlockedUsersResponse]
	getRelationshipStatus      cqrs.Dispatcher[*in.GetRelationshipStatusRequest, *out.GetRelationshipStatusResponse]
	getMutualFriends           cqrs.Dispatcher[*in.GetMutualFriendsRequest, *out.GetMutualFriendsResponse]
	getRelationshipSummary     cqrs.Dispatcher[*in.GetRelationshipSummaryRequest, *out.GetRelationshipSummaryResponse]
}

func NewHTTPServer(
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
) (infrahttp.HTTPServer, error) {
	return &relationshipHTTPServer{
		sendFriendRequest:          sendFriendRequest,
		cancelFriendRequest:        cancelFriendRequest,
		acceptFriendRequest:        acceptFriendRequest,
		rejectFriendRequest:        rejectFriendRequest,
		listIncomingFriendRequests: listIncomingFriendRequests,
		listOutgoingFriendRequests: listOutgoingFriendRequests,
		unfriendUser:               unfriendUser,
		listFriends:                listFriends,
		followUser:                 followUser,
		unfollowUser:               unfollowUser,
		listFollowers:              listFollowers,
		listFollowing:              listFollowing,
		blockUser:                  blockUser,
		unblockUser:                unblockUser,
		listBlockedUsers:           listBlockedUsers,
		getRelationshipStatus:      getRelationshipStatus,
		getMutualFriends:           getMutualFriends,
		getRelationshipSummary:     getRelationshipSummary,
	}, nil
}

func (s *relationshipHTTPServer) RegisterPublicRoutes(routes *gin.RouterGroup) {
	relationshiphttp.RegisterPublicRoutes(routes, s.sendFriendRequest, s.cancelFriendRequest, s.acceptFriendRequest, s.rejectFriendRequest, s.listIncomingFriendRequests, s.listOutgoingFriendRequests, s.unfriendUser, s.listFriends, s.followUser, s.unfollowUser, s.listFollowers, s.listFollowing, s.blockUser, s.unblockUser, s.listBlockedUsers, s.getRelationshipStatus, s.getMutualFriends, s.getRelationshipSummary)
}

func (s *relationshipHTTPServer) RegisterPrivateRoutes(routes *gin.RouterGroup) {
	relationshiphttp.RegisterPrivateRoutes(routes)
}

func (s *relationshipHTTPServer) RegisterSocketRoutes(routes *gin.RouterGroup) {
}

func (s *relationshipHTTPServer) Stop(ctx context.Context) error {
	return nil
}
