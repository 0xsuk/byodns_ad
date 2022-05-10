package model

import (
	"strings"
	"time"

	"github.com/0xsuk/byodns/util"
	"github.com/go-redis/redis"
)

const (
	v4Prefix = "ipv4:"
	v6Prefix = "ipv6:"
)

func SetDomainIPv4(domain string, ipv4 string, ttl time.Duration) {
	err := RedisInstance.Set(v4Prefix+domain, ipv4, ttl).Err()
	//basic nil
	if err != nil {
		util.Fatalln(err)
	}
}

func GetDomainIPv4(domain string) string {
	res, err := RedisInstance.Get(v4Prefix + domain).Result()
	//true when err exept res.Nil occured
	if err != nil && err != redis.Nil {
		util.Fatalln(err)
	}
	return strings.Replace(res, v4Prefix, "", 1)
}

func SetDomainIPv6(domain string, ipv6 string, ttl time.Duration) {
	err := RedisInstance.Set(v6Prefix+domain, ipv6, ttl).Err()
	if err != nil {
		util.Fatalln(err)
	}
}

func GetDomainIPv6(domain string) string {
	res, err := RedisInstance.Get(v6Prefix + domain).Result()
	//true when err exept res.Nil occured
	if err != nil && err != redis.Nil {
		util.Fatalln(err)
	}
	return strings.Replace(res, v6Prefix, "", 1)
}
