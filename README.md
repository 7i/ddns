# ddns
Simple ddns updater to update the current ip for a domain.
Will check external IP, works for servers behind NAT.
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
Domain : "example.com" 
DdnsUrl : "https://dyndns.binero.se/nic/update?hostname=*.example.com"
Username : "username"
Password : "password"
Frequency : 60 # Seconds between checks
Debug : false # Print out debug info 
```
