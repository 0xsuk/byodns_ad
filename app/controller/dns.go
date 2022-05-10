package controller

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/0xsuk/byodns/app/model"
	"github.com/0xsuk/byodns/config"
	"github.com/0xsuk/byodns/util"
	"github.com/miekg/dns"
)

func forEachQuestion(clientIP string, q dns.Question, msg *dns.Msg) {
	util.Printf("from %v: %v %v\n", clientIP, q.Name, q.Qtype)
	domain := q.Name[:strings.LastIndex(q.Name, ".")]
	var status string
	var qtype string

	if config.Cfg.LocalDNS.Ipv4only {
		if q.Qtype == dns.TypeA {
			qtype = "A"
		} else {
			return
		}
	} else {
		if q.Qtype == dns.TypeA {
			qtype = "A"
		} else if q.Qtype == dns.TypeAAAA {
			qtype = "AAAA"
		}
	}

	//took 0.002 sec
	//TODO qtype 65 https
	if model.IsBlocked(domain) {
		util.Println("BLOCKED!")
		//NOTE: always returning qtype A
		rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, "0.0.0.0"))
		if err != nil {
			util.Fatalln(err)
		}
		msg.Answer = []dns.RR{rr}
		status = "BLOCKED"
		go model.AddQuery(domain, clientIP, time.Now(), status, "yes")
		return
	}

	var dstIP string
	if q.Qtype == dns.TypeA { //for ipv4 req
		dstIP = model.GetDomainIPv4(domain)
	} else if q.Qtype == dns.TypeAAAA { //for ipv6 req
		dstIP = model.GetDomainIPv6(domain)
	}

	if dstIP != "" {
		util.Println(domain, "found in cache:", dstIP)
		rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", q.Name, qtype, dstIP))
		if err != nil {
			util.Fatalln(err)
		}
		msg.Answer = []dns.RR{rr}
		return
	}
	c := new(dns.Client)
	uppermsg := new(dns.Msg)

	uppermsg.SetQuestion(dns.Fqdn(domain), q.Qtype)
	uppermsg.RecursionDesired = true

	//resp, _, err := c.Exchange(uppermsg, net.JoinHostPort(config.Cfg.UpperDNS.IP, config.Cfg.UpperDNS.Port))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, _, err := c.ExchangeContext(ctx, uppermsg, net.JoinHostPort(config.Cfg.UpperDNS.IP, config.Cfg.UpperDNS.Port))

	if resp == nil {
		//fix invalid memory address or nil pointer dereference
		resp = new(dns.Msg)
	}

	if err != nil {
		//TODO deal with "read udp timeout"
		util.Println("\033[43m[WARNING]\033[00m: " + err.Error())
	}

	if resp.Rcode != dns.RcodeSuccess {
		util.Println("Domain Resolution Failed")
		status = "failed"
	} else {
		status = "ok"
	}

	for _, a := range resp.Answer {
		ans := strings.Split(a.String(), "\t")
		if len(ans) == 5 && ans[3] == qtype {
			// Save on cache
			ip := ans[4]
			ttl := time.Duration(a.Header().Ttl) * time.Second
			if q.Qtype == dns.TypeA {
				model.SetDomainIPv4(domain, ip, ttl)
			} else if q.Qtype == dns.TypeAAAA {
				model.SetDomainIPv6(domain, ip, ttl)
			}
		}
	}

	go model.AddQuery(domain, clientIP, time.Now(), status, "no")
	msg.Answer = resp.Answer

}

//parseQuery resolves question
func parseQuery(clientIP string, msg *dns.Msg) {
	//msg.Question contains slice of dns.Question
	///each took 0.05~0.3sec, but mostly its single dns.Question
	for _, q := range msg.Question {
		forEachQuestion(clientIP, q, msg)
	}
}

//handleDnsRequest writes response when dns question arrives
//handleDnsRequest is asyncronous already (probably)
func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	start := time.Now()
	msg := new(dns.Msg)
	msg.SetReply(r)
	clientIP := w.RemoteAddr().String()
	clientIP = clientIP[:strings.LastIndex(clientIP, ":")] //extract address before :

	switch r.Opcode {
	//for standard request
	case dns.OpcodeQuery:
		parseQuery(clientIP, msg)
	}

	err := w.WriteMsg(msg)
	if err != nil {
		util.Fatalln(err)
	}
	end := time.Now()
	util.Printf("TOOK: %fsec\n", end.Sub(start).Seconds())
}

//StartDNSServer starts on configured port
func StartDNSServer() {
	port := config.Cfg.LocalDNS.Port

	dns.HandleFunc(".", handleDnsRequest)

	server := &dns.Server{Addr: ":" + port, Net: "udp"}
	util.Println("Starting DNS on port", port)

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		server.Shutdown()
		util.Fatalln(err)
	}
}
