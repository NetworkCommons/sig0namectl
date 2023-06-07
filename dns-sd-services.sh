#!/bin/bash
#
# Sends DNS-DS updates for domain registration and browsing records
#------------------------------------------------------------------------------
SCRIPT_NAME="${0#*/}"
# Source env file for project-wide variable default values
ENV_FILE=${ENV_FILE:-".env"}
if [ -e ${ENV_FILE} ]
then
        . ${ENV_FILE}
        [[ -n ${DEBUG} ]] && echo "Sourced ${PWD}/${ENV_FILE} ..."
fi

# Source env file for script-wide default values
if [ -e ${ENV_FILE}.${SCRIPT_NAME} ]
then
	. ${ENV_FILE}.${SCRIPT_NAME}
	[[ -n ${DEBUG} ]] && echo "Sourced ${PWD}/${ENV_FILE}.${SCRIPT_NAME} ..."
fi

# Default existing ZONE fallback to $DOMAINNAME if set, else to domain set in $HOSTNAME, else error
# ZONE MUST correspond to an existing DNS ZONE, that is resolve with an SOA record
ZONE=${ZONE:-${DOMAINNAME}}
ZONE=${ZONE:-${HOSTNAME#*.}}
if [[ ! -n ${ZONE} ]]; then
        echo "Error: DNS zone \$ZONE environment variable is not set & could not be determined from \$DOMAINNAME or \$HOSTNAME"
        exit 1
fi

# Discover master (usually primary DNS server) of zone from SOA record
#
DIG_QUERY_PARAM=${DIG_QUERY_PARAM:-}
ZONE_SOA_MASTER=${ZONE_SOA_MASTER:-$(dig ${DIG_QUERY_PARAM} +short ${ZONE} SOA | cut -f1 -d' ')}
if [[ ! -n ${ZONE_SOA_MASTER} ]]; then
        echo "Warning: ZONE ${ZONE} SOA record does not resolve"
fi

#------------------------------------------------------------------------------
# test


NEW_SUBZONE=${NEW_SUBZONE:-""}

# select zone or subdomain within zone
if [[ -n ${NEW_SUBZONE} ]]; then
        DNSSD_DOMAIN="${NEW_SUBZONE}.${ZONE}"
else
	DNSSD_DOMAIN="${ZONE}"
fi

DNSSD_SERVICES=${DNSSD_SERVICES:-""}
if [[ ! -n ${DNSSD_SERVICES} ]]; then
        echo "Error: DNSSD SERVICES services \$DNSSD_SERVICES environment variable is not set. No services have been defined to browse in domain ${DNSSD_SERVICES}"
        exit 1
fi

# Form NSUPDATE input
#

# form SIG0 private auth key param
NSUPDATE_SIG0_KEY="/home/vortex/src/great-dane/test_go/Kzembla.zenr.io.+015+23799"
if [[ -n ${NSUPDATE_SIG0_KEY} ]]; then
	NSUPDATE_PARAM="${NSUPDATE_PARAM} -k ${NSUPDATE_SIG0_KEY}.key"
else
	NSUPDATE_PARAM=""
fi

# explicitly set server to send update to
NSUPDATE_SET_SERVER="server ${ZONE_SOA_MASTER}"

# generate nsupdate precondition to update (to avoid double entries)
# TODO test whether multiple prereqs work ...
# for now, just test for existance of lb._dns-sd._udp. in DNS_DOMAIN
NSUPDATE_ACTION=${NSUPDATE_ACTION:-"add"}
if [ "${NSUPDATE_ACTION}" == "add" ]; then
	NSUPDATE_PRECONDITION_SET="nxrrset"
else
	NSUPDATE_PRECONDITION_SET="yxrrset"
fi
NSUPDATE_PRECONDITION="prereq ${NSUPDATE_PRECONDITION_SET} lb._dns-sd._udp.${DNSSD_DOMAIN}. IN PTR"

# set RR TTLs
NSUPDATE_TTL="600"

# define DNS RR updates
NSUPDATE_INPUT="${NSUPDATE_SET_SERVER}\n"
for SERVICE in ${DNSSD_SERVICES}
do
	NSUPDATE_INPUT="${NSUPDATE_INPUT}update ${NSUPDATE_ACTION} _services._dns-sd._udp.${DNSSD_DOMAIN} ${NSUPDATE_TTL} PTR ${SERVICE}.${DNSSD_DOMAIN}.\n"
done
NSUPDATE_INPUT="${NSUPDATE_INPUT}send\n"
NSUPDATE_INPUT="${NSUPDATE_INPUT}quit\n"

if [[ -n ${DEBUG} ]]; then
	echo "---DEBUG"
	echo "SCRIPT_NAME		=	${SCRIPT_NAME}"
	echo "ZONE 			=	${ZONE}"
	echo "DNSSD_DOMAIN 		=	${DNSSD_DOMAIN}"
	echo "ZONE_SOA_MASTER 	=	${ZONE_SOA_MASTER}"
	echo "NSUPDATE_SIG0_KEY 	=	${NSUPDATE_SIG0_KEY}"
	echo
	echo "DEBUG: nsupdate settings"
	echo "  NSUPDATE_PARAM            = ${NSUPDATE_PARAM}"
	echo "  NSUPDATE_SET_SERVER       = ${NSUPDATE_SET_SERVER}"
	echo "  NSUPDATE_ACTION           = ${NSUPDATE_ACTION}"
	echo "  NSUPDATE_PRECONDITION_SET = ${NSUPDATE_PRECONDITION_SET} (not implemented for ${SCRIPT_NAME})"
	echo 
	echo "DEBUG: nsupdate commands to send"
	echo
	echo "NSUPDATE_SET_SERVER       = ${NSUPDATE_SET_SERVER}"
	echo "-- NSUPDATE_INPUT"
	echo -e "${NSUPDATE_INPUT}"
	echo "--"
fi



# send DNS updates via nsupdate
#

echo -e "${NSUPDATE_INPUT}" | nsupdate ${NSUPDATE_PARAM}

if [[ -n ${DEBUG} ]]; then
	echo
	echo "Browsable services"
	dig ${DIG_QUERY_PARAM} +short _services._dns-sd._udp.${DNSSD_DOMAIN} PTR
fi
