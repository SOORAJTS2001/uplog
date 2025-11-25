package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"database/sql"
	"uplog/config"
)
func DeleteCmd(sql_object *sql.DB,home, sessionID string) {
	if sessionID == "--all" {
		sql_object.Exec(`DELETE FROM sessions`)
		files, _ := filepath.Glob(filepath.Join(home, config.BaseDir, config.TmpDir, "*.log"))
		for _, f := range files {
			_ = os.Remove(f)
		}
		fmt.Println("deleted all uplog sessions")
		return
	}
	// delete from sqlite
	_, err := sql_object.Exec(`DELETE FROM sessions WHERE session_id = ?`, sessionID)
	if err != nil {
		fmt.Printf("failed to delete session: %v\n", err)
		os.Exit(1)
	}

	// 2. remove temp log file if exists
	tmpPath := filepath.Join(home, config.BaseDir, config.TmpDir, sessionID+".log")
	_ = os.Remove(tmpPath) // ignore error if not exists

	fmt.Printf("deleted session %s\n", sessionID)
}
