package constant

import (
	"time"

	"wechat-clone/core/shared/infra/lock"
)

const RealtimeMessageTopic = "room.realtime.message"

const VideoCallSessionTTL = 4 * time.Hour

func DefaultVideoCallLockOptions() lock.MultiLockOptions {
	opts := lock.DefaultMultiLockOptions()
	opts.KeyPrefix = "room:video_call"
	return opts
}
