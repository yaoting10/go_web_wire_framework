package util

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExecPath 当前执行文件所在的目录
func ExecPath() string {
	_path, err := os.Executable() // 获得程序路径
	if err != nil {
		panic(err)
	}
	return filepath.Dir(_path)
}

// CurrentPath 当前执行文件所在的目录
// func CurrentPath() string {
//	_, filename, _, _ := runtime.Caller(1)
//	return path.Dir(filename)
// }

// ProjectPath 获取工程根目录，首先按照模块工程获取，即工程根目录为 go.mod 所在的目录，找不到 go.mod，则返回 ExecPath 的值
func ProjectPath() string {
	// default linux/mac os
	var (
		sp = string(os.PathSeparator)
		ss []string
	)

	// GOMOD
	// in go source code:
	// Check for use of modules by 'go env GOMOD',
	// which reports a go.mod file path if modules are enabled.

	stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	p := string(bytes.TrimSpace(stdout))
	if p != "/dev/null" && p != "" {
		ss = strings.Split(p, sp)
		ss = ss[:len(ss)-1]
		return strings.Join(ss, sp)
	}
	return ExecPath()
}
