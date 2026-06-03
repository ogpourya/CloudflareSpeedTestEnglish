# CloudflareSpeedTest

[![Go Version](https://img.shields.io/github/go-mod/go-version/ogpourya/CloudflareSpeedTestEnglish.svg?style=flat-square&label=Go&color=00ADD8&logo=go)](https://github.com/ogpourya/CloudflareSpeedTestEnglish/)
[![GitHub license](https://img.shields.io/github/license/ogpourya/CloudflareSpeedTestEnglish.svg?style=flat-square&label=License&color=00ADD8&logo=github)](https://github.com/ogpourya/CloudflareSpeedTestEnglish/)

Many websites use Cloudflare CDN, but the IPs allocated to visitors can be slow, have high latency, or drop packets. This tool helps you test Cloudflare CDN latency and download speed to find the fastest IP (IPv4 and IPv6) for your connection.

This is an English translation and maintenance fork of the original [XIU2/CloudflareSpeedTest](https://github.com/ogpourya/CloudflareSpeedTestEnglish).

---

## ⚡ Quick Start

### Installation

You can install the tool directly using Go:

```bash
go install github.com/ogpourya/CloudflareSpeedTestEnglish@latest
mv $(go env GOPATH)/bin/CloudflareSpeedTestEnglish $(go env GOPATH)/bin/cfst
```

Alternatively, you can clone the repository and build it manually:

```bash
git clone https://github.com/ogpourya/CloudflareSpeedTestEnglish.git
cd CloudflareSpeedTestEnglish
go build -o cfst main.go
```

### Basic Usage

Run the tool without parameters to start a default test:

```bash
./cfst
```

Or run with specific parameters:

```bash
# Filter IPs with latency < 200ms and test download speed for top 20
./cfst -tl 200 -dn 20
```

> **Note:** If latencies are extremely low (e.g., `0.xx ms`), ensure your proxy/VPN is turned off, as the tool might be measuring the speed to your local proxy instead of the Cloudflare edge.

---

## ⚙️ Parameters

Run `cfst -h` to see all available options:

| Parameter | Default | Description |
| :--- | :--- | :--- |
| `-n` | `200` | Latency test concurrency (Max 1000). |
| `-t` | `4` | Number of latency tests per IP. |
| `-dn` | `10` | Number of IPs to test for download speed. |
| `-dt` | `10` | Max duration (seconds) for each download test. |
| `-tp` | `443` | Port used for testing. |
| `-url` | `https://...` | Test URL for latency and download speed. |
| `-httping` | `false` | Use HTTP protocol for latency tests instead of TCPing. |
| `-tl` | `9999` | Latency upper limit (ms). |
| `-tll` | `0` | Latency lower limit (ms). |
| `-tlr` | `1.00` | Packet loss rate upper limit (0.00 to 1.00). |
| `-sl` | `0.00` | Download speed lower limit (MB/s). |
| `-p` | `10` | Number of results to display. |
| `-f` | `ip.txt` | IP range data file. |
| `-ip` | `""` | Specify IP ranges via CLI (comma-separated). |
| `-o` | `result.csv` | File path to save results. |
| `-dd` | `false` | Disable download speed testing. |
| `-allip` | `false` | Test every single IP (IPv4 only) instead of random samples. |

---

## 📊 Result Example

The tool provides a real-time progress bar and outputs the results in a table:

```text
IP Address        Sent   Recv   Loss%   Avg Latency   Download Speed (MB/s)  Colo
104.27.200.69     4      4      0.00    146.23        28.64                  LAX
172.67.60.78      4      4      0.00    139.82        15.02                  SEA
...
```

Results are also saved to `result.csv` by default.

---

## 🏗️ Compilation

To compile with version information:

```bash
go build -ldflags "-s -w -X main.version=v1.0.0" -o cfst main.go
```

### Cross-Compilation

**For Linux (amd64):**
```bash
GOOS=linux GOARCH=amd64 go build -o cfst_linux main.go
```

**For Windows (amd64):**
```bash
GOOS=windows GOARCH=amd64 go build -o cfst.exe main.go
```

---

## 📜 License

Distributed under the GPL-3.0 License. See `LICENSE` for more information.
