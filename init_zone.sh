#!/bin/bash
#
NEW_SUBZONE=${NEW_SUBZONE:-"_test"}
PARENT_ZONE=${PARENT_ZONE:-"zenr.io"}

NEW_ZONE_SERVER_FQDN=${NEW_ZONE_SERVER_FQDN:-"dns-oarc.free2air.net"}
NEW_ZONE_SERVER_IP=${NEW_ZONE_SERVER_IP:-"128.140.34.230"}
NEW_ZONE_CONTACT=${NEW_ZONE_CONTACT:-"root.free2air.org."}



NEW_ZONE="${NEW_SUBZONE}.${PARENT_ZONE}"
ZONEFILE_PATH="/var/cache/bind/dynamic/${NEW_ZONE}" # default in Debian
ZONEFILE_PATH_OWNER="bind:bind"                     # default in Debian
NEW_ZONEFILE="named.${NEW_ZONE}"
KEY_ALGORITHM="RSASHA256"

echo "NEW_SUBZONE = ${NEW_SUBZONE}"
echo "PARENT_ZONE = ${PARENT_ZONE}"
echo "NEW_ZONE = ${NEW_ZONE}"
echo "ZONEFILE_PATH = ${ZONEFILE_PATH}"

# check path exists
mkdir -p ${ZONEFILE_PATH}
chown -R bind:bind ${ZONEFILE_PATH}
cd ${ZONEFILE_PATH}

cat <<EOF >${NEW_ZONEFILE} 
\$ORIGIN .
\$TTL 360        ; 6 minutes

${NEW_ZONE}                IN SOA  ${NEW_ZONE_SERVER_FQDN}. ${NEW_ZONE_CONTACT} (
                                2006243476 ; serial
                                10800      ; refresh (3 hours)
                                900        ; retry (15 minutes)
                                604800     ; expire (1 week)
                                86400      ; minimum (1 day)
                                )
\$TTL 600
                        NS      ${NEW_ZONE_SERVER_FQDN}.
                        A       ${NEW_ZONE_SERVER_IP}
EOF


# create zone signing key (ZSK)
#
dnssec-keygen -a ${KEY_ALGORITHM} -b 2048 -n ZONE ${NEW_ZONE}

# create key signing key (KSK)
#
dnssec-keygen -f KSK -a ${KEY_ALGORITHM} -b 4096 -n ZONE ${NEW_ZONE}

#  add the ZSK & KSK public keys to zone
#
for key in `ls K${NEW_ZONE}*.key`
do
echo "\$INCLUDE $key">> ${NEW_ZONEFILE}
done

# sign zone
#
SALT=`head -c 1000 /dev/random | sha1sum | cut -b 1-16`
dnssec-signzone -3 ${SALT} -A -N INCREMENT -o ${NEW_ZONE} -t ${NEW_ZONEFILE}

cd -
