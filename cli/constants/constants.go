package constants

import "os"

func userHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	// fallback
	dir, _ := os.UserHomeDir()
	return dir
}
var HomeDir string = userHomeDir()

const (
	BaseDir                      = ".uplog"
	ConfigDir                    = "config"
	TmpDir                       = "tmp"
	SqliteFileName               = "db.sqlite"
	Domain                       = "https://logs.uplog.com"
	LogsDomain                   = Domain + "/session="
	BackendDomain                = "http://127.0.0.1:8000"
	BackendUploadEndpoint        = BackendDomain + "/session/upload"
	BackendUserCreateEndpoint	 = BackendDomain + "/user/create"
	BackendSessionCreateEndpoint = BackendDomain + "/session/create"
	ChunkSize                    = 32 * 1024 // 32 KB
	BatchLimit                   = 10
	PollIntervalLimit            = 100
)
