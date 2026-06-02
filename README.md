# XIU2/CloudflareSpeedTest

[![Go Version](https://img.shields.io/github/go-mod/go-version/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Go&color=00ADD8&logo=go)](https://github.com/XIU2/CloudflareSpeedTest/)
[![Release Version](https://img.shields.io/github/v/release/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Release&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/releases/latest)
[![GitHub license](https://img.shields.io/github/license/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=License&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)
[![GitHub Star](https://img.shields.io/github/stars/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Star&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)
[![GitHub Fork](https://img.shields.io/github/forks/XIU2/CloudflareSpeedTest.svg?style=flat-square&label=Fork&color=00ADD8&logo=github)](https://github.com/XIU2/CloudflareSpeedTest/)

Many foreign websites use Cloudflare CDN, but the IPs allocated to visitors in mainland China can be slow, have high latency, or drop packets.

Although Cloudflare publishes all its [IP ranges](https://www.cloudflare.com/ips/), finding the fastest IP for yourself manually from so many candidates would be extremely tedious. This tool helps you do exactly that!

**Test Cloudflare CDN latency and download speed to find the fastest IP (IPv4 and IPv6) for you!** If you like it, please **give us a `⭐`**!

Other projects from the author:
- [TrackersList.com](https://github.com/XIU2/TrackersListCollection) - High-quality BT Tracker list to speed up BT downloads!
- [UserScript](https://github.com/XIU2/UserScript) - Tampermonkey scripts for GitHub fast download, eye care mode, etc.
- [SNIProxy](https://github.com/XIU2/SNIProxy) - A simple, platform-independent SNI proxy.

> [!IMPORTANT]
> Cloudflare CDN terms explicitly forbid using proxy servers through their CDN. Doing so is at your own risk. Do not rely on it completely. See [#382](https://github.com/XIU2/CloudflareSpeedTest/discussions/382) and [#383](https://github.com/XIU2/CloudflareSpeedTest/discussions/383).

---

## ⚡ Quick Start

### Download & Run

1. Download the compiled file from [GitHub Releases](https://github.com/XIU2/CloudflareSpeedTest/releases) and extract it.
2. Double-click `cfst.exe` (on Windows) and wait for the speed test to finish.

<details>
<summary><strong>Windows (Scoop installation)</strong></summary>

If you use Scoop, you can install it like this:
```sh
# Add the dorado bucket
scoop bucket add dorado https://github.com/chawyehsu/dorado
# Install cloudflare-speedtest
scoop install dorado/cloudflare-speedtest
```
</details>

<details>
<summary><strong>Linux / macOS Usage Example</strong></summary>

The following commands are examples. Check [Releases](https://github.com/XIU2/CloudflareSpeedTest/releases) for the latest version.

On macOS, you can download/extract using the finder or browser, but you still need to run it in the terminal (don't forget to grant execution permission).

```bash
# Create a new folder
mkdir cfst
cd cfst

# Download the archive (always points to the latest release)
wget -N https://github.com/XIU2/CloudflareSpeedTest/releases/latest/download/cfst_linux_amd64.tar.gz

# Extract
tar -zxf cfst_linux_amd64.tar.gz

# Grant execution permission
chmod +x cfst

# Run (no parameters)
./cfst

# Run with parameters example
./cfst -tl 200 -dn 20
```

> **Note:** If the average latency is extremely low (like `0.xx ms`), it means CFST is running through a proxy/VPN. Please turn off your proxy client first.
> If running on a router, turn off its proxy first, otherwise the results will be inaccurate.
</details>

> Running on mobile devices: **[Android Guide](https://github.com/XIU2/CloudflareSpeedTest/discussions/61), [Android APP (Hsia97)](https://github.com/Hsia97/CFSTAPP), [Android APP (xianshenglu)](https://github.com/xianshenglu/cloudflare-ip-tester-app), [iOS Guide](https://github.com/XIU2/CloudflareSpeedTest/discussions/321)**

> [!NOTE]
> This tool is for websites. **It does NOT support finding IPs for Cloudflare WARP** (which uses UDP). See [#392](https://github.com/XIU2/CloudflareSpeedTest/discussions/392).

---

### Result Example

By default, the program shows the **top 10 fastest IPs**:

```text
IP Address        Sent   Recv   Loss%   Avg Latency   Download Speed (MB/s)  Colo
104.27.200.69     4      4      0.00    146.23        28.64                  LAX
172.67.60.78      4      4      0.00    139.82        15.02                  SEA
104.25.140.153    4      4      0.00    146.49        14.90                  SJC
104.27.192.65     4      4      0.00    140.28        14.07                  LAX
172.67.62.214     4      4      0.00    139.29        12.71                  LAX
...
```

- The first line is the fastest IP overall (lowest latency and fastest download speed).
- The full results are saved in `result.csv` in the same directory. You can open it with a text editor or spreadsheet software (like Excel).

---

## ⚙️ Advanced Usage

Run `cfst -h` to see all available options:

```text
CloudflareSpeedTest vX.X.X
Test the latency and speed of CDN / website IPs to find the fastest one!
https://github.com/XIU2/CloudflareSpeedTest

Parameters:
    -n 200
        Latency test concurrency. More threads mean faster testing, but do not set too high on weak devices like routers (Default 200, Max 1000).
    -t 4
        Number of latency tests per IP (Default 4).
    -dn 10
        Number of IPs to test for download speed, starting from the lowest latency IP (Default 10).
    -dt 10
        Max duration (seconds) for each download speed test. Do not set too short (Default 10).
    -tp 443
        Port used for latency and download speed tests (Default 443).
    -url https://cf.xiu2.xyz/url
        The test URL used for latency (HTTPing) and download speed tests. It is recommended to host your own.
        During download speed tests, the tool extracts the node's airport code (Colo) from HTTP headers (Supports Cloudflare, AWS CloudFront, Fastly, Gcore, CDN77, Bunny, etc.).
    -httping
        Use HTTP protocol for latency tests instead of TCPing (Uses URL defined in -url).
        HTTPing behaves like a network scan. Lower the concurrency (-n) on servers to avoid temporary suspension.
    -httping-code 200
        Expected HTTP status code for HTTPing to be considered successful (Default 200, 301, 302).
    -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
        Filter by region/airport code. Comma separated. Case insensitive. Only works in HTTPing mode (Default: all regions).
        Supports:
          - Cloudflare, AWS CloudFront, Fastly: IATA 3-letter codes (e.g., HKG, LAX)
          - CDN77, Bunny: 2-letter country codes (e.g., US, CN)
          - Gcore: 2-letter city codes (e.g., FR, AM)
    -tl 200
        Latency upper limit (ms). Only output IPs with average latency below this value (Default 9999).
    -tll 40
        Latency lower limit (ms). Only output IPs with average latency above this value (Default 0).
    -tlr 0.2
        Packet loss rate upper limit. Range 0.00 to 1.00. Set to 0 to filter out any IPs with packet loss (Default 1.00).
    -sl 5
        Download speed lower limit (MB/s). Only output IPs faster than this. Stops when -dn target is reached (Default 0.00).
    -p 10
        Number of results to display on the command line. Set to 0 to exit without printing (Default 10).
    -f ip.txt
        The IP range data file. Supports other CDN IP ranges (Default ip.txt).
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        Directly specify IP ranges via command line, comma-separated (Default: empty).
    -o result.csv
        The file path to save results. Set to empty -o "" (or -o " ") to disable saving (Default result.csv).
    -dd
        Disable download speed testing. Results will only be sorted by latency (Default: enabled).
    -allip
        Test every single IP in the ranges (IPv4 only). By default, only one random IP is tested per /24 range.
    -debug
        Enable debug logging. Shows why tests fail or get interrupted.
    -v
        Print version and check for updates.
    -h
        Print help instructions.
```

---

## 📊 Understanding the Interface

Sometimes users get confused by the interface progress bars. Here is how it works:

<details>
<summary><strong>Click to expand interface explanation</strong></summary>

Let's say you run `./cfst -tll 40 -tl 150 -sl 1 -dn 5`. The output looks like:

```text
Start latency test (Mode: TCP, Port: 443, Range: 40 ~ 150 ms, Loss: 1.00)
321 / 321 [-----------------------------------------------------------] Available: 30
Start download test (Min Speed: 1.00 MB/s, Count: 5, Queue: 10)
3 / 5 [-----------------------------------------↗--------------------]
...
```

### Explanations:

1. **Why did `Available: 30` change to `Queue: 10`?**
   - During latency testing, `Available: 30` means 30 IPs successfully responded without timeout.
   - However, since you filtered with `-tll 40` and `-tl 150` (latency between 40-150ms), only 10 IPs met your criteria. This leaves 10 IPs in the `Queue` waiting for download speed tests.

2. **Why does it say `3 / 5` instead of `5`?**
   - You asked for 5 IPs that meet the speed requirement of `1.00 MB/s` (`-dn 5 -sl 1`).
   - The tool tested all 10 IPs in the queue, but only 3 of them were faster than `1.00 MB/s`. The other 7 were too slow. So the test ended with 3/5.

3. **What if the progress bar seems stuck?**
   - The progress bar for download speed only advances when an IP matches your `-sl` speed limit. If your `-sl` is set too high (e.g., `-sl 10` on a slow connection), the tool will keep testing IP after IP trying to find one that reaches 10 MB/s. Lower your `-sl` value if this happens.
</details>

---

## 💡 Practical Examples

Here are common ways to run CloudflareSpeedTest.

<details>
<summary><strong>Windows: Running with parameters</strong></summary>

### Method 1: Using CMD / PowerShell
1. Open the folder containing `cfst.exe`.
2. Type `cmd` in the address bar of Windows Explorer and press Enter. This opens CMD directly in that folder.
3. Type the command with parameters, for example:
   ```cmd
   cfst -tl 200 -dn 20
   ```
   *(For PowerShell, use `.\cfst -tl 200 -dn 20`)*

### Method 2: Create a Windows Shortcut
1. Right-click `cfst.exe` -> **Create shortcut**.
2. Right-click the shortcut -> **Properties**.
3. In the **Target** field, add your parameters outside the quotes, e.g.:
   `"D:\cfst\cfst.exe" -tl 200 -dn 20 -o " "`
   *(Note: Set -o " " with a space to disable result file saving without errors)*
</details>

<details>
<summary><strong>Testing IPv4 vs IPv6</strong></summary>

```bash
# Test IPv4 using the default ip.txt
cfst -f ip.txt

# Test IPv6 using ipv6.txt
cfst -f ipv6.txt

# Test specific IPs directly
cfst -ip 1.1.1.1,2606:4700::/32
```
</details>

<details>
<summary><strong>HTTPing Mode (HTTP Latency Test)</strong></summary>

HTTPing checks if you can establish an actual HTTP/HTTPS connection to a URL through the IP. Latencies usually follow: **ICMP < TCP < HTTP**.

```bash
# Switch to HTTPing mode
cfst -httping

# Specify a custom URL for testing
cfst -httping -url https://cf.xiu2.xyz/url

# Expecting HTTP status 200
cfst -httping -httping-code 200
```
</details>

<details>
<summary><strong>Filter by Location (Colo)</strong></summary>

Cloudflare uses Anycast IPs, which routing changes dynamically. You can filter IPs by their physical data centers using HTTPing mode:

```bash
# Only find IPs routed through Hong Kong, Tokyo, Los Angeles, Seattle, etc.
cfst -httping -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
```
See [cloudflarestatus.com](https://www.cloudflarestatus.com/) for a list of airport codes.
</details>

<details>
<summary><strong>Setting Speed and Latency Targets</strong></summary>

```bash
# Find IPs with latency below 200ms
cfst -tl 200

# Find IPs with latency between 60ms and 200ms
cfst -tll 60 -tl 200

# Filter out IPs that have any packet loss
cfst -tlr 0

# Stop when you find 10 IPs with download speed > 5 MB/s
cfst -sl 5 -dn 10
```
</details>

<details>
<summary><strong>Testing Single/Custom IPs</strong></summary>

Write the IPs or subnets into a text file, e.g., `my_ips.txt`:
```text
1.1.1.1
1.0.0.1/24
2606:4700::/32
```
Then run:
```bash
cfst -f my_ips.txt
```
</details>

---

## 🛠️ Speed Up All Cloudflare-backed Websites Automatically

Adding IPs to your `hosts` file one by one is tedious. You can use these methods to speed up all Cloudflare CDN-backed sites at once:
- **[3Proxy redirection guide](https://github.com/XIU2/CloudflareSpeedTest/discussions/71)**: Redirect all Cloudflare traffic to your fastest selected IP.
- **[Local DNS Server modification guide](https://github.com/XIU2/CloudflareSpeedTest/discussions/317)**: Use a local DNS server to map domains to your selected IP.
- **[Auto-Update Hosts script](https://github.com/XIU2/CloudflareSpeedTest/discussions/312)**: Scripts for Windows/Linux to automate updating hosts with the fastest IP.

---

## 🚀 Derivative Projects

Here are some mobile apps, router plugins, and rewrites created by the community:
- [CFSTAPP Android App (Hsia97)](https://github.com/Hsia97/CFSTAPP)
- [Cloudflare IP Tester App Android (xianshenglu)](https://github.com/xianshenglu/cloudflare-ip-tester-app)
- [OpenWrt App (mingxiaoyu)](https://github.com/mingxiaoyu/luci-app-cloudflarespeedtest)
- [OpenWrt Native Build (immortalwrt-collections)](https://github.com/immortalwrt-collections/openwrt-cdnspeedtest)
- [CloudflareST-Rust Rewrite](https://github.com/GuangYu-yu/CloudflareST-Rust)
- [CloudflareST Go Fork](https://github.com/masgzy/CloudflareST)

---

## 🏗️ Manual Compilation

Version numbers are injected during compilation. To compile the program manually, use `go build` with `-ldflags`:

```bash
go build -ldflags "-s -w -X main.version=v2.4.0"
```

To cross-compile for other systems/architectures from Windows:

**Compile for Linux (amd64):**
```cmd
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w -X main.version=v2.4.0"
```

**Compile for Windows (32-bit) from Linux:**
```bash
GOOS=windows
GOARCH=386
go build -ldflags "-s -w -X main.version=v2.4.0"
```

Use `go tool dist list` to see all supported platforms.

---

## 💖 Support the Project

If this project helped you, consider buying me a coffee!

![Donation QR Codes](https://github.com/XIU2/XIU2/blob/master/img/zs-01.png)![Donation QR Codes](https://github.com/XIU2/XIU2/blob/master/img/zs-02.png)

---

## 📜 License

The GPL-3.0 License. See [LICENSE](LICENSE) for details.
