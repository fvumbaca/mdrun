package browser

import (
	"io/ioutil"
	"path/filepath"
)

var filterForExtensions = []string{".md"}

func listMarkdownAndDirs(root string) ([]string, error) {
	var filenames []string
	infos, err := ioutil.ReadDir(root)
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
