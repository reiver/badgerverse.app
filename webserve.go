package main

import (
	"net/http"

	"github.com/reiver/badgerverse.app/cfg"
	"github.com/reiver/badgerverse.app/srv/http"
	"github.com/reiver/badgerverse.app/srv/log"

        // This import enables all the HTTP handlers.
      _ "github.com/reiver/badgerverse.app/www"
)

func webserve() {
	log := logsrv.Prefix("webserve")

	var tcpaddr string = cfg.WebServerTCPAddress()
	log.Informf("serving HTTP on TCP address: %q", tcpaddr)

	err := http.ListenAndServe(tcpaddr, &httpsrv.Mux)
	if nil != err {
		log.Errorf("ERROR: problem with serving HTTP on TCP address %q: %s", tcpaddr, err)
		panic(err)
	}
}
