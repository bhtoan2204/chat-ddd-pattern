package generator

import (
	"fmt"
	"strings"
)

type modulePaths struct {
	FsRoot     string
	ImportRoot string
}

func moduleForUsecase(usecaseName string) (modulePaths, error) {
	moduleName := strings.TrimSuffix(usecaseName, "Usecase")
	switch moduleName {
	case "Auth":
		return modulePaths{
			FsRoot:     "core/modules/account",
			ImportRoot: "wechat-clone/core/modules/account",
		}, nil
	case "Room", "Message":
		return modulePaths{
			FsRoot:     "core/modules/room",
			ImportRoot: "wechat-clone/core/modules/room",
		}, nil
	case "Notification":
		return modulePaths{
			FsRoot:     "core/modules/notification",
			ImportRoot: "wechat-clone/core/modules/notification",
		}, nil
	case "Payment":
		return modulePaths{
			FsRoot:     "core/modules/payment",
			ImportRoot: "wechat-clone/core/modules/payment",
		}, nil
	case "Ledger":
		return modulePaths{
			FsRoot:     "core/modules/ledger",
			ImportRoot: "wechat-clone/core/modules/ledger",
		}, nil
	default:
		return modulePaths{}, fmt.Errorf("unknown usecase: %s", usecaseName)
	}
}
