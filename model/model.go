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
	// Convert to transaction
	db.Exec(`create table if not exists text (id integer primary KEY, text blob)`)
	db.Exec(`create table if not exists revision (id integer primary key,
			pagetitle text, textid integer, comment text, userid int,
			usertext text, minor integer, deleted integer, len integer,
			parentid integer, timestamp integer, lendiff integer)`)
	db.Exec(`create table if not exists page (title text,
			namespace text, nicetitle text, redirect integer, revisionid integer,
			len integer, PRIMARY KEY(title, namespace))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS user (id integer PRIMARY KEY, name text,
			realname text text default "", password text, registration int, email text,
			admin bool, UNIQUE(id, name))`)
	db.Exec(`CREATE TABLE IF NOT EXISTS settings (name text PRIMARY KEY, value text)`)
	db.Exec(`INSERT INTO settings (name, value) values ("anonediting", "true")`)
	db.Exec(`INSERT INTO settings (name, value) values ("allowsignups", "true")`)
	logger = log.New("model")
}
