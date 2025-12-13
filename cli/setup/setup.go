package setup

import (
	"cli/api"
	"cli/db"
	"cli/utils"
	"fmt"
)
func Setup()(string, string, string){
		utils.SetupDirectories()
		db.InitDB()
		var apiKey,userId,sessionId, err = api.SetupSession()
		if err!=nil{
			fmt.Println("Error on settin up session",err)
		}
		return apiKey,userId,sessionId
}

var ApiKey,UserId,SessionId = Setup()
