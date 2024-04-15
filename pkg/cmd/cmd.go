package cmd

import (
	"goboot/pkg/printer"
	"os/exec"
)

func Exec(shell string, args ...string) string {
	cmd := exec.Command(shell, args...)
	printer.Printwln(8, "正在执行脚本: %s", cmd.String())
	bs, err := cmd.CombinedOutput()
	if err != nil {
		printer.Printwln(8, "执行脚本出错, 脚本: %s, 错误: %v", cmd.String(), err)
		return string(bs)
	}
	return string(bs)
}
