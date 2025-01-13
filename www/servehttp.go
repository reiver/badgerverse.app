package verboten

import (
	"net/http"

	"github.com/reiver/badgerverse.app/srv/http"
	"github.com/reiver/badgerverse.app/srv/log"
)

const path string = "/"

func init() {
	err := httpsrv.Mux.HandlePath(http.HandlerFunc(serveHTTP), path)
	if nil != err {
		panic(err)
	}
}

func serveHTTP(responsewriter http.ResponseWriter, request *http.Request) {
	log := logsrv.Prefix("www("+path+")").Begin()
	defer log.End()

	if nil == responsewriter {
		return
	}
	if nil == request {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("nil http-request")
		return
	}

	{
		_, err := responsewriter.Write(webpage)
		if nil != err {
			log.Errorf("problem writing http-response to http-client: %s", err)
		}
	}
}
