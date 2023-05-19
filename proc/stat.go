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
	Name      string
	State     string
	StartTime uint64 // The StartTime value is the process’s start time
	UTime     uint64 // The UTime value is the time the process has been running in user mode
	STime     uint64 // The STime value is the amount of time the process has been running in kernel mode
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
func getStatsForPID(PID int) (*Stats, error) {
	rawStats, fErr := readPIDStatFile(PID)
	uTime, uErr := strconv.ParseUint(rawStats[procUTimeIdx], 10, 64)
	sTime, sErr := strconv.ParseUint(rawStats[procSTimeIdx], 10, 64)
	startTime, stErr := strconv.ParseUint(rawStats[procStartTimeIdx], 10, 64)

	if err := errors.Join(fErr, uErr, sErr, stErr); err != nil {
		return nil, err
	}

	stats := Stats{
		Name:      rawStats[procNameIdx],
		State:     rawStats[procStateIdx],
		UTime:     uTime,
		STime:     sTime,
		StartTime: startTime,
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
