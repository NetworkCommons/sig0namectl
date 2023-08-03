#!/bin/bash

set_vars() {
	# set_vars
	#
	# sets environment variables from env files
	#
	SCRIPT_NAME="${0#*/}"

	# Source env file for project-wide variable default values
	ENV_FILE=${ENV_FILE:-"env"}
	if [ -e ${ENV_FILE} ]
	then
	        . ./${ENV_FILE}
	        [[ -n ${DEBUG_SOURCED} ]] && echo "Sourced ${PWD}/${ENV_FILE} ..."
	fi

	# Source env file for script-wide default values
	if [ -e ${ENV_FILE}.${SCRIPT_NAME} ]
	then
	        . ${ENV_FILE}.${SCRIPT_NAME}
	        [[ -n ${DEBUG_SOURCED} ]] && echo "Sourced ${PWD}/${ENV_FILE}.${SCRIPT_NAME} ..."
	fi

	# Default existing ZONE fallback to $DOMAINNAME if set, else to domain set in $HOSTNAME, else error
	# ZONE MUST correspond to an existing DNS ZONE, that is resolve with an SOA record
	ZONE=${ZONE:-${DOMAINNAME}}
	ZONE=${ZONE:-${HOSTNAME#*.}}
	if [[ ! -n ${ZONE} ]]; then
	        echo "Error: DNS zone \$ZONE environment variable is not set & could not be determined from \$DOMAINNAME or \$HOSTNAME"
	        exit 1
	fi

	# path to search for nsupdate format key pairs
	NSUPDATE_SIG0_KEYPATH=${NSUPDATE_SIG0_KEYPATH:-"${PWD}"}

	DIG_QUERY_PARAM=${DIG_QUERY_PARAM:-}
	AVAHI_BROWSE_PARAM=${AVAHI_BROWSE_PARAM:-"-brat"}
	
	# Discover master (usually primary DNS server) of zone from master field of SOA record
	#
	ZONE_SOA_MASTER=$( get_soa_master ${ZONE} )
	if [[ ! -n ${ZONE_SOA_MASTER} ]]; then
	        echo "Warning: ZONE ${ZONE} SOA record does not resolve"
	        exit 1
	fi
}
