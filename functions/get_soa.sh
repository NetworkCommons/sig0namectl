#!/bin/bash
#

get_soa() {
	# get_soa( fqdn )
	#
	# 	fqdn: FQDN portion of keypair filenames to search and return, eg host.exmple.com
	#
	#	returns string of most granular subdomain of fqdn parameter that contains a SOA resource record
	#
	local soa_fqdn="${1}"
	
	while [[ ! -n $(dig +short ${soa_fqdn}) ]]
	do
		soa_fqdn=${soa_fqdn#*.}
		[[ ! "${soa_fqdn}" == *"."* ]] && soa_fqdn="" && break 
	done
	echo ${soa_fqdn}
}

if [[ -n ${TEST} ]]; then
	printf "** TEST get_soa_fqdn()\n"
	DEBUG="true"
	for test in zenr.io test5.test4.test3.test2.test1.zembla.zenr.io test1.testzone.zenr.io test5.test4.test3.test2.test1.zembla.zenr.io.rubbish
	do
		test_ret=$(get_soa ${test})
		echo "'${test_ret}' is most granular domain with an SOA resource record for ${test}"
	done
fi
