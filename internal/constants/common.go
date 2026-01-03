package constants

var (
	WS_TIMEOUT_EXPIRETIME          = 3
	VIDEO_CONTROL_INTERVAL         = 1
	REDIS_TIMEOUT                  = 200
	VIDEO_STATE_TIME               = 2 * 60 * 60
	FILE_SIZE_1G                   = 1024 * 1024 * 1024
	FILE_SIZE_500M                 = 1024 * 1024 * 500
	UPLOAD_CHUNK_SIZE              = 5_242_880
	UPLOAD_TEMP_DIR                = "/temp/upload"
	MERGE_TEMP_DIR                 = "/temp/merge"
	MINIO_UPLOAD_TIMEOUT           = 1000
	UPLOAD_CHUNK_IDX_TIMEOUT       = 2 * 24
	ROOM_JOIN_USER_TIMEINTERVAL    = 5 * 1000
	USER_CREATE_GROUP_TIMEINTERVAL = 60 * 1000
	ROOM_JOINCODE_EFFECTIVE_TIME   = 5 * 60 * 1000
	KAFKA_TIMEOUT                  = 500
	VIDEOS_TRANSCODE_TIMEOUT       = 60 * 1000
)
