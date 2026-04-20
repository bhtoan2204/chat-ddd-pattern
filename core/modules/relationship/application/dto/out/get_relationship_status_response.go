// CODE_GENERATOR - do not edit: response
package out

type GetRelationshipStatusResponse struct {
	IsSelf                       bool `json:"is_self,omitempty"`
	IsFriend                     bool `json:"is_friend,omitempty"`
	IsFollowing                  bool `json:"is_following,omitempty"`
	IsFollower                   bool `json:"is_follower,omitempty"`
	HasBlocked                   bool `json:"has_blocked,omitempty"`
	IsBlockedBy                  bool `json:"is_blocked_by,omitempty"`
	OutgoingFriendRequestPending bool `json:"outgoing_friend_request_pending,omitempty"`
	IncomingFriendRequestPending bool `json:"incoming_friend_request_pending,omitempty"`
	CanSendFriendRequest         bool `json:"can_send_friend_request,omitempty"`
	CanFollow                    bool `json:"can_follow,omitempty"`
}
