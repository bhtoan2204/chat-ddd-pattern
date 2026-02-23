package generator

import "fmt"

type modulePaths struct {
	FsRoot     string
	ImportRoot string
}

func moduleForUsecase(usecaseName string) (modulePaths, error) {
	switch usecaseName {
	case "AuthUsecase":
		return modulePaths{
			FsRoot:     "core/account",
			ImportRoot: "go-socket/core/modules/account",
		}, nil
	case "RoomUsecase", "MessageUsecase":
		return modulePaths{
			FsRoot:     "core/room",
			ImportRoot: "go-socket/core/modules/room",
		}, nil
	default:
		return modulePaths{}, fmt.Errorf("unknown usecase: %s", usecaseName)
	}
}
