#/bin/bash
#
if [[ -n ${TEST_SEND_NSUPDATE} ]]; then
	source functions/get_soa.sh
	source functions/get_sig0_keyid.sh
fi

send_nsupdate() {
	# send_nsupdate zone sig0_auth_key_fqdn 
	#	sends updates via nsupdate, where
	#		zone:			required. the dns domain needed to send (reflected in the list of updates). this sets the server to send to via get_soa_master()
	#		nsupdate_items:		required. the nsupdate command lines to pass
	#		sig0_auth_key_fqdn: 	optional. the FQDN of the sig0 auth key to use with -k option of nsupdate command. uses get_sig0_keyid() 
	#		sig0_auth_keypath: 	optional, needed only wih sig0_auth_keyi, with global fallback default. path to search for sig0_keyid with get_sig0_keyid() 
	#
	# uses multiple environment variables as fallbacks
	#
	# local variables from function paramters
	local zone=${1:-"${ZONE}"}
	[[ ! -n ${zone} ]] && echo "Error: ${FUNCNAME[0]}(): requires a zone parmeter or ZONE environment variable to be set." && exit 1

	local nsupdate_items="${2}"

	local sig0_auth_fqdn=${3:-""}			# set sig0 authentication key FQDN with global var fallback default

	local sig0_auth_keypath=${4:-"${NSUPDATE_SIG0_KEYPATH}"}	# set keystore path with global var fallback default


	# local variables (derived)
	local zone_master="$( get_soa_master ${zone} )" 		# set zone_master to master field of most fine grained SOA record given for zone
	local sig0_auth_keyid
	get_sig0_keyid sig0_auth_keyid "${sig0_auth_fqdn}" "${sig0_auth_keypath}"
	# form SIG0 private auth key param
	local nsupdate_auth_param
	if [[ -n ${sig0_auth_keyid} ]]; then
		nsupdate_auth_param="-k ${sig0_auth_keypath}/${sig0_auth_keyid}"
	else
		nsupdate_auth_param=""
	fi
	
	local nsupdate_server="server ${zone_master}"
	local nsupdate_preamble=$(echo "debug")
	local nsupdate_postamble=$(echo "send";echo "quit")


	if [[ -n ${DEBUG} ]]; then
		echo
		echo "${FUNCNAME[0]}( zone='${zone}' nsupdate_items= sig0_auth_fqdn='${sig0_auth_fqdn}' sig0_auth_keypath=${sig0_auth_keypath} )"
		echo
		echo "${FUNCNAME[0]}(): zone_master          = ${zone_master}"
		echo
		echo "${FUNCNAME[0]}(): sig0_auth_keypath    = ${sig0_auth_keypath}"
		echo "${FUNCNAME[0]}(): sig0_auth_fqdn       = ${sig0_auth_fqdn}"
		echo "${FUNCNAME[0]}(): sig0_auth_keyid      = ${sig0_auth_keyid}"
		echo
		echo "${FUNCNAME[0]}():  nsupdate_auth_param = '${nsupdate_auth_param}'"

		echo
		echo "${FUNCNAME[0]}(): nsupdate commands to send"
		echo
		(echo "${nsupdate_premable}";echo "${nsupdate_server}"; echo "${nsupdate_items}";echo "${nsupdate_postamble}") | cat
	fi

	if [[ ! -n ${NSUPDATE_DISABLE} ]]; then
		[[ -n ${DEBUG} ]] && echo && echo "${FUNCNAME[0]}(): Sending zone SOA master server '${zone_master}' update requests via nsupdate..."
		(echo "${nsupdate_premable}";echo "${nsupdate_server}"; echo "${nsupdate_items}";echo "${nsupdate_postamble}") | nsupdate ${nsupdate_auth_param}
	else
		echo "${FUNCNAME[0]}(): Warning: sending updates via nsupdate is disabled. NSUPDATE_DISABLE=\"${NSUPDATE_DISABLE}\""
	fi
}

if [[ -n ${TEST_SEND_NSUPDATE} ]]; then
	# NSUPDATE_DISABLE="true"
	DEBUG="true"
	# NEW_SUBZONE="test"
	NSUPDATE_SIG0_KEYPATH="${PWD}/keystore"
	items=$(echo "update add test1 60 A 127.0.0.1";echo "update add test2 60 A 127.0.0.1")
	send_nsupdate "zenr.io" "${items}" "debug1.zenr.io"
fi
