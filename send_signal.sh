#!/bin/bash
#
# Send signal "requests" to _signal.<zone> ie CDS, CDNSKEY, and KEY records to zone master.
#------------------------------------------------------------------------------
echo
echo "Executing ${PWD}/$0 on host ${HOSTNAME}"

# Source env file for variable default values
ENV_FILE=${ENV_FILE:-".env"}
if [ -e ${ENV_FILE} ]
then
        echo "Sourcing ${PWD}/${ENV_FILE} ..."
        . ${ENV_FILE}
fi

# Default existing parent ZONE fallback to $DOMAINNAME if set, else to domain set in $HOSTNAME, else error
#
ZONE=${ZONE:-${DOMAINNAME}}
ZONE=${ZONE:-${HOSTNAME#*.}}
if [[ ! -n ${ZONE} ]]; then
        echo "Error: Parent DNS zone \$ZONE environment variable not set in command line, sourced in ${ENV_FILE} or determined from \$DOMAINNAME or \$HOSTNAME"
        exit 1
fi

# Discover master (usually primary DNS server) of parent zone from SOA record
#
DIG_QUERY_PARAM=${DIG_QUERY_PARAM:-}
ZONE_SOA_MASTER=${ZONE_SOA_MASTER:-$(dig ${DIG_QUERY_PARAM} +short ${ZONE} SOA | cut -f1 -d' ')}
if [[ ! -n ${ZONE_SOA_MASTER} ]]; then
        echo "Warning: Parent ZONE ${ZONE} SOA record does not resolve"
fi

# Define zone to install on local BIND server
#
NEW_SUBZONE=${NEW_SUBZONE:-"testzone"}
NEW_ZONE="${NEW_SUBZONE}.${ZONE}"

#------------------------------------------------------------------------------

# test
NEW_SUBZONE_UPDATE_KEYPAIR="/home/vortex/src/great-dane/test_go/Kzembla.zenr.io.+015+23799"
NEW_SUBZONE_UPDATE_KEY="`cat ${NEW_SUBZONE_UPDATE_KEYPAIR}.key`"

NEW_SUBZONE_TTL="600"
NSUPDATE_SET_SERVER="server ${ZONE_SOA_MASTER}"
NSUPDATE_ACTION=${NSUPDATE_ACTION:-"add"}
NSUPDATE_ITEM_UPDATE_KEY="update ${NSUPDATE_ACTION} _dsboot.${NEW_SUBZONE}._signal.${ZONE} ${NEW_SUBZONE_TTL} ${NEW_SUBZONE_UPDATE_KEY##*.}"
NSUPDATE_ITEM_CDS=""
NSUPDATE_ITEM_CDNSKEY=""
NSUPDATE_ITEM_PTR="update ${NSUPDATE_ACTION} _signal.${ZONE} ${NEW_SUBZONE_TTL} PTR ${NEW_SUBZONE}._signal.${ZONE}."

echo "ZONE              =	${ZONE}"
echo "ZONE_SOA_MASTER   =	${ZONE_SOA_MASTER}"

echo
echo "SUBZONE Requested =	${NEW_SUBZONE}"
echo "  DNS Update KEY pair =	${NEW_SUBZONE_UPDATE_KEYPAIR}"
echo "  DNS Update KEY file =	${NEW_SUBZONE_UPDATE_KEY}"

echo
echo "  NSUPDATE_SET_SERVER =	${NSUPDATE_SET_SERVER}"
echo "  NSUPDATE_ACTION = ${NSUPDATE_ACTION}"
echo "  NSUPDATE_ITEM_UPDATE_KEY = ${NSUPDATE_ITEM_UPDATE_KEY}"
echo "  NSUPDATE_ITEM_PTR = ${NSUPDATE_ITEM_PTR}"

echo
echo "nsupdate commands to send"
echo 
echo "${NSUPDATE_SET_SERVER}"
echo "${NSUPDATE_ITEM_UPDATE_KEY}"
echo "${NSUPDATE_ITEM_CDS}"
echo "${NSUPDATE_ITEM_CDNSKEY}"
echo "${NSUPDATE_ITEM_PTR}"
echo "send"

echo
echo "Sending zone master (${ZONE_SOA_MASTER}) 'update ${NSUPDATE_ACTION}' resource record requests via nsupdate..."
cat <<EOF | nsupdate
${NSUPDATE_SET_SERVER}
${NSUPDATE_ITEM_UPDATE_KEY}
${NSUPDATE_ITEM_CDNSKEY}
${NSUPDATE_ITEM_CDS}
${NSUPDATE_ITEM_PTR}
send
quit
EOF

echo
echo "Validate old entres ..."
echo "KEY     `dig @${ZONE_SOA_MASTER} +short ${NEW_SUBZONE}._signal.${ZONE} KEY`"
echo "CDNSKEY `dig @${ZONE_SOA_MASTER} +short ${NEW_SUBZONE}._signal.${ZONE} CDNSKEY`"
echo "CDS     `dig @${ZONE_SOA_MASTER} +short ${NEW_SUBZONE}._signal.${ZONE} CDS`"
echo "PTR     `dig @${ZONE_SOA_MASTER} +short _signal.${ZONE}. PTR`"
echo "Validate entres ..."
echo "KEY     `dig @${ZONE_SOA_MASTER} +short _dsboot.${NEW_SUBZONE}._signal.${ZONE} KEY`"
echo "CDNSKEY `dig @${ZONE_SOA_MASTER} +short _dsboot.${NEW_SUBZONE}._signal.${ZONE} CDNSKEY`"
echo "CDS     `dig @${ZONE_SOA_MASTER} +short _dsboot.${NEW_SUBZONE}._signal.${ZONE} CDS`"
echo "PTR     `dig @${ZONE_SOA_MASTER} +short _signal.${ZONE}. PTR`"
