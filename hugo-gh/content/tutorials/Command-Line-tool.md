+++
title = 'sig0namectl Command Line Usage'
date = 2024-06-29T14:17:22+02:00
draft = false
summary = 'Usage examples of the sig0namectl tool to query and update DNS resource records.'
+++

This section gives usage examples for the command line sig0namectl utility.

```
NAME:
   sig0namectl - sig0 name control - direct, secure dynamic DNS

USAGE:
   sig0namectl [global options] command [command options] 

COMMANDS:
   query, q        query <name>
   print-key, pk   add a task to the list
   update, u       
   requestKey, rk  requestKey <my.new.name>
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value                   [$SIG0_HOST]
   --zone value, -z value        the zone you want to update [$SIG0_ZONE]
   --server value, --srv value    [$SIG0_SERVER]
   --key-name value, --kn value  Kso.me.na.me.+aaa+bbbbb [$SIG0_SIG0_KEYFILES]
   --help, -h                    show help
```

# query command usage

```
vortex@siluvian2:~$
vortex@siluvian2:~$ sig0namectl --server doh.zenr.io --zone zenr.io  query  --type KEY zenr.io
[Querying] 25:zenr.io
(*dns.Msg)(0xc000024090)(;; opcode: QUERY, status: NOERROR, id: 6103
;; flags: qr aa rd; QUERY: 1, ANSWER: 3, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version 0; flags: do; udp: 1232

;; QUESTION SECTION:
;zenr.io.       IN       KEY

;; ANSWER SECTION:
zenr.io.        600     IN      KEY     512 3 15 UqXF3jrAR0GugKoJbVUebRDss8XnSbE2nfv6jv0pZuA=
zenr.io.        600     IN      RRSIG   KEY 8 2 600 20240821205305 20240807202747 29583 zenr.io. WF8H2U4qiyMCp86blpR409+IDxBr3xtkK6pKZKEQeT+7rXPcnaqagyAI0NMaoFsWdQ17dYdh7xPK8Wead2StmrtQJZ68GdzfhwyE/bfzm5j6EILBKTknwzgmvONUhwAh/rQ9Rx0qe9THlqnuU+0HdR5MvT3fpdu7WkjOYPu5Yj1yXYYAmVld16QSljsaGMttgQ9UInrsSek5oMMXxSZ0DOfOp3zDFlEWXEdxFO+Atk3YJk7YOx9ss4hVWOY5Dgns8VwqXg7QHs7klgdIyrm+tA2WtVukOE8eSzw9BkEs+TqaYdUQZDRy1EFJTW+yY9PbRxXZ+2M8C4+e6UWS9Ha8LQ==
zenr.io.        600     IN      RRSIG   KEY 13 2 600 20240821205305 20240807202747 36504 zenr.io. /FBZfXGQ4cuV98pL+OWnzT+onCn+N0sWtkmblgwaw4OH4nC8Y7szWTS8mJd4J5yHu39vcrlz6iaBV00Ri3uQEA==
)
vortex@siluvian2:~$
```

# update command usage


```
vortex@siluvian2:~$
vortex@siluvian2:~$ sig0namectl --kn keystore/Kzembla.zenr.io.+015+23799 --server doh.zenr.io --zone zenr.io --host test.zembla update 1.2.3.4
2024/06/24 16:00:12 -- Reading SIG(0) Keyfiles (dnssec-keygen format) --
2024/06/24 16:00:12 keystore/Kzembla.zenr.io.+015+23799.key import: zembla.zenr.io. 3600 1 25 512 3 15 duQIg/NgFjwsE8ZKUuXJUG2/NNFs4o4byuwnekT062U=
2024/06/24 16:00:12 -- Set dns.Msg Structure --
2024/06/24 16:00:12 -- Create, fill & attach SIG RR to dns.Msg Structure --
2024/06/24 16:00:12 -- Configure DoH client --
2024/06/24 16:00:12 -- Response from DNS server --
;; opcode: UPDATE, status: NOERROR, id: 3957
;; flags: qr; ZONE: 1, PREREQ: 0, UPDATE: 0, ADDITIONAL: 0

;; ZONE SECTION:
;zenr.io.       IN       SOA
vortex@siluvian2:~$
```

