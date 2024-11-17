package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		host, _ := os.Hostname()
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not read IP addresses: %s", err)
			return
		}

		var ips strings.Builder
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips.WriteString(ipnet.IP.String())
					ips.WriteString("(")
					ips.WriteString(ipnet.Network())
					ips.WriteString(")")
					ips.WriteString("<br/>")
				}
			}
		}

		fmt.Fprintf(w, "<p>Host: %s</p><p>IPs:<br/> %s</p>", host, ips.String())
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
