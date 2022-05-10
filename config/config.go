package config

import (
	"github.com/0xsuk/byodns/util"
	"gopkg.in/ini.v1"
)

const CONFIG = "/etc/byodns/config.ini"

type UpperDNS struct {
	IP   string
	Port string
}

type LocalDNS struct {
	Port     string
	Ipv4only bool
}

type WebServer struct {
	Port string
	ID   string
	Pass string
}

type Redis struct {
	IP   string
	Port string
	Pass string
	Db   int
}

type Cluster struct {
	//specified in MILLISECONDS
	Diff int64 //if domain was requested in less time than CLuster.Diff, byodns categorilze domain as belonging to organizer domain, otherwise categorize domain as organizer_domain
}

type Config struct {
	UpperDNS
	LocalDNS
	WebServer
	Redis
	Cluster
}

var Cfg *Config

func notEmpty(s string) string {
	if s == "" {
		util.Fatalln("Required configuration field is empty.\ncheck /etc/byodns/config.ini or run install.sh again")
	}

	return s
}

func init() {
	cfg, err := ini.Load(CONFIG)
	if err != nil {
		util.Fatalln(err, "\nTry running `bash ./install.sh`")
	}

	Cfg = &Config{
		UpperDNS: UpperDNS{
			IP:   notEmpty(cfg.Section("upperdns").Key("ip").String()),
			Port: notEmpty(cfg.Section("upperdns").Key("port").String()),
		},
		LocalDNS: LocalDNS{
			Port:     notEmpty(cfg.Section("localdns").Key("port").String()),
			Ipv4only: cfg.Section("localdns").Key("ipv4only").MustBool(true),
		},
		WebServer: WebServer{
			Port: notEmpty(cfg.Section("webserver").Key("port").String()),
			ID:   cfg.Section("webserver").Key("id").String(),
			Pass: cfg.Section("webserver").Key("pass").String(),
		},
		Redis: Redis{
			IP:   cfg.Section("redis").Key("ip").String(),
			Port: notEmpty(cfg.Section("redis").Key("port").String()),
			Pass: cfg.Section("redis").Key("pass").String(),
			Db:   cfg.Section("redis").Key("db").MustInt(0),
		},
		Cluster: Cluster{},
	}
	diff, err := cfg.Section("cluster").Key("diff").Int64()
	if err != nil {
		util.Fatalln("config.go [cluster]: diff must be int64.")
	}
	Cfg.Cluster.Diff = diff

	if Cfg.WebServer.ID == "" || Cfg.WebServer.Pass == "" {
		util.Println("[webserver] id or pass is empty. skipping basic auth")
	}
	if Cfg.LocalDNS.Ipv4only {
		util.Println("[localdns] ipv4only set to true. skipping resolution for ipv6 and others.")
	}

	util.Println("Initialized Configuration")

}
