:: --------------------------------------------------------------
::	Project: CloudflareSpeedTest automatically update DNS records
::	Version: 1.0.6
::	Author: XIU2
::	Project: https://github.com/XIU2/CloudflareSpeedTest
:: --------------------------------------------------------------
@echo off
Setlocal Enabledelayedexpansion

:: You can add or modify CFST runtime parameters here. echo.| auto-confirms program exit, so -p 0 is no longer needed.
echo.|cfst.exe -o "result_ddns.txt"

:: Check whether the result file exists; if not, the result count is 0
if not exist result_ddns.txt (
    echo.
    echo CFST speed test returned 0 IPs; skipping the following steps...
    goto :END
)

for /f "skip=1 tokens=1 delims=," %%i in (result_ddns.txt) do (
    Echo %%i
    if "%%i"=="" (
        echo.
        echo CFST speed test returned 0 IPs; skipping the following steps...
        goto :END
    )
::  API key mode (global permissions)
    curl -X PUT "https://api.cloudflare.com/client/v4/zones/zone ID/dns_records/DNS record ID" ^
            -H "X-Auth-Email: account email" ^
            -H "X-Auth-Key: previously obtained API key" ^
            -H "Content-Type: application/json" ^
            --data "{\"type\":\"A\",\"name\":\"full domain name\",\"content\":\"%%i\",\"ttl\":1,\"proxied\":true}"
::  API token mode (custom permissions). To use this mode, delete or comment out the lines above, then remove the leading "::" comment markers below.
::    curl -X PUT "https://api.cloudflare.com/client/v4/zones/zone ID/dns_records/DNS record ID" ^
::            -H "Authorization: Bearer previously obtained API token" ^
::            -H "Content-Type: application/json" ^
::            --data "{\"type\":\"A\",\"name\":\"full domain name\",\"content\":\"%%i\",\"ttl\":1,\"proxied\":true}"

        goto :END
)
:END
pause