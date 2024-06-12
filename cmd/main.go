package main

import (
	"os"

	"backupdb/handlers"
)

func main() {
	handlers.InitLogger()
	handlers.LoadConfig()
	handlers.BackupDB()
	os.Exit(0)
}
