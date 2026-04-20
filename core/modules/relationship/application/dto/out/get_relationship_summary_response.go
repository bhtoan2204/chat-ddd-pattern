// CODE_GENERATOR - do not edit: response
package out

type GetRelationshipSummaryResponse struct {
	FriendsCount       int64  `json:"friends_count,omitempty"`
	FollowersCount     int64  `json:"followers_count,omitempty"`
	FollowingCount     int64  `json:"following_count,omitempty"`
	MutualFriendsCount int64  `json:"mutual_friends_count,omitempty"`
	RelationshipStatus string `json:"relationship_status,omitempty"`
}
