package utils

import (
	"cli/constants"
	"fmt"
	"os"
)

func SetupDirectories() {
	fmt.Println(constants.CredentialsFile)
	_ = os.MkdirAll(constants.ConfigDir, 0o700)
	_ = os.MkdirAll(constants.TmpDir, 0o700)
	_,_ = os.OpenFile(constants.CredentialsFile,os.O_CREATE|os.O_EXCL,0o700)
}
