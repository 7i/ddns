# ddns
Simple ddns updater. 
Will update the IP for domains and subdomains to the current external IP.
Will check external IP via a DNS request to server resolver1.opendns.com with target myip.opendns.com.
Only supports basic auth.

# Install:
```
$ go get github.com/7i/ddns
```

# Example usage:
Move ddns.conf to desired location and adjust config path
```
$ ddns -config $GOPATH/src/github.com/7i/ddns/ddns.conf 
```

Use default config location /etc/ddns.conf.
```
$ ddns &
```

Specify your own location for your ddns.conf.
```
$ ddns -config /home/user/ddns.conf &
```

Turn on debugging.
```
$ ddns -v 1
```

Example ddns.conf:
```
Domains: # The following domains will be added at the end of the DdnsUrl string.
 - example.com
 - example2.com
DdnsUrl : "https://dyndns.binero.se/nic/update?hostname="
Username : "username"
Password : "password"
Frequency : 60 # Seconds between checks
Debug : false # Print out debug info 
```
