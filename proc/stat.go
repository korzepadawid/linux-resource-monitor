package proc

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	procPIDStatFmt   = "/proc/%d/stat"
	statFileSep      = " "
	procNameIdx      = 1
	procStateIdx     = 2
	procUTimeIdx     = 13
	procSTimeIdx     = 14
	procStartTimeIdx = 21
)

// Stats represents the statistics of the given
// process in the /proc/{proc_pid}/stat file
type Stats struct {
	Name      string  // The Name value is the name of the process
	State     string  // The State value is the state of the process
	StartTime uint64  // The StartTime value is the process’s start time
	UTime     uint64  // The UTime value is the time the process has been running in user mode
	STime     uint64  // The STime value is the amount of time the process has been running in kernel mode
	UpTime    float64 // The UpTime value is the system’s uptime
	ClkTck    float64 // The ClkTck value is the number of clock ticks in a second
}

// CPUUsage returns the CPU usage of the given process in %
func (s Stats) CPUUsage() float64 {
	uTimeInSec := float64(s.UTime) / s.ClkTck
	sTimeInSec := float64(s.STime) / s.ClkTck
	startTimeInSec := float64(s.StartTime) / s.ClkTck

	elapsedInSec := s.UpTime - startTimeInSec
	usageInSec := sTimeInSec + uTimeInSec
	return usageInSec * 100 / elapsedInSec
}

// readPIDs reads PIDs of currently running processes.
// (the function uses the /proc/ directory)
func readPIDs() ([]string, error) {
	regexPID := regexp.MustCompile("^[0-9]*$")
	dirEntries, err := os.ReadDir(procDir)
	if err != nil {
		return nil, err
	}

	dirNamesPID := make([]string, 0)
	for _, entry := range dirEntries {
		if entry.IsDir() && regexPID.MatchString(entry.Name()) {
			dirNamesPID = append(dirNamesPID, entry.Name())
		}
	}

	return dirNamesPID, nil
}

// getStatsForPID reads stats from the /proc/{proc_pid}/stat file.
func getStatsForPID(PID int, upTime float64, clkTck float64) (*Stats, error) {
	rawStats, fileErr := readPIDStatFile(PID)
	uTime, uTimeErr := strconv.ParseUint(rawStats[procUTimeIdx], 10, 64)
	sTime, sTimeErr := strconv.ParseUint(rawStats[procSTimeIdx], 10, 64)
	startTime, startTimeErr := strconv.ParseUint(rawStats[procStartTimeIdx], 10, 64)

	if err := errors.Join(
		fileErr,
		uTimeErr,
		sTimeErr,
		startTimeErr,
	); err != nil {
		return nil, err
	}

	stats := Stats{
		Name:      rawStats[procNameIdx][1 : len(rawStats[procNameIdx])-1],
		State:     rawStats[procStateIdx],
		StartTime: startTime,
		UTime:     uTime,
		STime:     sTime,
		UpTime:    upTime,
		ClkTck:    clkTck,
	}

	return &stats, nil
}

func readPIDStatFile(PID int) ([]string, error) {
	procPIDStatFilePath := fmt.Sprintf(procPIDStatFmt, PID)
	fileBytes, err := os.ReadFile(procPIDStatFilePath)
	if err != nil {
		return nil, err
	}

	line := string(fileBytes)
	rawStats := strings.Split(line, statFileSep)
	return rawStats, nil
}
