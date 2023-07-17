#!/bin/bash
#
# Sends DNS-DS updates for domain registration and browsing records
#------------------------------------------------------------------------------
SCRIPT_NAME="${0#*/}"
# Source env file for project variable default values
ENV_FILE=${ENV_FILE:-".env"}
if [ -e ${ENV_FILE} ]
then
        echo "Sourcing ${PWD}/${ENV_FILE} ..."
        . ${ENV_FILE}
fi

# Source env file for script default values
if [ -e ${ENV_FILE}.${SCRIPT_NAME} ]
then
	echo "Sourcing ${PWD}/${ENV_FILE}.${SCRIPT_NAME} ..."
	. ${ENV_FILE}.${SCRIPT_NAME}
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

# Form NSUPDATE input
#

# form SIG0 private auth key param
NSUPDATE_SIG0_KEY="/home/vortex/src/great-dane/test_go/Kzembla.zenr.io.+015+23799"
if [[ -n ${NSUPDATE_SIG0_KEY} ]]; then
	NSUPDATE_AUTH_PARAM="-k ${NSUPDATE_SIG0_KEY}.key"
else
	NSUPDATE_AUTH_PARAM=""
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
NSUPDATE_ITEM_DNS_SD_LB="update ${NSUPDATE_ACTION} lb._dns-sd._udp.${DNSSD_DOMAIN} ${NSUPDATE_TTL} PTR ${DNSSD_DOMAIN}."
NSUPDATE_ITEM_DNS_SD_DB="update ${NSUPDATE_ACTION} db._dns-sd._udp.${DNSSD_DOMAIN} ${NSUPDATE_TTL} PTR ${DNSSD_DOMAIN}."
NSUPDATE_ITEM_DNS_SD_B="update ${NSUPDATE_ACTION}  b._dns-sd._udp.${DNSSD_DOMAIN} ${NSUPDATE_TTL} PTR ${DNSSD_DOMAIN}."
NSUPDATE_ITEM_DNS_SD_DR="update ${NSUPDATE_ACTION} dr._dns-sd._udp.${DNSSD_DOMAIN} ${NSUPDATE_TTL} PTR ${DNSSD_DOMAIN}."
NSUPDATE_ITEM_DNS_SD_R="update ${NSUPDATE_ACTION}  r._dns-sd._udp.${DNSSD_DOMAIN} ${NSUPDATE_TTL} PTR ${DNSSD_DOMAIN}."
if [[ -n ${DEBUG} ]]; then
	echo "---DEBUG"
	echo "SCRIPT_NAME		=	${SCRIPT_NAME}"
	echo "ZONE 			=	${ZONE}"
	echo "NEW_SUBZONE 		=	${NEW_SUBZONE}"
	echo "ZONE_SOA_MASTER 	=	${ZONE_SOA_MASTER}"
	echo
	echo "DEBUG: nsupdate settings"
	echo
	echo "  NSUPDATE_SET_SERVER       = ${NSUPDATE_SET_SERVER}"
	echo "  NSUPDATE_ACTION           = ${NSUPDATE_ACTION}"
	echo "  NSUPDATE_PRECONDITION_SET = ${NSUPDATE_PRECONDITION_SET}"
	echo 
	echo "DEBUG: nsupdate commands to send"
	echo
	echo "NSUPDATE_SET_SERVER       = ${NSUPDATE_SET_SERVER}"
	echo "NSUPDATE_PRECONDITION     = ${NSUPDATE_PRECONDITION}"
	echo "NSUPDATE_ITEM_DNS_SD_LB   = ${NSUPDATE_ITEM_DNS_SD_LB}"
	echo "NSUPDATE_ITEM_DNS_SD_DB   = ${NSUPDATE_ITEM_DNS_SD_DB}"
	echo "NSUPDATE_ITEM_DNS_SD_B    = ${NSUPDATE_ITEM_DNS_SD_B}"
	echo "NSUPDATE_ITEM_DNS_SD_DR   = ${NSUPDATE_ITEM_DNS_SD_DR}"
	echo "NSUPDATE_ITEM_DNS_SD_DR   = ${NSUPDATE_ITEM_DNS_SD_DR}"
fi

# send DNS updates via nsupdate
#
cat <<EOF | nsupdate ${NSUPDATE_AUTH_PARAM}
${NSUPDATE_SET_SERVER}
${NSUPDATE_PRECONDITION}
${NSUPDATE_ITEM_DNS_SD_LB}
${NSUPDATE_ITEM_DNS_SD_DB}
${NSUPDATE_ITEM_DNS_SD_B}
${NSUPDATE_ITEM_DNS_SD_DR}
${NSUPDATE_ITEM_DNS_SD_R}
send
quit
EOF
DIG_QUERY_PARAM="+short"
for word in b lb db dr r _services
do
	echo "dig ${DIG_QUERY_PARAM} PTR ${word}._dns-sd._udp.${DNSSD_DOMAIN}."
	dig ${DIG_QUERY_PARAM} PTR ${word}._dns-sd._udp.${DNSSD_DOMAIN}.
done
