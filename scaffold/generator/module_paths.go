package generator

import (
	"strings"

	"wechat-clone/scaffold/utils"
)

type modulePaths struct {
	FsRoot     string
	ImportRoot string
}

func moduleForUsecase(usecaseName string) (modulePaths, error) {
	moduleName := strings.TrimSuffix(usecaseName, "Usecase")
	moduleDir := resolveModuleDir(moduleName)
	return modulePaths{
		FsRoot:     "core/modules/" + moduleDir,
		ImportRoot: "wechat-clone/core/modules/" + moduleDir,
	}, nil
}

func resolveModuleDir(moduleName string) string {
	switch moduleName {
	case "Auth":
		return "account"
	case "Message":
		return "room"
	default:
		return utils.Snake(moduleName)
	}
}
