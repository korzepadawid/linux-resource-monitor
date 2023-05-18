package proc

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	// https://man7.org/linux/man-pages/man5/proc.5.html
	procDir        = "/proc"
	procPIDStatFmt = "/proc/%d/stat"
	statFileSep    = " "
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

func ParsePIDStatFile(PID int) ([]string, error) {
	procPIDStatFilePath := fmt.Sprintf(procPIDStatFmt, PID)
	fileBytes, err := os.ReadFile(procPIDStatFilePath)
	if err != nil {
		return nil, err
	}

	line := string(fileBytes)
	rawStats := strings.Split(line, statFileSep)
	return rawStats, nil
}
