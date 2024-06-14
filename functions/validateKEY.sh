#!/bin/bash
function validateKEY()
{

	# validate KEY RR presence at key_fqdn
	#
	#       ip[6,4]: 1 or more IPv6 or IPv4 addresses
	#
	#       return value 0 = valid address
	#	non zero value = invalid address
	local key_fqdn=$1
	local query_answer=$(dig +short ${key_fqdn} KEY)
	local stat=$?
	# echo "validateKEY() DEBUG: resolve ${key_fqdn} KEY query_answer='${query_answer}'"
	# echo "length of query_answer=${#query_answer}"
	[[ ${#query_answer} > 0 ]] && stat=0 || stat=1
	return $stat
}

#function validateIPv6()
#{
#	local ip=$1
#	ipcalc -s -6 -c ${ip}
#	local stat=$?
#	#TODO: doesn't handle arbitrary zero compression yet
#	#local stat=1
#	#if [[ $ip =~ ^[0-9,a-f,A-F]{1,4}\:[0-9,a-f,A-F]{1,4}\:[0-9,a-f,A-F]{1,4}\:[0-9,a-f,A-F]{1,4}\:[0-9,a-f,A-F]{1,4}\:\:[0-9,a-f,A-F]{1,4}$ ]]; then
#	#	# OIFS=$IFS
#	#	# IFS=':'
#	#	# ip=($ip)
#	#	# IFS=$OIFS
#	#	# [[ ${ip[0]} -le 255 && ${ip[1]} -le 255 \
#	#	# && ${ip[2]} -le 255 && ${ip[3]} -le 255 ]]
#	#	stat=0
#	#fi
#	#echo "$stat"
#	return $stat
#}
#
#if [[ -n ${TEST} ]]; then
#	test_list=(
#		127.0.0.1
#		1.2.3.4
#		10.20.30.40
#		100.200.100.100
#		123.256.123.12
#		3.3.3.
#		3.3.3
#		5.5.5.5.5
#		fe80::eeb4:b20e:7677:6144
#		fe80::eeb4:b20e:7677:6144:8888:9999:AAAAA
#		fe80:eeb4:b20e:
#		2a01:4f8:c17:3dd5:8000::10
#	)
#	echo "-- TEST -- ipv4"
#	for ip in "${test_list[@]}"; do
#		res=$(validateIPv4 ${ip})
#		echo "${ip}= '${res}' : '$?' "
#	done
#	echo "-- TEST -- ipv6"
#	for ip in "${test_list[@]}"; do
#		res="$(validateIPv6 ${ip})"
#		echo "${ip}= '${res}' : '$?' "
#	done
#fi
