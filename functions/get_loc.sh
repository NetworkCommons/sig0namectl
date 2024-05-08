#!/bin/bash
function get_loc {
	HOST_OS="`uname -o`"
	case "${HOST_OS}" in
		"Android")
			LOC_PROVIDER="termux"
			;;
		"GNU/Linux")
			LOC_PROVIDER="gpsd"
			;;
		*)
			LOC_PROVIDER="none"
			;;
	esac

	case "${LOC_PROVIDER}" in
		"termux")
			LOC_JSON=`termux-location -p gps`
			#
			LAT=`echo $LOC_JSON | jq -r '.latitude'`
			LON=`echo $LOC_JSON | jq -r '.longitude'`
			ALT=`echo $LOC_JSON | jq -r '.altitude'`
			HACC=`echo $LOC_JSON | jq -r '.accuracy'`
			VACC=`echo $LOC_JSON | jq -r '.vertical_accuracy'`
			BEARING=`echo $LOC_JSON | jq -r '.bearing'`
			SPEED=`echo $LOC_JSON | jq -r '.speed'`
			#
			LATLON_DEG="${LAT} ${LON}"
			LATLON_DMS=`echo "${LATLON_DEG}" | GeoConvert -d -:`
			# echo "${LATLON_DMS}"
			RR_CONV="${LATLON_DMS//\:/\ }" 
			RR_CONV="${RR_CONV/N/\ N}" 
			RR_CONV="${RR_CONV/S/\ S}" 
			RR_CONV="${RR_CONV/E/\ E}" 
			RR_CONV="${RR_CONV/W/\ W}" 
			
			# 2 decimal places for alt, hacc, vacc
			ALT2="$(printf "%8.2f" ${ALT})"
			HACC2="$(printf "%8.2f" ${HACC})"
			VACC2="$(printf "%8.2f" ${VACC})"
			#
			LOC_RR_DATA="${RR_CONV} 0.00m ${ALT2}m ${HACC2}m ${VACC2}m"
			;;

		"gpsd")
			# pindrop --lon --lat --alt --error -v
			#
			LAT="$(pindrop --lat)"
			LON="$(pindrop --lon)"
			ALT="$(pindrop --alt)"
			HACC="$(pindrop --error | tr "\'" '\"' | jq -r '.s')"
			VACC="$(pindrop --error | tr "\'" '\"' | jq -r '.v')"
			#
			LATLON_DEG="${LAT} ${LON}"
			LATLON_DMS=`echo "${LATLON_DEG}" | GeoConvert -d -:`
			# echo "${LATLON_DMS}"
			RR_CONV="${LATLON_DMS//\:/\ }" 
			RR_CONV="${RR_CONV/N/\ N}" 
			RR_CONV="${RR_CONV/S/\ S}" 
			RR_CONV="${RR_CONV/E/\ E}" 
			RR_CONV="${RR_CONV/W/\ W}" 
			
			# 2 decimal places for alt, hacc, vacc
			ALT2="$(printf "%8.2f" ${ALT})"
			HACC2="$(printf "%8.2f" ${HACC})"
			VACC2="$(printf "%8.2f" ${VACC})"
			#
			LOC_RR_DATA="${RR_CONV} 0.00m ${ALT2}m ${HACC2}m ${VACC2}m"
			;;
		*)
			;;
	esac

	echo ${LOC_RR_DATA}
}
