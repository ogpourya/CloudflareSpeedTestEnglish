package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultOutput         = "result.csv"
	maxDelay              = 9999 * time.Millisecond
	minDelay              = 0 * time.Millisecond
	maxLossRate   float32 = 1.0
)

var (
	InputMaxDelay    = maxDelay
	InputMinDelay    = minDelay
	InputMaxLossRate = maxLossRate
	Output           = defaultOutput
	PrintNum         = 10
	Debug            = false // Whether to enable debug mode
)

// Whether to print test results
func NoPrintResult() bool {
	return PrintNum == 0
}

// Whether to output to file
func noOutput() bool {
	return Output == "" || Output == " "
}

type PingData struct {
	IP       *net.IPAddr
	Sended   int
	Received int
	Delay    time.Duration
	Colo     string
}

type CloudflareIPData struct {
	*PingData
	lossRate      float32
	DownloadSpeed float64
}

// Calculate packet loss rate
func (cf *CloudflareIPData) getLossRate() float32 {
	if cf.lossRate == 0 {
		pingLost := cf.Sended - cf.Received
		cf.lossRate = float32(pingLost) / float32(cf.Sended)
	}
	return cf.lossRate
}

func (cf *CloudflareIPData) toString() []string {
	result := make([]string, 7)
	result[0] = cf.IP.String()
	result[1] = strconv.Itoa(cf.Sended)
	result[2] = strconv.Itoa(cf.Received)
	result[3] = strconv.FormatFloat(float64(cf.getLossRate()), 'f', 2, 32)
	result[4] = strconv.FormatFloat(cf.Delay.Seconds()*1000, 'f', 2, 32)
	result[5] = strconv.FormatFloat(cf.DownloadSpeed/1024/1024, 'f', 2, 32)
	// If Colo is empty, use "N/A"
	if cf.Colo == "" {
		result[6] = "N/A"
	} else {
		result[6] = cf.Colo
	}
	return result
}

func ExportCsv(data []CloudflareIPData) {
	if noOutput() || len(data) == 0 {
		return
	}
	fp, err := os.Create(Output)
	if err != nil {
		log.Fatalf("Failed to create file [%s]: %v", Output, err)
		return
	}
	defer fp.Close()
	w := csv.NewWriter(fp) // Create a new file write stream
	_ = w.Write([]string{"IP Address", "Sent", "Received", "Loss Rate", "Avg Latency", "Download Speed(MB/s)", "Region Code"})
	_ = w.WriteAll(convertToString(data))
	w.Flush()
}

func convertToString(data []CloudflareIPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.toString())
	}
	return result
}

// Sort by delay and packet loss
type PingDelaySet []CloudflareIPData

// Filter by delay condition
func (s PingDelaySet) FilterDelay() (data PingDelaySet) {
	if InputMaxDelay > maxDelay || InputMinDelay < minDelay { // Skip filtering when delay conditions are outside the default range
		return s
	}
	if InputMaxDelay == maxDelay && InputMinDelay == minDelay { // Skip filtering when delay conditions are at default values
		return s
	}
	for _, v := range s {
		if v.Delay > InputMaxDelay { // Upper limit: delay exceeds max, no subsequent data will satisfy, break out
			break
		}
		if v.Delay < InputMinDelay { // Lower limit: delay below min, skip this entry
			continue
		}
		data = append(data, v) // Delay satisfies condition, add to new slice
	}
	return
}

// Filter by packet loss condition
func (s PingDelaySet) FilterLossRate() (data PingDelaySet) {
	if InputMaxLossRate >= maxLossRate { // Skip filtering when loss rate condition is at default
		return s
	}
	for _, v := range s {
		if v.getLossRate() > InputMaxLossRate { // Upper limit for packet loss rate
			break
		}
		data = append(data, v) // Loss rate satisfies condition, add to new slice
	}
	return
}

func (s PingDelaySet) Len() int {
	return len(s)
}
func (s PingDelaySet) Less(i, j int) bool {
	iRate, jRate := s[i].getLossRate(), s[j].getLossRate()
	if iRate != jRate {
		return iRate < jRate
	}
	return s[i].Delay < s[j].Delay
}
func (s PingDelaySet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort by download speed
type DownloadSpeedSet []CloudflareIPData

func (s DownloadSpeedSet) Len() int {
	return len(s)
}
func (s DownloadSpeedSet) Less(i, j int) bool {
	return s[i].DownloadSpeed > s[j].DownloadSpeed
}
func (s DownloadSpeedSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DownloadSpeedSet) Print() {
	if NoPrintResult() {
		return
	}
	if len(s) <= 0 { // Continue only if IP array length > 0
		fmt.Println("\n[Info] Total IP count in full test results is 0, skipping output.")
		return
	}
	dateString := convertToString(s) // Convert to multi-dimensional array [][]String
	if len(dateString) < PrintNum {  // If IP count is less than PrintNum, use IP count instead
		PrintNum = len(dateString)
	}
	headFormat := "%-16s%-5s%-5s%-5s%-6s%-12s%-5s\n"
	dataFormat := "%-18s%-8s%-8s%-8s%-10s%-16s%-8s\n"
	for i := 0; i < PrintNum; i++ { // If the output IPs contain IPv6, adjust spacing
		if len(dateString[i][0]) > 15 {
			headFormat = "%-40s%-5s%-5s%-5s%-6s%-12s%-5s\n"
			dataFormat = "%-42s%-8s%-8s%-8s%-10s%-16s%-8s\n"
			break
		}
	}
	Cyan.Printf(headFormat, "IP Address", "Sent", "Received", "Loss Rate", "Avg Latency", "Download Speed(MB/s)", "Region Code")
	for i := 0; i < PrintNum; i++ {
		fmt.Printf(dataFormat, dateString[i][0], dateString[i][1], dateString[i][2], dateString[i][3], dateString[i][4], dateString[i][5], dateString[i][6])
	}
	if !noOutput() {
		fmt.Printf("\nFull test results have been written to %v file. You can view it with a text editor or spreadsheet software.\n", Output)
	}
}

