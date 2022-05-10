package model

import (
	"database/sql"
	"sync"
	"time"

	"github.com/0xsuk/byodns/config"
	"github.com/0xsuk/byodns/util"
	"github.com/go-redis/redis"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DIR = "/etc/byodns"
	//db name
	DOMAINS = "domains.db"
)

var (
	DbDomains     *sql.DB
	RedisInstance *redis.Client
	Gravity       []string
	Blacklist     []string
	Whitelist     []string
	Queries       []Query
	mu            = &sync.Mutex{}
	//should be initialized from init() https://github.com/go-sql-driver/mysql/issues/150
	//used most often
	query_adder *sql.Stmt

	last_time        time.Time
	organizer_domain string = "first_organizer"

	//TODO if you change here, you have to modify rows.Scan
	table_gravity = sqlTable{
		TABLE: "gravity",
		READ:  "select domain from gravity",
		CREATE: `CREATE TABLE IF NOT EXISTS gravity
    			(
	   				domain TEXT NOT NULL 
    			)`,
	}
	table_query = sqlTable{
		TABLE:  "query",
		READ:   "select domain, clientip, timestamp, diff, organizer_domain, status, isblocked from query",
		INSERT: "insert into query (domain,clientip,timestamp,diff,organizer_domain,status,isblocked) values (?,?,?,?,?,?,?)",
		UPDATE: "update query set isblocked = ? where domain = ?",
		CREATE: `CREATE TABLE IF NOT EXISTS query
    			(
    				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    				domain TEXT NOT NULL,
    				clientip TEXT NOT NULL,
    				timestamp TEXT NOT NULL,
    				diff INT,
    				organizer_domain TEXT,
	  				status TEXT NOT NULL,
	  				isblocked TEXT NOT NULL
    			)`,
	}
	table_blacklist = sqlTable{
		TABLE:  "blacklist",
		READ:   "select domain from blacklist",
		INSERT: "insert into blacklist (domain) values (?)",
		UPDATE: "update blacklist set domain = ? where domain = ?",
		CREATE: `CREATE TABLE IF NOT EXISTS blacklist
    				(
     					domain TEXT NOT NULL 
    				)`,
		DELETE: "delete from blacklist where domain = ?",
	}
	table_whitelist = sqlTable{
		TABLE: "whitelist",
		READ:  "select domain from whitelist",
		CREATE: `CREATE TABLE IF NOT EXISTS whitelist
    			(
     				domain TEXT NOT NULL 
    			)`,
		INSERT: "insert into whitelist (domain) values (?)",
		UPDATE: "update whitelist set domain = ? where domain = ?",
		DELETE: "delete from whitelist where domain = ?",
	}
)

type sqlTable struct {
	TABLE  string
	READ   string
	INSERT string
	UPDATE string
	CREATE string
	DELETE string
}

type Query struct {
	Domain          string
	ClientIP        string
	Timestamp       string
	Diff            int64
	OrganizerDomain string
	Status          string
	IsBlocked       string //"yes" "no"
}

func initDomains() {
	var err error
	//No handling for path not exists
	DbDomains, err = sql.Open("sqlite3", DIR+"/"+DOMAINS)
	if err != nil {
		util.Fatalln(err)
	}

	var DomainsSchema = []sqlTable{table_gravity, table_query, table_blacklist, table_whitelist} //Create table if not exists
	for _, table := range DomainsSchema {
		_, err = DbDomains.Exec(table.CREATE)
		if err != nil {
			util.Fatalln(err)
		}
	}

	query_adder, _ = DbDomains.Prepare(table_query.INSERT)
}
func initRedis() {
	RedisInstance = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.IP + ":" + config.Cfg.Redis.Port,
		Password: config.Cfg.Redis.Pass,
		DB:       config.Cfg.Db,
	})
	err := RedisInstance.Ping().Err()
	if err != nil {
		util.Fatalln(err)
	}
	util.Println("Connected to", RedisInstance.Options().Addr)
}
func readDomains() {
	Gravity = ReadSliceFrom(table_gravity)
	Blacklist = ReadSliceFrom(table_blacklist)
	Whitelist = ReadSliceFrom(table_whitelist)
	Queries = ReadQuery()
}

func init() {
	initDomains()
	initRedis()
	readDomains()
}
