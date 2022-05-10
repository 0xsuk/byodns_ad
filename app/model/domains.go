package model

import (
	"time"

	"github.com/0xsuk/byodns/config"
	"github.com/0xsuk/byodns/util"
)

//file to interact with domain.db

func IsBlocked(domain string) bool {
	for _, blocked := range Gravity {
		if domain == blocked {
			return true
		}
	}
	for _, blocked := range Blacklist {
		if domain == blocked {
			return true
		}
	}
	return false
}

func ReverseSlice(a []string) []string {
	s := []string{}
	s = append(s, a...)
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}

	return s
}
func ReverseQuery(queries []Query) []Query {
	new_queries := []Query{}
	new_queries = append(new_queries, queries...)
	for i := len(new_queries)/2 - 1; i >= 0; i-- {
		opp := len(new_queries) - 1 - i
		new_queries[i], new_queries[opp] = new_queries[opp], new_queries[i]
	}

	return new_queries
}

//ReadSliceFrom reads single string from table and return []string
func ReadSliceFrom(table sqlTable) []string {
	rows, err := DbDomains.Query(table.READ)
	if err != nil {
		util.Fatalln(err)
	}
	s := []string{}
	for rows.Next() {
		var r string
		err = rows.Scan(&r)
		if err != nil {
			util.Fatalln(err)
		}
		s = append(s, r)
	}
	if err := rows.Close(); err != nil {
		util.Fatalln(err)
	}
	if err := rows.Err(); err != nil {
		util.Fatalln(err)
	}
	return s
}
func ReadQuery() []Query {
	rows, err := DbDomains.Query(table_query.READ)
	if err != nil {
		util.Fatalln(err)
	}

	queries := []Query{}
	for rows.Next() {
		q := Query{}
		err = rows.Scan(&q.Domain, &q.ClientIP, &q.Timestamp, &q.Diff, &q.OrganizerDomain, &q.Status, &q.IsBlocked)
		if err != nil {
			util.Fatalln(err)
		}
		queries = append(queries, q)
	}
	if err := rows.Close(); err != nil {
		util.Fatalln(err)
	}
	if err := rows.Err(); err != nil {
		util.Fatalln(err)
	}

	return queries
}

//TODO this function is called in go routine. meaning, var last_time is not accurately, thus organizer_domain does not make sense
//TODO employ better organizer: API domain is interfering
//Ideas
//- Special domain list
//- White listed domain -> organizer
//- Frequest domain that often requested first -> API???
func AddQuery(domain, clientip string, t time.Time, status string, isblocked string) {
	//although AddQuery is in qoroutine because we don't wait for it to finish, we lock it because sqlite only allows single writing operation.
	mu.Lock()
	defer mu.Unlock()
	var diff int64
	if !last_time.IsZero() {
		diff = t.Sub(last_time).Milliseconds()
		if diff > config.Cfg.Cluster.Diff {
			util.Println("\033[46mNew Organizer:\033[00m", domain)
			organizer_domain = domain
		}
	}
	last_time = t
	timestamp := t.Format(time.RFC3339Nano)
	//Because speed is important, don't ReadQuery() to update Queries
	Queries = append(Queries, Query{
		Domain:          domain,
		ClientIP:        clientip,
		Timestamp:       timestamp,
		Diff:            diff,
		OrganizerDomain: organizer_domain,
		Status:          status,
		IsBlocked:       isblocked,
	})
	_, err := query_adder.Exec(domain, clientip, timestamp, diff, organizer_domain, status, isblocked)
	if err != nil {
		util.Fatalln(err)
	}
}

func AddBlacklist(v string) []string {
	mu.Lock()
	_, err := DbDomains.Exec(table_blacklist.INSERT, v)
	if err != nil {
		util.Fatalln(err)
	}
	_, err = DbDomains.Exec(table_query.UPDATE, "yes", v)
	if err != nil {
		util.Fatalln(err)
	}
	mu.Unlock()

	Queries = ReadQuery()
	Blacklist = ReadSliceFrom(table_blacklist)

	s := []string{}
	s = append(s, Blacklist...)
	//#issue
	//s includes "" at the end of slice
	return s
}

func UpdateBlacklist(new string, old string) []string {
	mu.Lock()
	_, err := DbDomains.Exec(table_blacklist.UPDATE, new, old)
	if err != nil {
		util.Fatalln(err)
	}
	_, err = DbDomains.Exec(table_query.UPDATE, "no", old)
	if err != nil {
		util.Fatalln(err)
	}
	mu.Unlock()
	Queries = ReadQuery()
	Blacklist = ReadSliceFrom(table_blacklist)

	s := []string{}
	s = append(s, Blacklist...)
	return s
}

func DeleteBlacklist(v string) []string {
	mu.Lock()
	_, err := DbDomains.Exec(table_blacklist.DELETE, v)
	if err != nil {
		util.Fatalln(err)
	}
	_, err = DbDomains.Exec(table_query.UPDATE, "no", v)
	if err != nil {
		util.Fatalln(err)
	}
	mu.Unlock()
	Queries = ReadQuery()
	Blacklist = ReadSliceFrom(table_blacklist)

	s := []string{}
	s = append(s, Blacklist...)
	return s
}
