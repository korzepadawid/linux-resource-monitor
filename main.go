package main

import (
	"fmt"
	"log"
	"os"
	"resource-monitor/proc"
	"strconv"
	"time"
)

func main() {
	currentProcPid := os.Getpid()
	fmt.Printf("Current PID: %d\n", currentProcPid)

	for {
		select {
		case <-time.Tick(2 * time.Second):
			PIDs, err := proc.ReadPIDs()
			if err != nil {
				log.Fatal(err)
			}

			for _, PID := range PIDs {
				fmt.Println(PID)
				numPID, err := strconv.Atoi(PID)
				if err != nil {
					log.Fatal(err)
				}

				stats, err := proc.GetStatsForPID(numPID)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(stats)
				fmt.Printf("%s %s")
				fmt.Println()
			}

		}
	}
}
