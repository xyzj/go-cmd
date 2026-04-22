//go:build !windows
// +build !windows

package gocmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ProcessExist 检查进程是否存在（Unix/Linux 实现）
func ProcessExist(pid int) bool {
	err := syscall.Kill(pid, 0)

	if err == nil {
		return true // 进程存在且你有权限
	}

	if err == syscall.EPERM {
		return true // 进程存在，但你没权限
	}

	return false
}

// QueryProcess only for linux
func QueryProcess(name string) []*ProcessInfo {
	pi := make([]*ProcessInfo, 0)
	procs, err := os.ReadDir("/proc")
	if err != nil {
		return pi
	}
	for _, proc := range procs {
		if !proc.IsDir() {
			continue
		}
		pid, _ := strconv.ParseInt(proc.Name(), 10, 32)
		if pid == 0 {
			continue
		}
		cmd, _ := os.ReadFile("/proc/" + proc.Name() + "/cmdline")
		if len(cmd) == 0 {
			continue
		}
		cl := strings.Split(string(cmd), "\x00")
		if name != filepath.Base(cl[0]) {
			continue
		}
		pi = append(pi, &ProcessInfo{
			Name:    name,
			Pid:     int(pid),
			CmdLine: strings.Join(cl, " "),
		})
	}
	return pi
}
