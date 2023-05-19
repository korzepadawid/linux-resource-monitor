package proc

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	procDir           = "/proc"
	procUpTimeFile    = "/proc/uptime"
	procUpTimeFileSep = " "
)

type Info struct {
	ProcStats []*Stats
	UpTime    float64
}

// getProcUpTime returns the system's uptime in seconds
// (uses /proc/uptime)
func getProcUpTime() (float64, error) {
	bytes, err := os.ReadFile(procUpTimeFile)
	if err != nil {
		return 0, err
	}

	rawUpTimeFileContent := string(bytes)
	rawUpTime := strings.Split(rawUpTimeFileContent, procUpTimeFileSep)[0]

	upTimeNum, err := strconv.ParseFloat(rawUpTime, 64)
	if err != nil {
		return 0, err
	}

	return upTimeNum, nil
}

// ReadStat returns a slice with the parsed details of currently
// running processes
func ReadStat() (*Info, error) {
	PIDs, pidErr := readPIDs()
	upTime, upTimeErr := getProcUpTime()

	if err := errors.Join(pidErr, upTimeErr); err != nil {
		return nil, err
	}

	procStats := make([]*Stats, 0)
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
			procStats = append(procStats, stats)
			lock.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()

	info := &Info{
		ProcStats: procStats,
		UpTime:    upTime,
	}

	return info, nil
}
