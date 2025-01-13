package verboten

import (
	"net/http"
	"net/url"

	"github.com/reiver/badgerverse.app/srv/http"
	"github.com/reiver/badgerverse.app/srv/log"
)

const path string = "/profile"

const queryName string = "who"

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
	if nil == request.URL {
		const code int = http.StatusInternalServerError
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("nil http-request URL")
		return
	}

	var who string
	{
		var query url.Values = request.URL.Query()
		if len(query) <= 0 {
			const code int = http.StatusNotFound
			http.Error(responsewriter, http.StatusText(code), code)
			log.Error("empty http-request URL query parameters")
			return
		}

		who = query.Get(queryName)
		if "" == who {
			const code int = http.StatusNotFound
			http.Error(responsewriter, http.StatusText(code), code)
			log.Errorf("empty http-request URL query parameter %q", queryName)
			return
		}
	}
	log.Informf("who: %q", who)

	{
		_, err := responsewriter.Write([]byte("who: "+who))
		if nil != err {
			log.Errorf("problem writing http-response to http-client: %s", err)
		}
	}
}
