:: --------------------------------------------------------------
::	Project: CloudflareSpeedTest automatically update 3Proxy
::	Version: 1.0.6
::	Author: XIU2
::	Project: https://github.com/XIU2/CloudflareSpeedTest
:: --------------------------------------------------------------
@echo off
Setlocal Enabledelayedexpansion

:: Check whether administrator privileges have been obtained.
>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"

if '%errorlevel%' NEQ '0' (
    goto UACPrompt
) else ( goto gotAdmin )

:: Write a VBS script to relaunch this batch file as administrator.
:UACPrompt
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs"
    echo UAC.ShellExecute "%~s0", "", "", "runas", 1 >> "%temp%\getadmin.vbs"
    "%temp%\getadmin.vbs"
    exit /B

:: Delete the temporary VBS script if it exists.
:gotAdmin
    if exist "%temp%\getadmin.vbs" ( del "%temp%\getadmin.vbs" )
    pushd "%CD%"
    CD /D "%~dp0"

:: The script above checks for administrator privileges and requests them if missing.

:: If nowip_3proxy.txt does not exist, this is the first run.
if not exist "nowip_3proxy.txt" (
    echo This script runs CFST to get the fastest IP and replace the Cloudflare CDN IP in the 3Proxy config file.
    echo It can redirect all Cloudflare CDN IPs to the fastest IP, accelerating all sites that use Cloudflare CDN without adding domains to Hosts one by one.
    echo Before use, please read: https://github.com/XIU2/CloudflareSpeedTest/discussions/71
    echo.
    set /p nowip="Enter the current Cloudflare CDN IP used by 3Proxy and press Enter (this step will not be needed later):"
    echo !nowip!>nowip_3proxy.txt
    echo.
)

:: Read the current Cloudflare CDN IP from nowip_3proxy.txt.
set /p nowip=<nowip_3proxy.txt
echo Starting speed test...

:: RESET prepares the optional keep-testing-until-a-qualifying-IP-is-found behavior.
:: To use this feature, change the following 3 goto :STOP commands to goto :RESET.
:RESET

:: You can add or modify CFST runtime parameters here. echo.| auto-confirms program exit, so -p 0 is no longer needed.
echo.|cfst.exe -o "result_3proxy.txt"

:: Check whether the result file exists; if not, the result count is 0.
if not exist result_3proxy.txt (
    echo.
    echo CFST speed test returned 0 IPs; skipping the following steps...
    goto :STOP
)

:: Get the IP from the first result row.
for /f "skip=1 tokens=1 delims=," %%i in ('more result_3proxy.txt') do (
    SET bestip=%%i
    goto :END
)

:END

:: Check whether the obtained IP is empty or the same as the old IP.
if "%bestip%"=="" (
    echo.
    echo CFST speed test returned 0 IPs; skipping the following steps...
    goto :STOP
)
if "%bestip%"=="%nowip%" (
    echo.
    echo CFST speed test returned 0 IPs; skipping the following steps...
    goto :STOP
)

:: The following block is only needed for the keep-testing-until-a-qualifying-IP-is-found behavior.
:: When a minimum download speed is specified but no IP meets all conditions, CFST outputs all IP results.
:: Therefore, when specifying -sl, remove the leading :: comments below to check the file line count.
:: For example, if the download test count is 10, set the value below to 11.
::set /a v=0
::for /f %%a in ('type result_3proxy.txt') do set /a v+=1
::if %v% GTR 11 (
::    echo.
::    echo CFST did not find an IP that fully meets the conditions; testing again...
::    goto :RESET
::)

echo %bestip%>nowip_3proxy.txt
echo.
echo Old IP: %nowip%
echo New IP: %bestip%

:: Change the path below to your 3Proxy installation directory.
CD /d "D:\Program Files\3Proxy"
:: Make sure you have tested the 3Proxy service before running this script and that it is running.

echo.
echo Backing up 3proxy.cfg file (3proxy.cfg_backup)...
copy 3proxy.cfg 3proxy.cfg_backup
echo.
echo Starting replacement...
(
    for /f "tokens=*" %%i in (3proxy.cfg_backup) do (
        set s=%%i
        set s=!s:%nowip%=%bestip%!
        echo !s!
        )
)>3proxy.cfg

net stop 3proxy
net start 3proxy

echo Done...
echo.
:STOP
pause
