#!/bin/bash
#


NSUPDATE_SIG0_KEYPATH=${NSUPDATE_SIG0_KEYPATH:-"${PWD}"}

get_sig0_keyid() {
	# get_sig0_keyid( refvar, fqdn, search_path )
	#
	# 	refvar: global env var passed by reference will be set to *.[key|private] keypair filename prefix if found, or "" if not.
	#
	# 	fqdn: FQDN portion of keypair filenames to search and return, eg host.exmple.com
	#
	# 	search_path: path to search for *.key and *.private files
	#
	declare -n refvar="${1}"
	local fqdn="${2}"
	local search_path="${3:-${NSUPDATE_SIG0_KEYPATH}}"

	# get list of matching files
	#
	local FIND=$(find ${search_path:-${NSUPDATE_SIG0_KEYPATH}} -type f \( -iname K${fqdn}.\*.key -o -iname K${fqdn}.\*.private \))
	# split list into array
	local ar=(${FIND})
	local PREFIX=""
	[[ -n ${TEST} ]] && printf "==================== calling ${FUNCNAME[0]} ${1} ${2} ${3}\n"
	
	if (( ${#ar[@]} < 2 )); then
		refvar=""
		[[ -n ${DEBUG} ]] && printf "Warning: ${FUNCNAME[0]}(): path var set to \"${refvar}\": no keypair files for ${fqdn} found under ${search_path}.\n"
		return 1
	fi
	
	if (( ${#ar[@]} > 2 )); then
		refvar=""
		[[ -n ${DEBUG} ]] && printf "Warning: ${FUNCNAME[0]}(): path var set to \"${refvar}\": multiple keypair files for ${fqdn} found under ${search_path}.\n"
		return 1
	fi
	
	if [ "${ar[0]%.*}" != "${ar[1]%.*}" ]; then
		refvar=""
		[[ -n ${DEBUG} ]] && printf "Warning: ${FUNCNAME[0]}(): path var set to \"${refvar}\": no matching keypair files for ${fqdn} found under ${search_path}.\n"
		return 1
	fi

	PREFIX="${ar[0]%.*}"
	refvar="${PREFIX##*/}"
	[[ -n ${DEBUG} ]] && printf "Info: ${FUNCNAME[0]}(): set to \"${refvar}\": Unique keypair prefix found: ${PREFIX} in ${search_path}\n"

}


if [[ -n ${TEST} ]]; then
	printf "** TEST get_sig0_keyid()\n"
	DEBUG="true"
	get_sig0_keyid SIG0_KEYID zembla.zenr.io ${NSUPDATE_SIG0_KEYPATH}
	echo "SIG0_KEYID set to '${SIG0_KEYID}'"
	get_sig0_keyid SIG0_KEYID vortex.zenr.io
	echo "SIG0_KEYID set to '${SIG0_KEYID}'"
	get_sig0_keyid SIG0_KEYID no.such.name ${NSUPDATE_SIG0_KEYPATH}
	echo "SIG0_KEYID set to '${SIG0_KEYID}'"
	exit
fi
