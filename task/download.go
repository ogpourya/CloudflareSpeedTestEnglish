package task

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/utils"

	"github.com/VividCortex/ewma"
)

const (
	bufferSize                     = 1024
	defaultURL                     = "https://cf.xiu2.xyz/url"
	defaultTimeout                 = 10 * time.Second
	defaultDisableDownload         = false
	defaultTestNum                 = 10
	defaultMinSpeed        float64 = 0.0
)

var (
	URL     = defaultURL
	Timeout = defaultTimeout
	Disable = defaultDisableDownload

	TestCount = defaultTestNum
	MinSpeed  = defaultMinSpeed
)

func checkDownloadDefault() {
	if URL == "" {
		URL = defaultURL
	}
	if Timeout <= 0 {
		Timeout = defaultTimeout
	}
	if TestCount <= 0 {
		TestCount = defaultTestNum
	}
	if MinSpeed <= 0.0 {
		MinSpeed = defaultMinSpeed
	}
}

func TestDownloadSpeed(ipSet utils.PingDelaySet) (speedSet utils.DownloadSpeedSet) {
	checkDownloadDefault()
	if Disable {
		return utils.DownloadSpeedSet(ipSet)
	}
	if len(ipSet) <= 0 { // Only proceed with download speed test if the IP array length (number of IPs) is greater than 0
		utils.Yellow.Println("[Info] Latency test result IP count is 0, skipping download speed test.")
		return
	}
	testNum := TestCount                        // Queue size waiting for download speed test, defaults to the download test count (-dn)
	if len(ipSet) < TestCount || MinSpeed > 0 { // If the IP array length after latency filtering is less than the download test count (-dn) (i.e. -dn expected count is insufficient), or a download speed minimum (-sl) is specified (which may require testing all IPs until enough qualifying ones are found or all are tested), correct the queue size to the number of IPs
		testNum = len(ipSet)
	}
	if testNum < TestCount { // If the queue size is less than the download test count (-dn) (clearly -dn expected count is insufficient), correct the download test count (-dn) to the queue size
		TestCount = testNum
	}

	utils.Cyan.Printf("Starting download speed test (minimum: %.2f MB/s, count: %d, queue: %d)\n", MinSpeed, TestCount, testNum)
	// Control the download speed test progress bar length to match the latency test progress bar length
	bar_a := len(strconv.Itoa(len(ipSet)))
	bar_b := "     "
	for i := 0; i < bar_a; i++ {
		bar_b += " "
	}
	bar := utils.NewBar(TestCount, bar_b, "")
	for i := 0; i < testNum; i++ {
		speed, colo := downloadHandler(ipSet[i].IP)
		ipSet[i].DownloadSpeed = speed
		if ipSet[i].Colo == "" { // Only write Colo if it is empty; otherwise it was already set during an httping speed test
			ipSet[i].Colo = colo
		}
		// After each IP download speed test, filter results by [download speed minimum]
		if speed >= MinSpeed*1024*1024 {
			bar.Grow(1, "")
			speedSet = append(speedSet, ipSet[i]) // Add to new array if above the download speed minimum
			if len(speedSet) == TestCount {        // Break out of loop once enough qualifying IPs are found (download test count -dn)
				break
			}
		}
	}
	bar.Done()
	if MinSpeed == 0.00 { // If no download speed minimum is specified, return all speed test data directly
		speedSet = utils.DownloadSpeedSet(ipSet)
	} else if utils.Debug && len(speedSet) == 0 { // If a download speed minimum is specified, debug mode is on, and no qualifying IPs were found, return all speed test data so the user can review current results and adjust expectations accordingly
		utils.Yellow.Println("[Debug] No IPs meet the download speed minimum condition. Ignoring condition and returning all speed test data (to help adjust conditions for next test).")
		speedSet = utils.DownloadSpeedSet(ipSet)
	}
	// Sort by speed
	sort.Sort(speedSet)
	return
}

func getDialContext(ip *net.IPAddr) func(ctx context.Context, network, address string) (net.Conn, error) {
	var fakeSourceAddr string
	if isIPv4(ip.String()) {
		fakeSourceAddr = fmt.Sprintf("%s:%d", ip.String(), TCPPort)
	} else {
		fakeSourceAddr = fmt.Sprintf("[%s]:%d", ip.String(), TCPPort)
	}
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, fakeSourceAddr)
	}
}

// Unified debug output for request errors
func printDownloadDebugInfo(ip *net.IPAddr, err error, statusCode int, url, lastRedirectURL string, response *http.Response) {
	finalURL := url // Default final URL so output works even when response is nil
	if lastRedirectURL != "" {
		finalURL = lastRedirectURL // If lastRedirectURL is not empty, a redirect occurred; prefer outputting the last redirect target
	} else if response != nil && response.Request != nil && response.Request.URL != nil {
		finalURL = response.Request.URL.String() // If response is not nil and Request and URL are not nil, get the last successful response address
	}
	if url != finalURL { // If URL and final address differ, a redirect occurred and the error originated from the redirected address
		if statusCode > 0 { // If status code is greater than 0, the error was caused by a subsequent HTTP status code
			utils.Red.Printf("[Debug] IP: %s, download speed test terminated, HTTP status code: %d, download URL: %s, redirected URL with error: %s\n", ip.String(), statusCode, url, finalURL)
		} else {
			utils.Red.Printf("[Debug] IP: %s, download speed test failed, error: %v, download URL: %s, redirected URL with error: %s\n", ip.String(), err, url, finalURL)
		}
	} else { // If URL and final address are the same, no redirect occurred
		if statusCode > 0 { // If status code is greater than 0, the error was caused by a subsequent HTTP status code
			utils.Red.Printf("[Debug] IP: %s, download speed test terminated, HTTP status code: %d, download URL: %s\n", ip.String(), statusCode, url)
		} else {
			utils.Red.Printf("[Debug] IP: %s, download speed test failed, error: %v, download URL: %s\n", ip.String(), err, url)
		}
	}
}

// return download Speed
func downloadHandler(ip *net.IPAddr) (float64, string) {
	var lastRedirectURL string // Records the last redirect target to output on access error
	client := &http.Client{
		Transport: &http.Transport{DialContext: getDialContext(ip)},
		Timeout:   Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			lastRedirectURL = req.URL.String() // Record each redirect target to output on access error
			if len(via) > 10 {                 // Limit to at most 10 redirects
				if utils.Debug { // Output more info in debug mode
					utils.Red.Printf("[Debug] IP: %s, too many redirects for download URL, terminating test, download URL: %s\n", ip.String(), req.URL.String())
				}
				return http.ErrUseLastResponse
			}
			if req.Header.Get("Referer") == defaultURL { // When using the default download URL, do not carry Referer on redirect
				req.Header.Del("Referer")
			}
			return nil
		},
	}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		if utils.Debug { // Output more info in debug mode
			utils.Red.Printf("[Debug] IP: %s, failed to create download speed test request, error: %v, download URL: %s\n", ip.String(), err, URL)
		}
		return 0.0, ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		if utils.Debug { // Output more info in debug mode
			printDownloadDebugInfo(ip, err, 0, URL, lastRedirectURL, response)
		}
		return 0.0, ""
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		if utils.Debug { // Output more info in debug mode
			printDownloadDebugInfo(ip, nil, response.StatusCode, URL, lastRedirectURL, response)
		}
		return 0.0, ""
	}

	// Get the region code from response headers
	colo := getHeaderColo(response.Header)

	timeStart := time.Now()           // Start time (now)
	timeEnd := timeStart.Add(Timeout) // End time = start time + download test duration

	contentLength := response.ContentLength // File size
	buffer := make([]byte, bufferSize)

	var (
		contentRead     int64 = 0
		timeSlice             = Timeout / 100
		timeCounter           = 1
		lastContentRead int64 = 0
	)

	var nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
	e := ewma.NewMovingAverage()

	// Loop to calculate; exit loop (stop test) if file download is complete (both values are equal)
	for contentLength != contentRead {
		currentTime := time.Now()
		if currentTime.After(nextTime) {
			timeCounter++
			nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
			e.Add(float64(contentRead - lastContentRead))
			lastContentRead = contentRead
		}
		// Exit loop (stop test) if download test time is exceeded
		if currentTime.After(timeEnd) {
			break
		}
		bufferRead, err := response.Body.Read(buffer)
		if err != nil {
			if err != io.EOF { // If an error occurs during download (e.g. Timeout) and it is not because the file finished downloading, exit loop (stop test)
				break
			} else if contentLength == -1 { // File download complete and file size unknown, exit loop (stop test). For example: https://speed.cloudflare.com/__down?bytes=200000000 — if download completes within 10 seconds, results will be significantly lower or even show 0.00 (when download speed is too fast)
				break
			}
			// Get the previous time slice
			last_time_slice := timeStart.Add(timeSlice * time.Duration(timeCounter-1))
			// Downloaded data / (current time - previous time slice / time slice)
			e.Add(float64(contentRead-lastContentRead) / (float64(currentTime.Sub(last_time_slice)) / float64(timeSlice)))
		}
		contentRead += int64(bufferRead)
	}
	return e.Value() / (Timeout.Seconds() / 120), colo
}
