package io

import (
	"github.com/hankmor/gotools/errs"
	"goboot/cmd/cli/util"
	"os"
	"path/filepath"
	"strings"
)

func WalkAllFiles(path string) []string {
	fi, err := os.Stat(path)
	errs.Throw(err)

	if fi.IsDir() {
		root := util.FmtDirPath(path)
		return walkDir(root)
	} else {
		return []string{path}
	}
}

func walkDir(dir string) []string {
	dir = util.FmtDirPath(dir)
	// 读取文件
	es, err := os.ReadDir(dir)
	errs.Throw(err)
	var files []string
	for _, entry := range es {
		if entry.IsDir() {
			files = append(files, walkDir(dir+entry.Name())...)
		} else {
			fi, _ := entry.Info()
			if strings.ToLower(filepath.Ext(entry.Name())) != ".txt" {
				continue
			}
			file := dir + fi.Name()
			files = append(files, file)
		}
	}
	return files
}
