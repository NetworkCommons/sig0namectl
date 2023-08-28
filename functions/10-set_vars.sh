#!/bin/bash

set_vars() {
	# set_vars
	#
	# sets environment variables from env files
	# set_vars $*
	# 	function also handles invocation paramters & requires parameters passed
	SCRIPT_NAME="${0#*/}"

	# Source env file for project-wide variable default values
	ENV_FILE=${ENV_FILE:-"env"}
	if [ -e ${ENV_FILE} ]; then
		. ./${ENV_FILE}
		[[ -n ${DEBUG_SET_VARS} ]] && echo "Sourced ${PWD}/${ENV_FILE} ..."
	fi

	# Source env file for script-wide default values
	if [ -e ${ENV_FILE}.${SCRIPT_NAME} ]; then
		source ${ENV_FILE}.${SCRIPT_NAME}
		[[ -n ${DEBUG_SET_VARS} ]] && echo "Sourced ${PWD}/${ENV_FILE}.${SCRIPT_NAME} ..."
	fi
	
	while getopts ":dhk:s:" ARG;
	do
		case "${ARG}" in
			d)
				NSUPDATE_ACTION="delete"
				[[ -n ${DEBUG_SET_VARS} ]] && echo "-d passed: delete action set"
				;;
			s)
				NSUPDATE_SIG0_KEYPATH="${OPTARG}"
				[[ -n ${DEBUG_SET_VARS} ]] && echo "-s passed: explicit sig0 keystore path parameter '${OPTARG}' given"
				;;
			k)
				NSUPDATE_AUTH_SIG0_KEY_FQDN="${OPTARG}"
				[[ -n ${DEBUG_SET_VARS} ]] && echo "-k passed: explicit sig0 key fqdn parameter '${OPTARG}' given"
				;;
			h)
				[[ -n ${DEBUG_SET_VARS} ]] && echo "-h passed: print help & exit '${OPTARG}' given"
				echo -e "USAGE: ${SCRIPT_NAME} [-d] [-s keystore_path] [-k keypair_fqdn] fqdn"
				echo -e "\t -d deletion action request (default action request is add)"
				echo -e "\t -s specify explicit keystore path"
				echo -e "\t -k specify explicit keypair FQDN for server action authentication"
				echo -e "\t fqdn specifies the fully qualified domain name (FQDN) to act upon (or under)"
				echo -e "\n All options overide default environment variable values, set on command line or in:"
				echo -e "\t${PWD}/${ENV_FILE}: for project-wide scope"
				echo -e "\t${PWD}/${ENV_FILE}.${SCRIPT_NAME}: for script-wide scope"
				exit 1
				;;
		esac
	done
	shift "$((OPTIND-1))"

	NEW_FQDN="${NEW_FQDN:-${1}}"
	if [[ ! -n ${NEW_FQDN} ]]; then
		if [[ -n ${ZONE} ]]; then
			if [[ -n ${NEW_SUBZONE} ]]; then
				NEW_FQDN="${NEW_SUBZONE}.${ZONE}"
			else
				NEW_FQDN="${ZONE}"
			fi
		else
			echo "No ZONE var set or FQDN argument given" && exit 1
		fi
	fi
	shift 1
	CMDLINE_EXTRA_PARAMS=$*

	# Last attempt to set default sig0 keystore path
	NSUPDATE_SIG0_KEYPATH=${NSUPDATE_SIG0_KEYPATH:-"${PWD}"}
	ZONE="${ZONE:-$(get_soa ${NEW_FQDN})}" # sanity check?
	NEW_SUBZONE=${NEW_SUBZONE:-${NEW_FQDN%*.${ZONE}}}
	[[ -n ${DEBUG_SET_VARS} ]] && echo "NEW_FQDN='${NEW_FQDN}'" && echo "ZONE='${ZONE}'" && echo "NEW_SUBZONE='${NEW_SUBZONE}'" && echo "NSUPDATE_SIG0_KEYPATH='${NSUPDATE_SIG0_KEYPATH}'"
	if [[ ${#ZONE} < 2 ]]; then
		echo "Error: DNS zone ZONE='${ZONE}' environment variable is not set & could not be determined from \$DOMAINNAME or \$HOSTNAME"
		echo "DEBUG: NEW_FQDN='${NEW_FQDN}'"
		exit 1
	fi

	# Discover master (usually primary DNS server) of zone from master field of SOA record
	#
	ZONE_SOA_MASTER=$( get_soa_master ${ZONE} )
	if [[ ! -n ${ZONE_SOA_MASTER} ]]; then
		echo "Error: Could not resolve SOA record for ZONE '${ZONE}'"
		exit 1
	fi

	# External tool defaults
	DIG_QUERY_PARAM=${DIG_QUERY_PARAM:-}
	AVAHI_BROWSE_PARAM=${AVAHI_BROWSE_PARAM:-"-brat"}
}
