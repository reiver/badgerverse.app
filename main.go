package main

import (
	"github.com/reiver/badgerverse.app/srv/log"
)

func main() {
	log := logsrv.Prefix("main")

	log.Inform("badgerverse.app ⚡")
	yell()

	log.Inform("Here we go…")
	webserve()
}
