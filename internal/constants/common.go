package constants

var (
	WS_TIMEOUT_EXPIRETIME          = 3
	VIDEO_CONTROL_INTERVAL         = 500
	REDIS_TIMEOUT                  = 200
	UPLOAD_CHUNK_SIZE              = 5_242_880
	UPLOAD_TEMP_DIR                = "/temp/upload"
	MERGE_TEMP_DIR                 = "/temp/merge"
	UPLOAD_CHUNK_IDX_TIMEOUT       = 2 * 24
	ROOM_JOIN_USER_TIMEINTERVAL    = 5 * 1000
	USER_CREATE_GROUP_TIMEINTERVAL = 60 * 1000
	ROOM_JOINCODE_EFFECTIVE_TIME   = 5 * 60 * 1000
	VIDEOS_TRANSCODE_TIMEOUT       = 60 * 1000
	VIDEO_STATE_TTL                = 10
	VIDEO_URL_TTL                  = 30
)
