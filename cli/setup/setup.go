package setup

import (
	"cli/db"
	"cli/constants"
	"cli/utils"
)
func Setup(){
		db.InitDB(constants.HomeDir)
		utils.SetupDirectories(constants.HomeDir)
}
