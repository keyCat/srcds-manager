package steamcmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func UpdateForPath(path string) {
}

func resolveAppManifestPath(path string) (string, error) {
	fsys := os.DirFS(path)
	matches, err := fs.Glob(fsys, "steamapps/appmanifest_*.acf")
	if err != nil {
		return "", err
	}
	if matches == nil {
		return "", errors.New(fmt.Sprintf("App Manifset not found by glob %s/steamapps/appmanifest_*.acf", path))
	}
	return matches[0], nil
}
