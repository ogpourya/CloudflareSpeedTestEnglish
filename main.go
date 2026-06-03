package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ogpourya/CloudflareSpeedTestEnglish/task"
	"github.com/ogpourya/CloudflareSpeedTestEnglish/utils"
)

var (
	version, versionNew string
)

func init() {
	var printVersion bool
	var help = `
CloudflareSpeedTest ` + version + `
Test the latency and speed of all IPs of various CDNs or websites to get the fastest IP (IPv4+IPv6)!
https://github.com/ogpourya/CloudflareSpeedTestEnglish

Parameters:
    -n 200
        Latency test threads; the more threads, the faster the latency test. Do not set too high on weak devices (such as routers); (default 200, max 1000)
    -t 4
        Latency test times; number of latency tests for a single IP; (default 4)
    -dn 10
        Download test quantity; number of IPs to test download speed from the lowest latency up after sorting; (default 10)
    -dt 10
        Download test duration; maximum time for download speed test of a single IP, should not be too short; (default 10 seconds)
    -tp 443
        Specify test port; port used during latency/download test; (default 443)
    -url https://cf.xiu2.xyz/url
        Specify test URL; URL used during latency test (HTTPing)/download test, the default URL is not guaranteed to be available, self-hosting is recommended;

    -httping
        Switch test mode; change latency test mode to HTTP protocol, using the [-url] parameter; (default TCPing)
    -httping-code 200
        Valid HTTP status code; the valid HTTP status code returned by the webpage during HTTPing, only one allowed; (default 200 301 302)
    -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
        Match specified locations; IATA airport codes, country or city codes, separated by English commas, only available in HTTPing mode; (default all locations)

    -tl 200
        Average latency upper limit; only output IPs lower than the specified average latency, limits can be used in combination; (default 9999 ms)
    -tll 40
        Average latency lower limit; only output IPs higher than the specified average latency; (default 0 ms)
    -tlr 0.2
        Packet loss rate upper limit; only output IPs lower than/equal to the specified packet loss rate, range 0.00~1.00, 0 filters out any IP with packet loss; (default 1.00)
    -sl 5
        Download speed lower limit; only output IPs higher than the specified download speed, testing stops once the target number [-dn] is met; (default 0.00 MB/s)

    -p 10
        Display result quantity; display specified number of results after speed test, if 0, exit directly without displaying results; (default 10)
    -f ip.txt
        IP range data file; enclose in quotes if the path contains spaces; supports other CDN IP ranges; (default ip.txt)
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        Specify IP range data; specify IP range data directly via parameters, separated by English commas; (default empty)
    -o result.csv
        Write to result file; enclose in quotes if the path contains spaces; if empty, do not write to file [-o ""]; (default result.csv)

    -dd
        Disable download test; when disabled, results are sorted by latency (default sorted by download speed); (default enabled)
    -allip
        Test all IPs; test every IP in the IP range (only supports IPv4); (default random IP test per /24 range)

    -debug
        Debug output mode; outputs more logs in unexpected situations to help determine the cause; (default disabled)

    -v
        Print program version + check for updates
    -h
        Print help instructions
`
	var minDelay, maxDelay, downloadTime int
	var maxLossRate float64
	flag.IntVar(&task.Routines, "n", 200, "Latency test threads")
	flag.IntVar(&task.PingTimes, "t", 4, "Latency test times")
	flag.IntVar(&task.TestCount, "dn", 10, "Download test quantity")
	flag.IntVar(&downloadTime, "dt", 10, "Download test duration")
	flag.IntVar(&task.TCPPort, "tp", 443, "Specify test port")
	flag.StringVar(&task.URL, "url", "https://cf.xiu2.xyz/url", "Specify test URL")

	flag.BoolVar(&task.Httping, "httping", false, "Switch test mode")
	flag.IntVar(&task.HttpingStatusCode, "httping-code", 0, "Valid HTTP status code")
	flag.StringVar(&task.HttpingCFColo, "cfcolo", "", "Match specified locations")

	flag.IntVar(&maxDelay, "tl", 9999, "Average latency upper limit")
	flag.IntVar(&minDelay, "tll", 0, "Average latency lower limit")
	flag.Float64Var(&maxLossRate, "tlr", 1, "Packet loss rate upper limit")
	flag.Float64Var(&task.MinSpeed, "sl", 0, "Download speed lower limit")

	flag.IntVar(&utils.PrintNum, "p", 10, "Display result quantity")
	flag.StringVar(&task.IPFile, "f", "ip.txt", "IP range data file")
	flag.StringVar(&task.IPText, "ip", "", "Specify IP range data")
	flag.StringVar(&utils.Output, "o", "result.csv", "Output result file")

	flag.BoolVar(&task.Disable, "dd", false, "Disable download test")
	flag.BoolVar(&task.TestAll, "allip", false, "Test all IPs")

	flag.BoolVar(&utils.Debug, "debug", false, "Debug output mode")

	flag.BoolVar(&printVersion, "v", false, "Print program version")
	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()

	if task.MinSpeed > 0 && time.Duration(maxDelay)*time.Millisecond == utils.InputMaxDelay {
		utils.Yellow.Println("[Tip] When using the [-sl] parameter, it is recommended to pair it with the [-tl] parameter to avoid continuous testing if the target number [-dn] cannot be met...")
	}
	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)
	task.Timeout = time.Duration(downloadTime) * time.Second
	task.HttpingCFColomap = task.MapColoMap()

	if printVersion {
		println(version)
		fmt.Println("Checking for updates...")
		checkUpdate()
		if versionNew != "" {
			utils.Yellow.Printf("*** New version [%s] found! Please go to [https://github.com/ogpourya/CloudflareSpeedTestEnglish] to update! ***", versionNew)
		} else {
			utils.Green.Println("Current version [" + version + "] is up to date!")
		}
		os.Exit(0)
	}
}

func main() {
	task.InitRandSeed() // Set random seed

	fmt.Printf("# XIU2/CloudflareSpeedTest %s \n\n", version)

	// Start latency test + filter latency/packet loss
	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	// Start download test
	speedData := task.TestDownloadSpeed(pingData)
	utils.ExportCsv(speedData) // Output file
	speedData.Print()          // Print results
	endPrint()                 // Choose exit method based on situation (for Windows)
}

// Choose exit method based on situation (for Windows)
func endPrint() {
	if utils.NoPrintResult() { // If no need to print speed test results, exit directly
		return
	}
	if runtime.GOOS == "windows" { // If it is Windows, require pressing Enter or Ctrl+C to exit (to prevent direct closing when run by double-clicking after the test finishes)
		fmt.Printf("Press Enter or Ctrl+C to exit.")
		fmt.Scanln()
	}
}

// Check for updates
func checkUpdate() {
	timeout := 10 * time.Second
	client := http.Client{Timeout: timeout}
	res, err := client.Get("https://api.xiu2.xyz/ver/cloudflarespeedtest.txt")
	if err != nil {
		return
	}
	// Read resource data body: []byte
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	// Close resource stream
	defer res.Body.Close()
	if string(body) != version {
		versionNew = string(body)
	}
}
