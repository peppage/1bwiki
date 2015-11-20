package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mgutz/logxi/v1"
)

var db *sqlx.DB
var logger log.Logger

func init() {
	var err error
	db, err = sqlx.Connect("sqlite3", "./1bwiki.db")
	if err != nil {
		panic(err)
	}
	logger = log.New("model")
}
