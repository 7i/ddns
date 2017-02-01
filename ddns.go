package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v2"
)

var debug bool

type Config struct {
	// Your domain names
	Domains []string `yaml:"Domains"`
	// Your ddns update url eg. "https://dyndns.binero.se/nic/update?hostname="
	DdnsUrl  string `yaml:"DdnsUrl"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	// How many seconds in between checks
	Frequency uint `yaml:"Frequency"`
	Debug     bool `yaml:"Debug"`
}

/* Exampel ddns.conf:

Domains:
 - example.com
 - example2.com
DdnsUrl: "https://dyndns.binero.se/nic/update?hostname="
Username: "user"
Password: "asdf"
Frequency: 60 # Seconds between checks
Debug: false # Print out debug info
*/
func main() {
	path := flag.String("config", "/etc/ddns.conf", "path to config file, eg. /home/user/ddns.conf")
	v := flag.Int("v", -1, "Print out debug messages, 0=no debug, 1=debug on")
	flag.Parse()

	conf, err := ioutil.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	var c Config
	err = yaml.Unmarshal(conf, &c)
	if err != nil {
		log.Fatal(err)
	}

	if *v != -1 { // command line takes presendence over config file
		if *v == 0 {
			debug = false
		} else {
			debug = true
		}
	} else {
		debug = c.Debug // default false
	}

	if c.Frequency < 1 { // Default value is 60s update check
		c.Frequency = 60
	}

	go updateService(c)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "exit" || scanner.Text() == "quit" || scanner.Text() == "q" {
			break
		}
	}
}

func updateService(c Config) {
	if len(c.Domains) < 1 {
		log.Fatalln("No domains in config file")
	}
	if debug {
		for _, domain := range c.Domains {
			fmt.Println("Starting updateService for " + domain)
		}
	}

	for {
		externalIP := getExternalIP()

		for _, domain := range c.Domains {
			ips, err := net.LookupIP(domain)
			if externalIP == "" || len(ips) == 0 {
				if debug {
					fmt.Println("DNS Error: ", err)
				}
				continue
			}
			set := true
			for _, i := range ips {
				if i.String() == externalIP {
					set = false
					break
				}
			}
			if set {
				mainDomain, err := http.NewRequest("GET", c.DdnsUrl+domain, nil)
				if err != nil {
					log.Println("Error while constructing request:", err)
				} else {
					setIP(mainDomain, c.Username, c.Password)
					if debug {
						log.Println("New IP for", domain, ":", externalIP)
					}
				}

				subDomains, err := http.NewRequest("GET", c.DdnsUrl+"*."+domain, nil)
				if err != nil {
					log.Println("Error while constructing request:", err)
				} else {
					setIP(subDomains, c.Username, c.Password)
					if debug {
						log.Println("New IP for", "*."+domain, ":", externalIP)
					}
				}
			}
		}
		time.Sleep(time.Duration(c.Frequency) * time.Second)
	}
}

func setIP(req *http.Request, username, password string) {
	client := &http.Client{}
	req.SetBasicAuth(username, password)
	res, err := client.Do(req)
	if err != nil && debug {
		log.Println(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil && debug {
		log.Println(err)
	} else if debug {
		log.Println("Responce from ddns request:\n", string(resBody))
	}
}

func getExternalIP() (ip string) {
	target := "myip.opendns.com"
	server := "resolver1.opendns.com"

	cl := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)
	r, _, err := cl.Exchange(&m, server+":53")
	if err != nil && debug {
		log.Println("DNS error:", err)
	} else if len(r.Answer) == 0 && debug {
		log.Println("No results from DNS request")
	} else {
		ip = fmt.Sprintf("%s", r.Answer[0].(*dns.A).A)
	}
	return ip
}
