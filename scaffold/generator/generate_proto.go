package generator

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func GenerateProto(protoRoot string) (string, error) {
	files, err := collectProtoFiles(protoRoot)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "generated 0 proto file(s)", nil
	}
	sort.Strings(files)

	args := []string{
		"-I", protoRoot,
		"--go_out=.",
		"--go_opt=module=wechat-clone",
		"--go-grpc_out=.",
		"--go-grpc_opt=module=wechat-clone",
	}
	args = append(args, files...)

	cmd := exec.Command("protoc", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("protoc failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return fmt.Sprintf("generated %d proto file(s)", len(files)), nil
}

func collectProtoFiles(root string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".proto" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return files, nil
}
