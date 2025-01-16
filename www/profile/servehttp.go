package verboten

import (
	htmltemplate "html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/reiver/go-fediverseid"
	mstdnent    "github.com/reiver/go-mstdn/ent"
	mstdnlookup "github.com/reiver/go-mstdn/api/v1/accounts/lookup"

	"github.com/reiver/badgerverse.app/lib/label"
	"github.com/reiver/badgerverse.app/srv/acct"
	"github.com/reiver/badgerverse.app/srv/demo"
	"github.com/reiver/badgerverse.app/srv/html"
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

	var fediverseID fediverseid.FediverseID
	{
		var err error
		fediverseID, err = fediverseid.ParseFediverseID(who)
		if nil != err {
			const code int = http.StatusBadRequest
			http.Error(responsewriter, http.StatusText(code), code)
			log.Errorf("problem parsing fediverse-id %q: %s", who, err)
			return
		}
	}

	var acctname string
	var accthost string
	{
		acctname = fediverseID.NameElse("")

		accthost = fediverseID.HostElse("")
		if "" == accthost {
			const code int = http.StatusBadRequest
			http.Error(responsewriter, http.StatusText(code), code)
			log.Errorf("empty fediverse-id host in %q", who)
			return
		}
	}
	log.Informf("fediverse-id name: %q", acctname)
	log.Informf("fediverse-id host: %q", accthost)

	var tmplParams struct {
		COVERIMAGE htmltemplate.URL
		ICONIMAGE htmltemplate.URL
		LABELS []label.Label
		NAME string
		SUMMARY htmltemplate.HTML
		WHO string
	}
	switch strings.ToLower(who) {
	case demosrv.REIVER_WHO:
		tmplParams.COVERIMAGE = htmltemplate.URL(demosrv.REIVER_COVERIMAGE)
		tmplParams.ICONIMAGE  = htmltemplate.URL(demosrv.REIVER_ICONIMAGE)
		tmplParams.LABELS     = demosrv.REIVER_LABELS
		tmplParams.NAME       = demosrv.REIVER_NAME
		tmplParams.SUMMARY    = htmltemplate.HTML(demosrv.REIVER_SUMMARY)
		tmplParams.WHO        = demosrv.REIVER_WHO
	default:
		var account mstdnent.Account

		err := mstdnlookup.Get(&account, accthost, acctname)
		if nil != err {
			const code int = http.StatusNotFound
			http.Error(responsewriter, http.StatusText(code), code)
			log.Errorf("problem with GET lookup to %q %q: %s", accthost, acctname, err)
			return
		}

		tmplParams.COVERIMAGE = htmltemplate.URL(acctsrv.DEFAULT_COVERIMAGE)
		account.Header.WhenSomething(func(value string){
			if "" != value {
				tmplParams.COVERIMAGE = htmltemplate.URL(value)
			}
		})
		tmplParams.ICONIMAGE  = htmltemplate.URL(acctsrv.DEFAULT_ICONIMAGE)
		account.Avatar.WhenSomething(func(value string){
			if "" != value {
				tmplParams.ICONIMAGE = htmltemplate.URL(value)
			}
		})
					
		tmplParams.LABELS     = nil //@TODO
					
		tmplParams.NAME       = account.DisplayName.GetElse("")
		tmplParams.SUMMARY    = htmltemplate.HTML(htmlsrv.SanitizeString(account.Note.GetElse("")))
		tmplParams.WHO        = who
	}

	{
		err := template.Execute(responsewriter, tmplParams)
		if nil != err {
			log.Errorf("problem rendering template to HTML and writing http-response to http-client: %s", err)
		}
	}
}
