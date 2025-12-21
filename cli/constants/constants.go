package constants

import ("os"
"path/filepath"
)

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
	base                         = ".uplog"
	config                       = "config"
	tmpDir                       = "tmp"
	credentialsFileName          = "credentials.json"
	SqliteFileName               = "db.sqlite"
	Domain                       = "http://localhost:8080/live-logs/"
	LogsDomain                   = Domain + "/session="
	BackendDomain                = "http://127.0.0.1:8000"
	BackendUploadEndpoint        = BackendDomain + "/session/upload"
	BackendUserCreateEndpoint	 = BackendDomain + "/user/create"
	BackendSessionCreateEndpoint = BackendDomain + "/session/create"
	ChunkSize                    = 32 * 1024 // 32 KB
	BatchLimit                   = 1
	PollIntervalLimit            = 100
)

var BaseDir = filepath.Join(HomeDir,base)
var ConfigDir = filepath.Join(BaseDir,config)
var CredentialsFile = filepath.Join(ConfigDir,credentialsFileName)
var TmpDir = filepath.Join(BaseDir,tmpDir)
