package setup

import (
	"cli/api"
	"cli/db"
	"cli/utils"
	"fmt"
	"database/sql"

)

var (
	ApiKey    string
	UserId    string
	SessionId string
	DB_con *sql.DB
)

func Setup(setup bool)(string, string, string,*sql.DB){
		utils.SetupDirectories()
		DB_con = db.InitDB()
		var err error
		if setup{
		ApiKey,UserId,SessionId,err = api.SetupSession()
		if err!=nil{
			fmt.Println("Error on setting up session",err)
		}
		return ApiKey,UserId,SessionId,DB_con
	}
		return "","","",DB_con
}
