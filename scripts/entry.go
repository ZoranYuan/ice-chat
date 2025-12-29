package scripts

import "github.com/redis/go-redis/v9"

var (
	UpdateLatestTimestamp = redis.NewScript(`
		local old = redis.call("GET", KEYS[1])
		if not old or tonumber(ARGV[2]) > tonumber(old)
		then
			redis.call("SET", KEYS[1], ARGV[2], "EX", ARGV[1])
			redis.call("SET", KEYS[2], ARGV[3], "EX", ARGV[1])
			return 1
		end 
			return 0
	`)
)
