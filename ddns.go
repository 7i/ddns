package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v2"
)

var debug bool

type Config struct {
	// Your domain name
	Domain string `yaml:"Domain"`
	// Your ddns update url eg. "https://dyndns.binero.se/nic/update?hostname=*.example.com"
	DdnsUrl  string `yaml:"DdnsUrl"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	// How many seconds in between checks
	Frequency uint `yaml:"Frequency"`
	Debug     bool `yaml:"Debug"`
}

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

	client := &http.Client{}
	req, err := http.NewRequest("GET", c.DdnsUrl, nil)
	if err != nil {
		log.Fatalln("Error while constructing request:", err)
	}

	if c.Frequency < 1 { // Default value is 60s update check
		c.Frequency = 60
	}

	for {
		ip := getExternalIP()
		ips, err := net.LookupIP(c.Domain)
		if ip == "" || len(ips) == 0 {
			if debug {
				fmt.Println("DNS Error: ", err)
			}
			continue
		}
		set := true
		for _, i := range ips {
			if i.String() == ip {
				set = false
			}
		}
		if set {
			setIP(client, req, c)
			if debug {
				log.Println("New IP for", c.Domain, ":", ip)
			}
		}
		time.Sleep(time.Duration(c.Frequency) * time.Second)
	}
}

func setIP(client *http.Client, req *http.Request, c Config) {
	req.SetBasicAuth(c.Username, c.Password)
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
