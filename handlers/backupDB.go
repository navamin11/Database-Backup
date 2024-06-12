package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/joho/sqltocsv"
)

var ctx = context.Background()

func BackupDB() {
	var (
		wg      sync.WaitGroup
		sqlStmt string
	)

	mainDir, err := CreateDirectory()
	if err != nil {
		Logger.Sugar().Fatalf("%s", err)
	}

	for _, config := range DBconfig.Database {

		Logger.Sugar().Infof("%s", "-----------------------------------------------")
		Logger.Sugar().Infof("%s to %s", "Database Connection", config.DbName)
		Logger.Sugar().Infof("%s", "-----------------------------------------------")

		switch config.Driver {
		case "sqlserver":

			db, err := ConnectMSSQL(config.Driver, config.Username, config.Password, config.Host, config.Port, config.DbName)
			if err != nil {
				Logger.Sugar().Panicf("[%s] %s - %s", "X", config.DbName, err)
			}
			defer db.Close()

			subDir, err := CreateSubDirectory(*mainDir, config.Host, config.DbName)
			if err != nil {
				Logger.Sugar().Panicf("%s", err)
			}

			for _, table := range config.Tables {
				wg.Add(1)

				if len(strings.TrimSpace(table.Where)) != 0 {
					sqlStmt = fmt.Sprintf(`SELECT %s FROM %s.%s.%s WHERE %s`, table.Select, config.DbName, config.Schema, table.TbName, table.Where)
				} else {
					sqlStmt = fmt.Sprintf(`SELECT %s FROM %s.%s.%s`, table.Select, config.DbName, config.Schema, table.TbName)
				}

				go func(sqlStmt, TbName string, subdir *string) {
					defer wg.Done()

					defer func() {
						if r := recover(); r != nil {
							return
						}
					}()

					rows, err := db.Query(sqlStmt)
					if err != nil {
						Logger.Sugar().Errorf("[%s] %s - %s", "X", sqlStmt, err)
					}
					defer rows.Close()

					if sqltocsv.WriteFile(*subdir+"/"+TbName+".csv", rows); err != nil {
						Logger.Sugar().Errorf("[%s] %s - %s", "X", sqlStmt, err)
					} else {
						Logger.Sugar().Infof("[%s] %s", "/", sqlStmt)
					}
				}(sqlStmt, table.TbName, subDir)
			}
			wg.Wait()

		case "postgres":

			db, err := ConnectPostgreSQL(config.Driver, config.Username, config.Password, config.Host, config.Port, config.DbName)
			if err != nil {
				Logger.Sugar().Panicf("[%s] %s - %s", "X", config.DbName, err)
			}
			defer db.Close()

			subDir, err := CreateSubDirectory(*mainDir, config.Host, config.DbName)
			if err != nil {
				Logger.Sugar().Panicf("%s", err)
			}

			for _, table := range config.Tables {

				wg.Add(1)
				if len(strings.TrimSpace(table.Where)) != 0 {
					sqlStmt = fmt.Sprintf(`SELECT %s FROM %s.%s WHERE %s`, table.Select, config.Schema, table.TbName, table.Where)
				} else {
					sqlStmt = fmt.Sprintf(`SELECT %s FROM %s.%s`, table.Select, config.Schema, table.TbName)
				}

				go func(sqlStmt, TbName string, subdir *string) {
					defer wg.Done()

					defer func() {
						if r := recover(); r != nil {
							return
						}
					}()

					rows, err := db.Query(sqlStmt)
					if err != nil {
						Logger.Sugar().Errorf("[%s] %s - %s", "X", sqlStmt, err)
					}
					defer rows.Close()

					if sqltocsv.WriteFile(*subdir+"/"+TbName+".csv", rows); err != nil {
						Logger.Sugar().Errorf("[%s] %s - %s", "X", sqlStmt, err)
					} else {
						Logger.Sugar().Infof("[%s] %s", "/", sqlStmt)
					}
				}(sqlStmt, table.TbName, subDir)
			}
			wg.Wait()

		case "mysql":

			db, err := ConnectMySQL(config.Driver, config.Username, config.Password, config.Host, config.Port, config.DbName)
			if err != nil {
				Logger.Sugar().Panicf("[%s] %s - %s", "X", config.DbName, err)
			}
			defer db.Close()

			subDir, err := CreateSubDirectory(*mainDir, config.Host, config.DbName)
			if err != nil {
				Logger.Sugar().Panicf("%s", err)
			}

			for _, table := range config.Tables {

				wg.Add(1)
				if len(strings.TrimSpace(table.Where)) != 0 {
					sqlStmt = fmt.Sprintf(`SELECT %s FROM %s.%s WHERE %s`, table.Select, config.Schema, table.TbName, table.Where)
				} else {
					sqlStmt = fmt.Sprintf(`SELECT %s FROM %s.%s`, table.Select, config.Schema, table.TbName)
				}

				go func(sqlStmt, TbName string, subdir *string) {
					defer wg.Done()

					defer func() {
						if r := recover(); r != nil {
							return
						}
					}()

					rows, err := db.Query(sqlStmt)
					if err != nil {
						Logger.Sugar().Errorf("[%s] %s - %s", "X", sqlStmt, err)
					}
					defer rows.Close()

					if sqltocsv.WriteFile(*subdir+"/"+TbName+".csv", rows); err != nil {
						Logger.Sugar().Errorf("[%s] %s - %s", "X", sqlStmt, err)
					} else {
						Logger.Sugar().Infof("[%s] %s", "/", sqlStmt)
					}
				}(sqlStmt, table.TbName, subDir)
			}
			wg.Wait()

		default:
			Logger.Sugar().Warnf("[%s] Driver %s not supported.", "X", config.DbName)
		}
	}

	Logger.Sugar().Infof("%s", "-----------------------------------------------")
	Logger.Sugar().Infof("%s", "Generate Zip Archive")
	Logger.Sugar().Infof("%s", "-----------------------------------------------")

	if err := ZipSource(*mainDir, *mainDir); err != nil {
		Logger.Sugar().Errorf("[%s] Can't archive file %s.zip - %s", "X", *mainDir, err)
	} else {
		if err := DeleteDirectory(*mainDir); err != nil {
			Logger.Sugar().Errorf("[%s] Can't delete directory %s - %s", "X", *mainDir, err)
		} else {
			Logger.Sugar().Infof("[%s] Archive file %s.zip successfully", "/", *mainDir)
		}
	}
}
