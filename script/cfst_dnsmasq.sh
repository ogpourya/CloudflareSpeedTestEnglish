#!/usr/bin/env bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH
# --------------------------------------------------------------
#	Project: CloudflareSpeedTest automatically update dnsmasq config file
#	Version: 1.0.1
#	Author: XIU2,Sving1024
#	Project: https://github.com/XIU2/CloudflareSpeedTest
# --------------------------------------------------------------

_UPDATE() {
	echo -e "Starting speed test..."
	BESTIP=""
	BESTIP_IPV6="::"
	# You can add or modify CFST runtime parameters here
	./cfst -o "result_hosts.txt"
	# Uncomment to test IPv6
	#./cfst -o "result_hosts_ipv6.txt" -f ipv6.txt

	# If you need to keep testing until a qualifying IP is found, change the two exit 0 commands below to _UPDATE
	[[ ! -e "result_hosts.txt" ]] && echo "CFST speed test returned 0 IPs; skipping the following steps..." && exit 0

	# The following line is only needed for the keep-testing-until-a-qualifying-IP-is-found behavior
	# When a minimum download speed is specified but no IP meets all conditions, CFST outputs all IP results
	# Therefore, when specifying the -sl parameter, remove the leading # comment marker below to check the file line count (for example, if the download test count is 10, set the value below to 11)
	#[[ $(cat result_hosts.txt|wc -l) > 11 ]] && echo "CFST did not find an IP that fully meets the conditions; testing again..." && _UPDATE

	BESTIP=$(sed -n "2,1p" result_hosts.txt | awk -F, '{print $1}')
	# Uncomment to test IPv6
	#BESTIP_IPV6=$(sed -n "2,1p" result_hosts_ipv6.txt | awk -F, '{print $1}')

	if [[ -z "${BESTIP}" ]]; then
		echo "CFST speed test returned 0 IPs; skipping the following steps..."
		exit 0
	fi
	echo ${BESTIP} > nowip_hosts.txt
	echo -e "Best IPv4 IP: ${BESTIP}\n"
	# Uncomment to test IPv6
	#echo -e "Best IPv6 IP: ${BESTIP_IPV6}\n"

    [[ -f cloudflare.conf ]] && rm cloudflare.conf

    cat site.conf | while read domain
    do
        if [[ ${domain:0:1} != "#" && ${domain} != "" ]]; then 
			echo "address=/${domain}/${BESTIP}" >> "cloudflare.conf"
			echo "address=/${domain}/${BESTIP_IPV6}" >> "cloudflare.conf"
		fi
    done

    [[ -f /etc/dnsmasq.d/cloudflare.conf ]] && rm /etc/dnsmasq.d/cloudflare.conf
    cp cloudflare.conf /etc/dnsmasq.d/cloudflare.conf
    systemctl restart dnsmasq.service
}

_UPDATE