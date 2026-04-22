//go:build windows
// +build windows

package gocmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// ProcessExist 检查进程是否存在（Windows 实现）
func ProcessExist(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// 在 Windows 上，FindProcess 成功通常意味着进程存在
	process.Release()
	return true
}

// QueryProcess Windows 实现
// 自动判断系统版本：在新版 Windows 上使用 PowerShell，在旧版上使用 wmic
func QueryProcess(name string) []*ProcessInfo {
	// 去掉可能存在的 .exe 后缀，因为查询逻辑会统一处理
	cleanName := strings.TrimSuffix(name, ".exe")

	v, _ := syscall.GetVersion()
	major := uint8(v)
	// Windows 10/Server 2016 及以上版本 major 为 10
	if major >= 10 {
		return queryProcessWithPowerShell(cleanName)
	}
	return queryProcessWithWmic(name)
}

// 使用 PowerShell 获取进程信息 (推荐用于 Win10/Server 2016+)
func queryProcessWithPowerShell(name string) []*ProcessInfo {
	pi := make([]*ProcessInfo, 0)
	// Get-CimInstance Win32_Process 比 Get-Process 拿 CommandLine 更可靠
	psCmd := fmt.Sprintf("Get-CimInstance Win32_Process -Filter \"name like '%s.exe'\" | Select-Object ProcessId, Name, CommandLine | ConvertTo-Csv -NoTypeInformation", name)
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCmd)
	output, err := cmd.Output()
	if err != nil {
		return pi
	}

	lines := strings.Split(string(output), "\r\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ",") {
			continue
		}
		// CSV: "ProcessId","Name","CommandLine"
		parts := strings.Split(strings.ReplaceAll(line, "\"", ""), ",")
		if len(parts) < 3 {
			continue
		}

		pid, _ := strconv.Atoi(parts[0])
		if pid == 0 {
			continue
		}

		pi = append(pi, &ProcessInfo{
			Pid:     pid,
			Name:    parts[1],
			CmdLine: parts[2],
		})
	}
	return pi
}

// 使用 wmic 获取进程信息 (用于旧版系统)
func queryProcessWithWmic(name string) []*ProcessInfo {
	pi := make([]*ProcessInfo, 0)
	// 确保查询名带有 .exe 后缀以便于 wmic 精确匹配
	searchName := name
	if !strings.HasSuffix(strings.ToLower(name), ".exe") {
		searchName = name + ".exe"
	}

	cmd := exec.Command("wmic", "process", "where", fmt.Sprintf("name='%s'", searchName), "get", "CommandLine,ProcessId,Name", "/format:csv")
	output, err := cmd.Output()
	if err != nil {
		return pi
	}

	lines := strings.Split(string(output), "\r\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(strings.ToLower(line), "node") {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 4 {
			continue
		}
		// wmic csv: Node,CommandLine,Name,ProcessId
		cmdLine := parts[1]
		procName := parts[2]
		pidStr := parts[3]

		pid, _ := strconv.Atoi(pidStr)
		if pid == 0 {
			continue
		}

		pi = append(pi, &ProcessInfo{
			Name:    procName,
			Pid:     pid,
			CmdLine: cmdLine,
		})
	}
	return pi
}
