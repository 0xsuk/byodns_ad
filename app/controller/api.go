package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/0xsuk/byodns/app/model"
	"github.com/0xsuk/byodns/util"
	"github.com/gorilla/mux"
)

//interact with model stuff
func processQuery(action Action, w http.ResponseWriter, r *http.Request) {
	switch action {
	case READ:
		if r.Method != "GET" {
			w.WriteHeader(400)
			return
		}
		queries := model.ReverseQuery(model.Queries)
		jqueries, err := json.Marshal(queries)
		if err != nil {
			util.Fatalln(err)
		}
		w.Write(jqueries)

	case SEARCH:
		if r.Method != "GET" {
			w.WriteHeader(400)
			return
		}
		util.Println("searching query")
		s := []model.Query{}
		for _, query := range model.Queries {
			if strings.Contains(query.Domain, r.FormValue("value")) {
				s = append(s, query)
			}
		}
		//Fix
		s = model.ReverseQuery(s)
		js, err := json.Marshal(s)
		if err != nil {
			util.Fatalln(err)
		}
		w.Write(js)

	default:
		w.Write([]byte("No such action exists"))
		w.WriteHeader(400)

	}
}

func processBlacklist(action Action, w http.ResponseWriter, r *http.Request) {
	switch action {
	case CREATE:
		//TODO input check
		if r.Method != "POST" {
			w.WriteHeader(400)
			return
		}
		v := r.FormValue("value")
		if v == "" {
			w.WriteHeader(400)
			return
		}
		util.Println("adding blacklist:", v)
		domains := model.AddBlacklist(v)
		domains = model.ReverseSlice(domains)
		jdomains, err := json.Marshal(domains)
		if err != nil {
			util.Fatalln(err)
		}
		w.Write(jdomains)

	case READ:
		if r.Method != "GET" {
			w.WriteHeader(400)
			return
		}
		util.Println("reading blacklist")
		domains := model.ReverseSlice(model.Blacklist) //issue: Blacklist reversed
		jdomains, err := json.Marshal(domains)
		if err != nil {
			util.Fatalln(err)
		}
		w.Write(jdomains)

	case UPDATE:
		if r.Method != "PUT" {
			w.WriteHeader(400)
			return
		}
		new := r.FormValue("value")
		old := r.FormValue("old")
		if new == "" || old == "" {
			w.WriteHeader(400)
			return
		}
		util.Println("updateing from", old, "to", new)
		domains := model.UpdateBlacklist(new, old)
		domains = model.ReverseSlice(domains)
		jdomains, err := json.Marshal(domains)
		if err != nil {
			util.Fatalln(err)
		}
		w.Write(jdomains)

	case DELETE:
		//TODO input check
		if r.Method != "DELETE" {
			w.WriteHeader(400)
			return
		}
		v := r.FormValue("value")
		if v == "" {
			w.WriteHeader(400)
			return
		}
		domains := model.DeleteBlacklist(v)
		domains = model.ReverseSlice(domains)
		jdomains, err := json.Marshal(domains)
		if err != nil {
			util.Fatalln(err)
		}
		w.Write(jdomains)

	case SEARCH:
		if r.Method != "GET" {
			w.WriteHeader(400)
			return
		}
		v := r.FormValue("value")
		if v == "" {
			js, err := json.Marshal(model.Blacklist)
			if err != nil {
				util.Fatalln(err)
			}
			w.Write(js)
			return
		}
		s := []string{}
		for _, domain := range model.Blacklist {
			if strings.Contains(domain, v) {
				s = append(s, domain)
			}
		}
		s = model.ReverseSlice(s)
		js, err := json.Marshal(s)
		if err != nil {
			util.Fatalln(err)
		}
		w.Write(js)

	default:
		w.Write([]byte("No such action exists"))
		w.WriteHeader(400)
	}
}

func processWhitelist(action Action, w http.ResponseWriter, r *http.Request) {

}

//interact with wrinting stuff
func apiHandler(w http.ResponseWriter, r *http.Request) {
	util.Println("HTTP request from", r.RemoteAddr)

	obj := Obj(mux.Vars(r)["obj"])
	action := Action(mux.Vars(r)["action"])

	switch obj {
	case QUERY:
		processQuery(action, w, r)
	case BLACKLIST:
		processBlacklist(action, w, r)
	case WHITELIST:
		processWhitelist(action, w, r)
	// case STATS:
	// 	processStats(action, w, r)
	default:
		util.Println("Unsupported obj:", obj)
		w.WriteHeader(400)
	}

}
