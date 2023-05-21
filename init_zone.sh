#!/bin/bash
#
NEW_ZONE="testzone.zenr.io"
ZONEFILE_PATH="/var/cache/bind/dynamic/${NEW_ZONE}" # def Debian
NEW_ZONEFILE="named.${NEW_ZONE}"
KEY_ALGORITHM="RSASHA256"

mkdir -p ${ZONEFILEPATH}
chown -R bind:bind ${ZONEFILEPATH}
cd ${ZONEFILEPATH}

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
$SALT=`head -c 1000 /dev/random | sha1sum | cut -b 1-16`
dnssec-signzone -3 ${SALT} -A -N INCREMENT -o ${NEW_ZONE} -t ${NEW_ZONEFILE}
