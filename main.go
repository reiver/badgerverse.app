package main

import (
	"github.com/reiver/badgerverse.app/srv/log"
)

func main() {
	log := logsrv.Prefix("main")

	log.Inform("badgerverse.app ⚡")
	tell()

	log.Inform("Here we go…")
	webserve()
}
