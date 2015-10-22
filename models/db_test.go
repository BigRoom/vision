package models

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Parse()

	Init(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_SERVICE"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db := *DB

	dbName := fmt.Sprintf("testing%v", time.Now().Unix())
	DB.DB.MustExec("CREATE DATABASE " + dbName + ";")

	Init(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_SERVICE"), os.Getenv("DB_PORT"), dbName)

	r := m.Run()

	DB.DB.Close()

	db.DB.MustExec("DROP DATABASE " + dbName + ";")

	os.Exit(r)
}
