package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUStats struct {
	User      int
	Nice      int
	System    int
	Idle      int
	Iowait    int
	Irq       int
	Softirq   int
	Steal     int
	TotalTime int
}

func readCPUStats() (map[string]CPUStats, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	cpuStats := make(map[string]CPUStats)

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			if len(fields) < 8 {
				continue
			}
			name := fields[0]
			values := make([]int, 0, len(fields)-1)
			for _, field := range fields[1:] {
				value, err := strconv.Atoi(field)
				if err != nil {
					return nil, err
				}
				values = append(values, value)
			}

			totalTime := sum(values)
			cpuStats[name] = CPUStats{
				User:      values[0],
				Nice:      values[1],
				System:    values[2],
				Idle:      values[3],
				Iowait:    values[4],
				Irq:       values[5],
				Softirq:   values[6],
				Steal:     values[7],
				TotalTime: totalTime,
			}
		}
	}

	return cpuStats, nil
}

func sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func calculateCPUUsage(stat1, stat2 CPUStats) float64 {
	totalDiff := stat2.TotalTime - stat1.TotalTime
	idleDiff := stat2.Idle - stat1.Idle

	if totalDiff == 0 {
		return 0.0
	}

	return 100.0 * float64(totalDiff-idleDiff) / float64(totalDiff)
}

func main() {
	for {
		stat1, err := readCPUStats()
		if err != nil {
			fmt.Printf("Error reading CPU stats: %v\n", err)
			return
		}

		time.Sleep(200 * time.Millisecond)

		stat2, err := readCPUStats()
		if err != nil {
			fmt.Printf("Error reading CPU stats: %v\n", err)
			return
		}

		fmt.Print("\033[H\033[2J")
		for cpu, stats1 := range stat1 {
			stats2 := stat2[cpu]
			usage := calculateCPUUsage(stats1, stats2)
			fmt.Printf("%s: %.2f%%\n", cpu, usage)
		}
	}
}
