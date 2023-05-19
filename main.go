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

			for _, s := range stats {
				fmt.Printf("{%s, CPU=%f}\n", s.Name, s.CPUUsage())
			}
			fmt.Println()
		}

	}
}
