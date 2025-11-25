package config


import (
	"time"

)
const (
	BaseDir                      = ".uplog"
	ConfigDir                    = "config"
	TmpDir                       = "tmp"
	SqliteFileName               = "db.sqlite"
	ConfigFileName				 = "config.json"
	Domain                       = "https://logs.uplog.com"
	LogsDomain                   = Domain + "/session="
	BackendDomain                = "http://127.0.0.1:8000"
	BackendUploadEndpoint        = BackendDomain + "/session/upload"
	BackendSessionCreateEndpoint = BackendDomain + "/session/create"
	ChunkSize                    = 1024 * 1024 // 32 KB
	BatchLimit                   = 200
	PollIntervalLimit            = 200 * time.Millisecond
)
