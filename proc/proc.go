package proc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	// https://man7.org/linux/man-pages/man5/proc.5.html
	procDir        = "/proc"
	procPIDStatFmt = "/proc/%d/stat"
	statFileSep    = " "
	procNameIdx    = 1
	procStateIdx   = 2
	procUTimeIdx   = 13
	procSTimeIdx   = 14
)

var ()

// ProcessStats represents the statistics of the given
// process in the /proc/{proc_pid}/stat file
type ProcessStats struct {
	Name  string
	State string
	UTime uint64 // The UTime value is the time the process has been running in user mode
	STime uint64 // The STime value is the amount of time the process has been running in kernel mode
}

// GetProcStat returns a slice with the parsed details of currently
// running processes
func GetProcStat() ([]*ProcessStats, error) {
	PIDs, err := readPIDs()
	if err != nil {
		return nil, err
	}

	res := make([]*ProcessStats, 0)
	lock := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for _, p := range PIDs {
		wg.Add(1)
		p := p
		go func() {
			numPID, err := strconv.Atoi(p)
			if err != nil {
				log.Fatal(err)
			}

			stats, err := getStatsForPID(numPID)
			if err != nil {
				log.Fatal(err)
			}

			lock.Lock()
			res = append(res, stats)
			lock.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()

	return res, nil
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
func getStatsForPID(PID int) (*ProcessStats, error) {
	rawStats, fErr := readPIDStatFile(PID)
	uTime, uErr := strconv.ParseUint(rawStats[procUTimeIdx], 10, 64)
	sTime, sErr := strconv.ParseUint(rawStats[procSTimeIdx], 10, 64)

	if err := errors.Join(fErr, uErr, sErr); err != nil {
		return nil, err
	}

	stats := ProcessStats{
		Name:  rawStats[procNameIdx],
		State: rawStats[procStateIdx],
		UTime: uTime,
		STime: sTime,
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
