package helpers

import (
	"os"
	"strings"
)

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func DeepDirectory(path string, targetDir string) (string, bool) {
	info, err := os.Stat(path + "/" + targetDir)
	if os.IsNotExist(err) {
		newPath, _, ok := CutLastDirectoryPath(path)
		if !ok {
			return "", false
		}
		return DeepDirectory(newPath, targetDir)
	}
	return path, info.IsDir()
}

func WalkDeepDirectory(path string, targetDir string) ([]string, bool) {
	info, err := os.Stat(path + "/" + targetDir)
	if os.IsNotExist(err) {
		newPath, _, ok := CutLastDirectoryPath(path)
		if !ok {
			return nil, false
		}

		dirs, ok := WalkDeepDirectory(newPath, targetDir)
		if !ok {
			return []string{}, false
		}

		return append(dirs, path), true
	}

	return []string{path}, info.IsDir()
}

func CutLastDirectoryPath(path string) (prefix string, cutted string, ok bool) {
	if strings.LastIndex(path, "/") == -1 {
		return "", "", false
	}

	return path[:strings.LastIndex(path, "/")], path[strings.LastIndex(path, "/")+1:], true
}
