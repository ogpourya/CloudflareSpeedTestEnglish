#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
# --------------------------------------------------------------
#	Project: CloudflareSpeedTest automatically update Hosts
#	Version: 1.0.4
#	Author: XIU2
#	Project: https://github.com/XIU2/CloudflareSpeedTest
# --------------------------------------------------------------

_CHECK() {
	while true
		do
		if [[ ! -e "nowip_hosts.txt" ]]; then
			echo -e "This script runs CFST to get the fastest IP and replace the Cloudflare CDN IP in Hosts.\nBefore use, please read: https://github.com/XIU2/CloudflareSpeedTest/issues/42#issuecomment-768273848"
			echo -e "On first use, change all Cloudflare CDN IPs in Hosts to the same IP."
			read -e -p "Enter that Cloudflare CDN IP and press Enter (this step will not be needed later): " NOWIP
			if [[ ! -z "${NOWIP}" ]]; then
				echo ${NOWIP} > nowip_hosts.txt
				break
			else
				echo "The IP cannot be empty!"
			fi
		else
			break
		fi
	done
}

_UPDATE() {
	echo -e "Starting speed test..."
	NOWIP=$(head -1 nowip_hosts.txt)

	# You can add or modify CFST runtime parameters here
	./cfst -o "result_hosts.txt"

	# If you need to keep testing until a qualifying IP is found, change the two exit 0 commands below to _UPDATE
	[[ ! -e "result_hosts.txt" ]] && echo "CFST speed test returned 0 IPs; skipping the following steps..." && exit 0

	# The following line is only needed for the keep-testing-until-a-qualifying-IP-is-found behavior
	# When a minimum download speed is specified but no IP meets all conditions, CFST outputs all IP results
	# Therefore, when specifying the -sl parameter, remove the leading # comment marker below to check the file line count (for example, if the download test count is 10, set the value below to 11)
	#[[ $(cat result_hosts.txt|wc -l) > 11 ]] && echo "CFST did not find an IP that fully meets the conditions; testing again..." && _UPDATE


	BESTIP=$(sed -n "2,1p" result_hosts.txt | awk -F, '{print $1}')
	if [[ -z "${BESTIP}" ]]; then
		echo "CFST speed test returned 0 IPs; skipping the following steps..."
		exit 0
	fi
	echo ${BESTIP} > nowip_hosts.txt
	echo -e "\nOld IP: ${NOWIP}\nNew IP: ${BESTIP}\n"

	echo "Backing up Hosts file (hosts_backup)..."
	\cp -f /etc/hosts /etc/hosts_backup

	echo -e "Starting replacement..."
	sed -i '' 's/'${NOWIP}'/'${BESTIP}'/g' /etc/hosts
	echo -e "Done..."
}

_CHECK
_UPDATE