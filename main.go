package main

import (
	"fmt"
	"os"
	"resource-monitor/proc"
	"time"
)

func main() {
	currentProcPid := os.Getpid()
	fmt.Printf("Current PID: %d\n", currentProcPid)

	for {
		select {
		case <-time.Tick(time.Second):
			stats, err := proc.ReadStat()
			if err != nil {
				panic(err)
			}

			fmt.Println(stats.UpTime)

			for _, stat := range stats.ProcStats {
				fmt.Printf("{%s %s %d %d %d}\n", stat.Name, stat.State, stat.UTime, stat.STime, stat.StartTime)
			}
			fmt.Println()
		}

	}
}
