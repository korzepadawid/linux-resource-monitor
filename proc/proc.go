package proc

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	// https://man7.org/linux/man-pages/man5/proc.5.html
	procDir        = "/proc"
	procPIDStatFmt = "/proc/%d/stat"
	statFileSep    = " "
	procNameIdx    = 1
	procStateIdx   = 2
	procUtimeIdx   = 13
	procStimeIdx   = 14
)

var (
	regexPID = regexp.MustCompile("^[0-9]*$")
)

// ReadPIDs reads PIDs of currently running processes.
// (the function uses the /proc/ directory)
func ReadPIDs() ([]string, error) {
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

type ProcessStats struct {
	Name  string
	State string
	UTime uint64
	STime uint64
}

func GetStatsForPID(PID int) (*ProcessStats, error) {
	rawStats, err := readPIDStatFile(PID)
	if err != nil {
		log.Fatal(err)
	}

	utime, err := strconv.ParseUint(rawStats[procUtimeIdx], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	stime, err := strconv.ParseUint(rawStats[procStimeIdx], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	stats := ProcessStats{
		Name:  rawStats[procNameIdx],
		State: rawStats[procStateIdx],
		UTime: utime,
		STime: stime,
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
