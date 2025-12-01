package utils

import (
	"path/filepath"
	"os"
	"cli/constants"
)

func SetupDirectories(home string) {
	base := filepath.Join(home, constants.BaseDir)
	_ = os.MkdirAll(filepath.Join(base, constants.ConfigDir), 0o700)
	_ = os.MkdirAll(filepath.Join(base, constants.TmpDir), 0o700)
}
