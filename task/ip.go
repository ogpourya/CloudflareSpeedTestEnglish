package task

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultInputFile = "ip.txt"

var (
	// TestAll test all ip
	TestAll = false
	// IPFile is the filename of IP Ranges
	IPFile = defaultInputFile
	IPText string
)

func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

func isIPv4(ip string) bool {
	return strings.Contains(ip, ".")
}

func randIPEndWith(num byte) byte {
	if num == 0 { // For a single IP like /32
		return byte(0)
	}
	return byte(rand.Intn(int(num)))
}

type IPRanges struct {
	ips     []*net.IPAddr
	mask    string
	firstIP net.IP
	ipNet   *net.IPNet
}

func newIPRanges() *IPRanges {
	return &IPRanges{
		ips: make([]*net.IPAddr, 0),
	}
}

// If it is a single IP, append a subnet mask; otherwise extract the subnet mask (r.mask)
func (r *IPRanges) fixIP(ip string) string {
	// If it does not contain '/', it is not a CIDR range but a single IP, so append /32 or /128 subnet mask
	if i := strings.IndexByte(ip, '/'); i < 0 {
		if isIPv4(ip) {
			r.mask = "/32"
		} else {
			r.mask = "/128"
		}
		ip += r.mask
	} else {
		r.mask = ip[i:]
	}
	return ip
}

// Parse the IP range to get the IP, IP range, and subnet mask
func (r *IPRanges) parseCIDR(ip string) {
	var err error
	if r.firstIP, r.ipNet, err = net.ParseCIDR(r.fixIP(ip)); err != nil {
		log.Fatalln("ParseCIDR err", err)
	}
}

func (r *IPRanges) appendIPv4(d byte) {
	r.appendIP(net.IPv4(r.firstIP[12], r.firstIP[13], r.firstIP[14], d))
}

func (r *IPRanges) appendIP(ip net.IP) {
	r.ips = append(r.ips, &net.IPAddr{IP: ip})
}

// Return the minimum value and available count of the fourth IP octet
func (r *IPRanges) getIPRange() (minIP, hosts byte) {
	minIP = r.firstIP[15] & r.ipNet.Mask[3] // Minimum value of the fourth IP octet

	// Get the number of hosts based on the subnet mask
	m := net.IPv4Mask(255, 255, 255, 255)
	for i, v := range r.ipNet.Mask {
		m[i] ^= v
	}
	total, _ := strconv.ParseInt(m.String(), 16, 32) // Total available IPs
	if total > 255 {                                 // Correct the available IP count for the fourth octet
		hosts = 255
		return
	}
	hosts = byte(total)
	return
}

func (r *IPRanges) chooseIPv4() {
	if r.mask == "/32" { // Single IP needs no randomization; add it directly
		r.appendIP(r.firstIP)
	} else {
		minIP, hosts := r.getIPRange()    // Get the minimum value and available count of the fourth IP octet
		for r.ipNet.Contains(r.firstIP) { // Keep looping as long as the IP has not exceeded the IP range
			if TestAll { // If testing all IPs
				for i := 0; i <= int(hosts); i++ { // Iterate through the last octet from min to max
					r.appendIPv4(byte(i) + minIP)
				}
			} else { // Randomize the last octet 0.0.0.X
				r.appendIPv4(minIP + randIPEndWith(hosts))
			}
			r.firstIP[14]++ // 0.0.(X+1).X
			if r.firstIP[14] == 0 {
				r.firstIP[13]++ // 0.(X+1).X.X
				if r.firstIP[13] == 0 {
					r.firstIP[12]++ // (X+1).X.X.X
				}
			}
		}
	}
}

func (r *IPRanges) chooseIPv6() {
	if r.mask == "/128" { // Single IP needs no randomization; add it directly
		r.appendIP(r.firstIP)
	} else {
		var tempIP uint8                  // Temporary variable to record the previous byte value
		for r.ipNet.Contains(r.firstIP) { // Keep looping as long as the IP has not exceeded the IP range
			r.firstIP[15] = randIPEndWith(255) // Randomize the last octet
			r.firstIP[14] = randIPEndWith(255) // Randomize the second-to-last octet

			targetIP := make([]byte, len(r.firstIP))
			copy(targetIP, r.firstIP)
			r.appendIP(targetIP) // Add IP address to the pool

			for i := 13; i >= 0; i-- { // Randomize from the third-to-last byte onward toward the front
				tempIP = r.firstIP[i]              // Save the previous byte value
				r.firstIP[i] += randIPEndWith(255) // Randomize 0~255 and add to the current byte
				if r.firstIP[i] >= tempIP {        // If the current byte is greater than or equal to the previous byte, randomization succeeded; exit the loop
					break
				}
			}
		}
	}
}

func loadIPRanges() []*net.IPAddr {
	ranges := newIPRanges()
	if IPText != "" { // Get IP range data from the parameter
		IPs := strings.Split(IPText, ",") // Split by comma into array and iterate
		for _, IP := range IPs {
			IP = strings.TrimSpace(IP) // Remove leading and trailing whitespace (spaces, tabs, newlines, etc.)
			if IP == "" {              // Skip empty entries (i.e. leading, trailing, or consecutive commas)
				continue
			}
			ranges.parseCIDR(IP) // Parse the IP range to get IP, range, and subnet mask
			if isIPv4(IP) {      // Generate all IPv4 / IPv6 addresses to test (single / random / all)
				ranges.chooseIPv4()
			} else {
				ranges.chooseIPv6()
			}
		}
	} else { // Get IP range data from file
		if IPFile == "" {
			IPFile = defaultInputFile
		}
		file, err := os.Open(IPFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() { // Iterate through each line of the file
			line := strings.TrimSpace(scanner.Text()) // Remove leading and trailing whitespace (spaces, tabs, newlines, etc.)
			if line == "" {                           // Skip empty lines
				continue
			}
			ranges.parseCIDR(line) // Parse the IP range to get IP, range, and subnet mask
			if isIPv4(line) {      // Generate all IPv4 / IPv6 addresses to test (single / random / all)
				ranges.chooseIPv4()
			} else {
				ranges.chooseIPv6()
			}
		}
	}
	return ranges.ips
}
