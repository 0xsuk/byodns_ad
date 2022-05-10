package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/0xsuk/byodns/app/model"
	"github.com/0xsuk/byodns/config"
	"github.com/0xsuk/byodns/util"
	"github.com/gorilla/mux"
)

type Action string
type Obj string

const (
	CREATE Action = "create"
	READ   Action = "read"
	UPDATE Action = "update"
	DELETE Action = "delete"
	SEARCH Action = "search"

	QUERY     Obj = "query"
	BLACKLIST Obj = "blacklist"
	WHITELIST Obj = "whitelist"
)

var templates = make(map[string]*template.Template)

type Page struct {
	Path string
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	//deep copy
	if err := templates["query"].Execute(w, Page{
		Path: "query",
	}); err != nil {
		util.Fatalln(err)
	}
}

func blacklistHandler(w http.ResponseWriter, r *http.Request) {

	if err := templates["blacklist"].Execute(w, Page{
		Path: "blacklist",
	}); err != nil {
		util.Fatalln(err)
	}
}
func whitelistHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates["whitelist"].Execute(w, Page{
		Path: "whitelist",
	}); err != nil {
		util.Fatalln(err)
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates["dashboard"].Execute(w, Page{
		Path: "",
	}); err != nil {
		util.Fatalln(err)
	}
}

func debugBlacklist(w http.ResponseWriter, r *http.Request) {
	jdomains, _ := json.Marshal(model.Blacklist)
	w.Write(jdomains)
}

func middleWare(h http.Handler) http.Handler {
	serverID, serverPass := config.Cfg.WebServer.ID, config.Cfg.WebServer.Pass
	if serverID == "" || serverPass == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID, clientPass, ok := r.BasicAuth()
		if !ok || clientID != serverID || clientPass != serverPass {
			w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func StartHTTPServer() {
	Port := config.Cfg.WebServer.Port
	util.Printf("Starting HTTP server on port %v ...", Port)

	r := mux.NewRouter()
	r.HandleFunc("/api/{obj}/{action}", apiHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("app/view/"))))
	r.HandleFunc("/query", queryHandler).Methods("GET")
	r.HandleFunc("/blacklist", blacklistHandler).Methods("GET")
	r.HandleFunc("/whitelist", whitelistHandler).Methods("GET")
	r.HandleFunc("/debug/blacklist", debugBlacklist)
	r.HandleFunc("/", baseHandler).Methods("GET")
	util.Fatalln(http.ListenAndServe(":"+Port, middleWare(r)))
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles("app/view/"+name+".html",
		"app/view/_header.html",
		"app/view/_footer.html",
		"app/view/_searchbar.html",
	)
	if err != nil {
		util.Fatalln(err)
	}
	return t
}

func init() {
	util.Println("Parsing templates...")
	templates["dashboard"] = loadTemplate("dashboard")
	templates["query"] = loadTemplate("query")
	templates["blacklist"] = loadTemplate("blacklist")
	templates["whitelist"] = loadTemplate("whitelist")
	util.Println("Parsed templates!")
}
