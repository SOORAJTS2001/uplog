package utils

import (
	"cli/constants"
	"cli/db"
	"cli/models"
	"database/sql"
	"fmt"
	"time"
	"os"
)

func SetupDirectories() {
	_ = os.MkdirAll(constants.ConfigDir, 0o700)
	_ = os.MkdirAll(constants.TmpDir, 0o700)
	_,_ = os.OpenFile(constants.CredentialsFile,os.O_CREATE|os.O_EXCL,0o700)
}

func Help(){
	fmt.Println("Visit https://www.uplog.live/docs for full documentation")
	fmt.Println()
	fmt.Println("Reserved commands")
	fmt.Println("Usage `uplog <command>`")
	fmt.Println("\t\t\t","list - ", "To list your entire session data")
	fmt.Println("\t\t\t","purge - ","To purge the entire session data including the temp log files")
	fmt.Println("\t\t\t","delete - ","To delete an existing session data")
	fmt.Println("Reserved flags")
	fmt.Println("Usage `uplog --flag`")
	fmt.Println("\t\t\t","batch - ","Batch size to upload for the session")
	fmt.Println("\t\t\t","poll - ","How often should it be polled for update")
	fmt.Println("\t\t\t","tag - ","Tag name for the session")
}

func NewSession(sessionId string,tag string,userId string) *models.Session {
	now := time.Now()
	var url string = constants.Domain+"subject."+userId+"."+sessionId
	return &models.Session{
		SessionId:  sessionId,
		CreatedAt:  now,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		LineCount:  0,
		SizeBytes:  0,
		IsUploaded: false,
		Mode:       "default",
		Tag:       	tag,
		Url:        url,
	}
}


func RelativeTime(t time.Time) string {
	now := time.Now()
	d := t.Sub(now)

	abs := d
	if abs < 0 {
		abs = -abs
	}

	switch {
	case abs < time.Minute:
		if d < 0 {
			return "just now"
		}
		return "right now"

	case abs < time.Hour:
		n := int(abs.Minutes())
		if d < 0 {
			return fmt.Sprintf("%d minutes ago", n)
		}
		return fmt.Sprintf("%d minutes", n)

	case abs < 24*time.Hour:
		n := int(abs.Hours())
		if d < 0 {
			return fmt.Sprintf("%d hours ago", n)
		}
		return fmt.Sprintf("%d hours", n)

	case abs < 7*24*time.Hour:
		n := int(abs.Hours() / 24)
		if d < 0 {
			return fmt.Sprintf("%d days ago", n)
		}
		return fmt.Sprintf("%d days", n)

	default:
		n := int(abs.Hours() / (24 * 7))
		if d < 0 {
			return fmt.Sprintf("%d weeks ago", n)
		}
		return fmt.Sprintf("%d weeks", n)
	}
}


func ListLogs(db_con *sql.DB){
	sessions,err:=db.ListSessions(db_con)
	if err!=nil{
		fmt.Println("Listing sessions failed due to",err)
	}
	fmt.Println("SessionId","| Tag","| Created at","| Expires at","| Size","| Lines","| Url")
	for _, s := range sessions {
	fmt.Println(s.SessionId,s.Tag,RelativeTime(s.CreatedAt),RelativeTime(s.ExpiresAt),s.SizeBytes,s.LineCount,s.Url)
}
}
func PurgeLogs(db_con *sql.DB){
	err:=db.DeleteAllSessions(db_con)
	if err!=nil{
		fmt.Println("Purging sessions failed due to ",err)
	}
	fmt.Println("Successfully purged all sessions")
}
func DeleteLog(sessionId string,db_con *sql.DB){
	err:=db.DeleteSessionById(sessionId,db_con)
	if err!=nil{
		fmt.Println("Deleting log failed due to",err)
	}
	fmt.Println("Successfully deleted session")
}
func UploadLog(args []string){

}
