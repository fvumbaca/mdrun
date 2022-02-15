package browser

import (
	"io/fs"
	"path/filepath"
)

type Handler struct {
}

var filterForExtensions = []string{".md"}

func listMarkdownAndDirs(root fs.FS) ([]string, error) {
	var filenames []string
	infos, err := fs.ReadDir(root, ".")
	if err != nil {
		return filenames, err
	}
	for _, i := range infos {
		if i.IsDir() || strIn(filterForExtensions, filepath.Ext(i.Name())) {
			filenames = append(filenames, i.Name())
		}
	}
	return filenames, nil
}

func strIn(list []string, s string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}
