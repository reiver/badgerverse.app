package verboten

import (
	htmltemplate "html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/reiver/badgerverse.app/lib/label"
	"github.com/reiver/badgerverse.app/srv/demo"
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

	var tmplParams struct {
		COVERIMAGE htmltemplate.URL
		ICONIMAGE htmltemplate.URL
		LABELS []label.Label
		NAME string
		SUMMARY string
		WHO string
	}
	switch strings.ToLower(who) {
	case demosrv.REIVER_WHO:
		tmplParams.COVERIMAGE = htmltemplate.URL(demosrv.REIVER_COVERIMAGE)
		tmplParams.ICONIMAGE  = htmltemplate.URL(demosrv.REIVER_ICONIMAGE)
		tmplParams.LABELS     = demosrv.REIVER_LABELS
		tmplParams.NAME       = demosrv.REIVER_NAME
		tmplParams.SUMMARY    = demosrv.REIVER_SUMMARY
		tmplParams.WHO        = demosrv.REIVER_WHO
	default:
		const code int = http.StatusNotFound
		http.Error(responsewriter, http.StatusText(code), code)
		log.Error("TODO: fetch profile")
		return
	}

	{
		err := template.Execute(responsewriter, tmplParams)
		if nil != err {
			log.Errorf("problem rendering template to HTML and writing http-response to http-client: %s", err)
		}
	}
}
