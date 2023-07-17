#!/bin/bash
#
get_soa() {
	# get_soa( fqdn )
	#
	# 	fqdn: FQDN to search and return superdomain that has SOA, eg host.example.com
	#
	#	returns string of most granular superdomain of fqdn parameter that contains a SOA resource record
	#
	local soa_fqdn="${1}"
	
	while [[ ! -n $(dig +short ${soa_fqdn} SOA) ]]
	do
		soa_fqdn=${soa_fqdn#*.}
		[[ ! "${soa_fqdn}" == *"."* ]] && soa_fqdn="" && break 
	done
	echo ${soa_fqdn}
}

get_soa_master() {
	        # get_soa_master( fqdn )
        #
        #       fqdn: FQDN portion of keypair filenames to search and return, eg host.exmple.com
        #
        #       returns master DNS server FQDN of most granular superdomain $fqdn with SOA record
	#
	local soa=""
	local master=""
	soa_fqdn=$( get_soa "${1}" )
	[[ ! "${soa_fqdn}" == "" ]] && master=$(dig ${DIG_QUERY_PARAM} +short ${soa_fqdn} SOA | cut -f1 -d' ')
	echo ${master}
}

if [[ -n ${TEST} ]]; then
	printf "** TEST get_soa()\n"
	DEBUG="true"
	for test in zenr.io test5.test4.test3.test2.test1.zembla.zenr.io test1.testzone.zenr.io test5.test4.test3.test2.test1.zembla.zenr.io.rubbish
	do
		test_ret=$(get_soa ${test})
		echo "'${test_ret}' is most granular domain with an SOA resource record for ${test}"
	done
fi

if [[ -n ${TEST} ]]; then
	printf "** TEST get_soa_master()\n"
	DEBUG="true"
	for test in zenr.io test5.test4.test3.test2.test1.zembla.zenr.io test1.testzone.zenr.io test5.test4.test3.test2.test1.zembla.zenr.io.rubbish
	do
		test_ret=$(get_soa_master ${test})
		echo "'${test_ret}' is master DNS server for DNS name ${test}"
	done
fi

