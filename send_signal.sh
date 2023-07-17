#!/bin/bash
#
# Send signal "requests" to _signal.<zone> ie CDS, CDNSKEY, and KEY records to zone master.
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

# load helpful functions
for i in functions/*.sh
do
        . ${i}
        [[ -n ${DEBUG} ]] && echo "Sourced ${PWD}/functions/$i ..."
done

# Default existing ZONE fallback to $DOMAINNAME if set, else to domain set in $HOSTNAME, else error
# ZONE MUST correspond to an existing DNS ZONE, that is resolve with an SOA record
ZONE=${ZONE:-${DOMAINNAME}}
ZONE=${ZONE:-${HOSTNAME#*.}}
if [[ ! -n ${ZONE} ]]; then
        echo "Error: DNS zone \$ZONE environment variable is not set & could not be determined from \$DOMAINNAME or \$HOSTNAME"
        exit 1
fi

DIG_QUERY_PARAM=${DIG_QUERY_PARAM:-}
AVAHI_BROWSE_PARAM=${AVAHI_BROWSE_PARAM:-"-brat"}

# Discover master (usually primary DNS server) of zone from master field of SOA record
#
ZONE_SOA_MASTER=$( get_soa_master ${ZONE} )
if [[ ! -n ${ZONE_SOA_MASTER} ]]; then
        echo "Warning: ZONE ${ZONE} SOA record does not resolve"
        exit 1
fi

#------------------------------------------------------------------------------

# Define zone to install on local BIND server
#
if [[ ! -n ${NEW_SUBZONE} ]]; then
        echo "Error: NEW_SUBZONE ${NEW_SUBZONE} environment variable is undefined."
fi
NEW_SUBZONE=${NEW_SUBZONE:-"testzone"}
NEW_ZONE="${NEW_SUBZONE}.${ZONE}"

# A check should be made over whether the requested domain RRs already exist and warn the requester

# test
NSUPDATE_SIG0_KEY="/home/vortex/src/great-dane/test_go/Kzembla.zenr.io.+015+23799"
NEW_SUBZONE_UPDATE_KEY="`cat ${NSUPDATE_SIG0_KEY}.key`"

UPDATE_TTL=${UPDATE_TTL:-"60"}

NSUPDATE_SET_SERVER="server ${ZONE_SOA_MASTER}"

NSUPDATE_ACTION=${NSUPDATE_ACTION:-"add"}

NSUPDATE_ITEM_SIG0_KEY="update ${NSUPDATE_ACTION} _dsboot.${NEW_SUBZONE}._signal.${ZONE} ${UPDATE_TTL} ${NEW_SUBZONE_UPDATE_KEY##*.}"
NSUPDATE_ITEM_CDS=""
NSUPDATE_ITEM_CDNSKEY=""
NSUPDATE_ITEM_PTR="update ${NSUPDATE_ACTION} _signal.${ZONE} ${UPDATE_TTL} PTR ${NEW_SUBZONE}._signal.${ZONE}."

if [[ -n ${DEBUG} ]]; then
	echo "ZONE              =	${ZONE}"
	echo "ZONE_SOA_MASTER   =	${ZONE_SOA_MASTER}"
	echo "NEW_SUBZONE   =	${NEW_SUBZONE}"
	
	echo
	echo "SUBZONE Requested =	${NEW_SUBZONE}"
	echo "  DNS Update KEY pair =	${NSUPDATE_SIG0_KEY}"
	echo "  DNS Update KEY file =	${NEW_SUBZONE_UPDATE_KEY}"
	
	echo
	echo "  NSUPDATE_SET_SERVER =	${NSUPDATE_SET_SERVER}"
	echo "  NSUPDATE_ACTION = ${NSUPDATE_ACTION}"
	echo "  NSUPDATE_ITEM_SIG0_KEY = ${NSUPDATE_ITEM_SIG0_KEY}"
	echo "  NSUPDATE_ITEM_PTR = ${NSUPDATE_ITEM_PTR}"
	
	echo
	echo "nsupdate commands to send"
	echo 
	echo "${NSUPDATE_SET_SERVER}"
	echo "${NSUPDATE_ITEM_SIG0_KEY}"
	echo "${NSUPDATE_ITEM_CDS}"
	echo "${NSUPDATE_ITEM_CDNSKEY}"
	echo "${NSUPDATE_ITEM_PTR}"
	echo "send"
	echo "quit"
fi

echo
echo "Sending zone master (${ZONE_SOA_MASTER}) 'update ${NSUPDATE_ACTION}' resource record requests via nsupdate..."
cat <<EOF | nsupdate
${NSUPDATE_SET_SERVER}
${NSUPDATE_ITEM_SIG0_KEY}
${NSUPDATE_ITEM_CDNSKEY}
${NSUPDATE_ITEM_CDS}
${NSUPDATE_ITEM_PTR}
send
quit
EOF

echo
echo "Validate current entres via 'dig' ..."
echo "KEY     `dig @${ZONE_SOA_MASTER} +short ${NEW_SUBZONE}._signal.${ZONE} KEY`"
echo "CDNSKEY `dig @${ZONE_SOA_MASTER} +short ${NEW_SUBZONE}._signal.${ZONE} CDNSKEY`"
echo "CDS     `dig @${ZONE_SOA_MASTER} +short ${NEW_SUBZONE}._signal.${ZONE} CDS`"
echo "PTR     `dig @${ZONE_SOA_MASTER} +short _signal.${ZONE}. PTR`"
# echo "Validate entries via avahi-browse"
# echo "KEY     `dig @${ZONE_SOA_MASTER} +short _dsboot.${NEW_SUBZONE}._signal.${ZONE} KEY`"
# echo "CDNSKEY `dig @${ZONE_SOA_MASTER} +short _dsboot.${NEW_SUBZONE}._signal.${ZONE} CDNSKEY`"
# echo "CDS     `dig @${ZONE_SOA_MASTER} +short _dsboot.${NEW_SUBZONE}._signal.${ZONE} CDS`"
# echo "PTR     `dig @${ZONE_SOA_MASTER} +short _signal.${ZONE}. PTR`"
