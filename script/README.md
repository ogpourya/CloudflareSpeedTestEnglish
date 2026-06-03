# XIU2/CloudflareSpeedTest - Scripts

These scripts call **CFST** and extend it with additional customized features.

****
> [!TIP]
> CFST was designed as a **command-line program** for **versatility**. It is not practical to put every possible requirement into the main program, especially personalized or niche needs. Doing so would increase maintenance effort and make the software bloated. One advantage of a command-line program is that it can be easily used together with other software and scripts.

The scripts below implement several requirements externally.

> In other words, a script calls CFST for speed testing, obtains the results, and then lets you decide how to process those results for your own needs, such as modifying Hosts.

Overall, these scripts are simple and focused. Besides meeting some users' needs, they serve more as example references for using CFST with scripts. Users who can write scripts or software can implement their own customized workflows.

If you have useful personal scripts, you can also send them through [**Issues**](https://github.com/XIU2/CloudflareSpeedTest/issues), [**Discussions**](https://github.com/XIU2/CloudflareSpeedTest/discussions), or **Pull requests** so they can be added here for more people to use.

> Tip: click the three-line icon in the upper-right corner to view the table of contents.

****
## 📑 cfst_hosts.sh / cfst_hosts.bat (included in the archive)

The script runs CFST to get the fastest IP and replaces the old CDN IP in the Hosts file.

> **Author:** [@XIU2](https://github.com/xiu2)  
> **Usage / feedback:** https://github.com/XIU2/CloudflareSpeedTest/discussions/312

<details>
<summary><code><strong>Changelog</strong></code></summary>

****

#### December 15, 2025, version v1.0.5 (cfst_hosts.bat)
 - **1. Fixed** an issue where the first IP row could not be obtained in newer CFST versions.

#### December 17, 2021, version v1.0.4
 - **1. Optimized** the keep-testing-until-a-qualifying-IP-is-found feature so it retests correctly when a minimum download speed is specified (commented out by default).

#### December 17, 2021, version v1.0.3
 - **1. Added** the keep-testing-until-a-qualifying-IP-is-found feature (commented out by default).
 - **2. Optimized** code.

#### September 29, 2021, version v1.0.2
 - **1. Fixed** an issue where the script did not exit when the speed test result contained 0 IPs.

#### April 29, 2021, version v1.0.1
 - **1. Optimized** behavior so `-p 0` is no longer required to avoid exiting on Enter.

#### January 28, 2021, version v1.0.0
 - **1. Released** the first version.

</details>

****

## 📑 cfst_3proxy.bat (included in the archive)

The script runs CFST to get the fastest IP and replaces the old Cloudflare CDN IP in the 3Proxy config file.  
It can redirect all Cloudflare CDN IPs to the fastest IP, accelerating all sites that use Cloudflare CDN without adding domains to Hosts one by one.

> **Author:** [@XIU2](https://github.com/xiu2)  
> **Usage / feedback:** https://github.com/XIU2/CloudflareSpeedTest/discussions/71

<details>
<summary><code><strong>Changelog</strong></code></summary>

****

#### December 15, 2025, version v1.0.6
 - **1. Fixed** an issue where the first IP row could not be obtained in newer CFST versions.

#### December 17, 2021, version v1.0.5
 - **1. Optimized** the keep-testing-until-a-qualifying-IP-is-found feature so it retests correctly when a minimum download speed is specified (commented out by default).

#### December 17, 2021, version v1.0.4
 - **1. Added** the keep-testing-until-a-qualifying-IP-is-found feature (commented out by default).
 - **2. Optimized** code.

#### September 29, 2021, version v1.0.3
 - **1. Fixed** an issue where the script did not exit when the speed test result contained 0 IPs.

#### April 29, 2021, version v1.0.2
 - **1. Optimized** behavior so `-p 0` is no longer required to avoid exiting on Enter.

#### March 16, 2021, version v1.0.1
 - **1. Optimized** code and comments.

#### March 13, 2021, version v1.0.0
 - **1. Released** the first version.

</details>

****

## 📑 cfst_dnspod.sh

If your domain is hosted on **DNSPod**, you can use the official DNSPod API to automatically update DNS records.  
The script runs CFST to get the fastest IP and updates DNS records to that fastest IP.

> **Author:** [@imashen](https://github.com/imashen)  
> **Usage / feedback:** https://github.com/XIU2/CloudflareSpeedTest/pull/533

<details>
<summary><code><strong>Changelog</strong></code></summary>

****

#### August 6, 2024, version v1.0.0
 - **1. Released** the first version.

</details>

****

## 📑 cfst_ddns.sh / cfst_ddns.bat

If your domain is hosted on **Cloudflare**, you can use the official Cloudflare API to automatically update DNS records.  
The script runs CFST to get the fastest IP and updates DNS records to that fastest IP through the Cloudflare API.

> **Author:** [@XIU2](https://github.com/xiu2)  
> **Usage / feedback:** https://github.com/XIU2/CloudflareSpeedTest/discussions/481

<details>
<summary><code><strong>Changelog</strong></code></summary>

****

#### December 15, 2025, version v1.0.6 (cfst_ddns.bat)
 - **1. Fixed** an issue where the first IP row could not be obtained in newer CFST versions.

#### October 6, 2024, version v1.0.5
 - **1. Added** API token support. Compared with globally scoped API keys, API tokens allow more flexible permission control.

#### December 17, 2021, version v1.0.4
 - **1. Added** the keep-testing-until-a-qualifying-IP-is-found feature (commented out by default).
 - **2. Optimized** code.

#### September 29, 2021, version v1.0.3
 - **1. Fixed** an issue where the script did not exit when the speed test result contained 0 IPs.

#### April 29, 2021, version v1.0.2
 - **1. Optimized** behavior so `-p 0` is no longer required to avoid exiting on Enter.

#### January 27, 2021, version v1.0.1
 - **1. Optimized** configuration loading from a file.

#### January 26, 2021, version v1.0.0
 - **1. Released** the first version.

</details>

****

## 📑 cfst_dnsmasq.sh

The script runs CFST to get the fastest IP and replaces the old Cloudflare CDN IP in the dnsmasq config file.

> **Author:** [@Sving1024](https://github.com/Sving1024)  
> **Usage / feedback:** https://github.com/XIU2/CloudflareSpeedTest/discussions/566

<details>
<summary><code><strong>Changelog</strong></code></summary>

****

#### January 22, 2025, version v1.0.1
 - **1. Fixed** an IPv6 issue.

#### December 28, 2024, version v1.0.0
 - **1. Released** the first version.

</details>

****

## Feature Suggestions / Issue Feedback

If you encounter problems while using these scripts, first check the corresponding **usage** discussion to see whether someone else has already asked about it.  
If you cannot find a similar issue, comment directly in the corresponding **usage** discussion and ask the author.
