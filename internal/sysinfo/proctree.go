package sysinfo

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ProcessTreeRSS returns the resident memory in bytes of the process root plus every process
// started below it (its children, grandchildren, and so on). That is the memory Orbit uses:
// the agent and the apps and databases it launched. It returns 0 when /proc is unavailable.
func ProcessTreeRSS(root int) uint64 {
	children := map[int][]int{}
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return 0
	}
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		if ppid, ok := readPPID(pid); ok {
			children[ppid] = append(children[ppid], pid)
		}
	}

	var total uint64
	queue := []int{root}
	seen := map[int]bool{}
	for len(queue) > 0 {
		pid := queue[0]
		queue = queue[1:]
		if seen[pid] {
			continue
		}
		seen[pid] = true
		total += readRSS(pid)
		queue = append(queue, children[pid]...)
	}
	return total
}

// readPPID returns the parent pid from /proc/<pid>/stat. The command name is in parentheses and
// can contain spaces, so the fields are read after the last ')'.
func readPPID(pid int) (int, bool) {
	b, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return 0, false
	}
	s := string(b)
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return 0, false
	}
	fields := strings.Fields(s[i+1:])
	if len(fields) < 2 {
		return 0, false
	}
	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, false
	}
	return ppid, true
}

// readRSS returns VmRSS for a process in bytes, or 0 if it has already exited.
func readRSS(pid int) uint64 {
	f, err := os.Open("/proc/" + strconv.Itoa(pid) + "/status")
	if err != nil {
		return 0
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) >= 2 && fields[0] == "VmRSS:" {
			kb, _ := strconv.ParseUint(fields[1], 10, 64)
			return kb * 1024
		}
	}
	return 0
}
