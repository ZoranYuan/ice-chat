package scripts

import "github.com/redis/go-redis/v9"

var (
	UpdateLatestTimestamp = redis.NewScript(`
	if tonumber(ARGV[2]) > tonumber(ARGV[1])
		then
			redis.call("SET", KEYS[1], ARGV[3])
			return 1
		end 
			return 0
	`)

	UpdateUploadChunkIdx = redis.NewScript(`
		local current = redis.call("GET", KEYS[1])
		if not current
		then
			redis.call("SET", KEYS[1], 0, "EX", ARGV[1])
			return 0
		end
			return current
	`)

	UpdateVideoState = redis.NewScript(`
		local oldState = redis.call("GET", KEYS[1])
		if not oldState then
			return nil
		end
		
		local data = cjson.decode(oldState)
		local now = tonumber(ARGV[1])
		local action = ARGV[2]
		local time = tonumber(ARGV[3])
		local speed = ARGV[4]
		local timeStamp = tonumber(ARGV[5])


		if timeStamp <= data.updateAt then
			return nil
		end

		if data.isPlaying then 
			local delta = now - data.updateAt
			data.currentTime = data.currentTime + delta * data.speed
		end

		if action == 'play' then
			data.isPlaying = true
		elseif action == 'pause' then
			data.isPaused = false
		elseif action == 'seek' then
			data.currentTime = time
		elseif action == 'change_rate' then
			data.speed = speed
		end

		data.updateAt = now
		redis.call("SET", KEYS[1], cjson.encode(data), "EX", ARGV[6])
		
		return cjson.encode(data)
	`)
)
