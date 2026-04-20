// CODE_GENERATOR: module-http-server-builder
package assembly

import (
	"context"

	appCtx "wechat-clone/core/context"
	relationshipcommand "wechat-clone/core/modules/relationship/application/command"
	relationshipquery "wechat-clone/core/modules/relationship/application/query"
	relationshiprepo "wechat-clone/core/modules/relationship/infra/persistent/repository"
	relationshipserver "wechat-clone/core/modules/relationship/transport/server"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/transport/http"
)

func buildHTTPServer(_ context.Context, appContext *appCtx.AppContext) (http.HTTPServer, error) {
	relationshipRepos := relationshiprepo.NewRepoImpl(appContext)
	sendFriendRequest := cqrs.NewDispatcher(relationshipcommand.NewSendFriendRequest(appContext, relationshipRepos))
	cancelFriendRequest := cqrs.NewDispatcher(relationshipcommand.NewCancelFriendRequest(appContext, relationshipRepos))
	acceptFriendRequest := cqrs.NewDispatcher(relationshipcommand.NewAcceptFriendRequest(appContext, relationshipRepos))
	rejectFriendRequest := cqrs.NewDispatcher(relationshipcommand.NewRejectFriendRequest(appContext, relationshipRepos))
	listIncomingFriendRequests := cqrs.NewDispatcher(relationshipquery.NewListIncomingFriendRequests(appContext, relationshipRepos))
	listOutgoingFriendRequests := cqrs.NewDispatcher(relationshipquery.NewListOutgoingFriendRequests(appContext, relationshipRepos))
	unfriendUser := cqrs.NewDispatcher(relationshipcommand.NewUnfriendUser(appContext, relationshipRepos))
	listFriends := cqrs.NewDispatcher(relationshipquery.NewListFriends(appContext, relationshipRepos))
	followUser := cqrs.NewDispatcher(relationshipcommand.NewFollowUser(appContext, relationshipRepos))
	unfollowUser := cqrs.NewDispatcher(relationshipcommand.NewUnfollowUser(appContext, relationshipRepos))
	listFollowers := cqrs.NewDispatcher(relationshipquery.NewListFollowers(appContext, relationshipRepos))
	listFollowing := cqrs.NewDispatcher(relationshipquery.NewListFollowing(appContext, relationshipRepos))
	blockUser := cqrs.NewDispatcher(relationshipcommand.NewBlockUser(appContext, relationshipRepos))
	unblockUser := cqrs.NewDispatcher(relationshipcommand.NewUnblockUser(appContext, relationshipRepos))
	listBlockedUsers := cqrs.NewDispatcher(relationshipquery.NewListBlockedUsers(appContext, relationshipRepos))
	getRelationshipStatus := cqrs.NewDispatcher(relationshipquery.NewGetRelationshipStatus(appContext, relationshipRepos))
	getMutualFriends := cqrs.NewDispatcher(relationshipquery.NewGetMutualFriends(appContext, relationshipRepos))
	getRelationshipSummary := cqrs.NewDispatcher(relationshipquery.NewGetRelationshipSummary(appContext, relationshipRepos))

	server, err := relationshipserver.NewHTTPServer(
		sendFriendRequest,
		cancelFriendRequest,
		acceptFriendRequest,
		rejectFriendRequest,
		listIncomingFriendRequests,
		listOutgoingFriendRequests,
		unfriendUser,
		listFriends,
		followUser,
		unfollowUser,
		listFollowers,
		listFollowing,
		blockUser,
		unblockUser,
		listBlockedUsers,
		getRelationshipStatus,
		getMutualFriends,
		getRelationshipSummary,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return server, nil
}
