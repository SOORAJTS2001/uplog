package commands

import (
	"fmt"
	"os"
	"database/sql"
	"uplog/db"
	"uplog/utils"
)
func ListCmd(sql_object*sql.DB, home string) {
	sessions, err := db.ListSessions(sql_object)
	if err != nil {
		fmt.Printf("failed to list sessions: %v\n", err)
		os.Exit(1)
	}
	if len(sessions) == 0 {
		fmt.Println("no sessions found")
		return
	}
	for _, s := range sessions {
		ts := s.CreatedAt.Format("2006-01-02 15:04:05")
		uploaded := "no"
		if s.IsUploaded {
			uploaded = "yes"
		}
		fmt.Printf("%s | %s | %d bytes | %d lines | uploaded: %s | %s\n",
			ts, s.Tag, s.SizeBytes, s.LineCount, uploaded, utils.ConstructShareURL(s.SessionID))
	}
}
