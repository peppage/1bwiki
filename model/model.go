package model

import (
	"crypto/rand"

	"github.com/GeertJohan/go.rice"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mgutz/logxi/v1"
)

var db *sqlx.DB
var logger log.Logger

func init() {
	logger = log.New("model")
	var err error
	db, err = sqlx.Connect("sqlite3", "./1bwiki.db")
	if err != nil {
		logger.Error("connecting to db", "err", err)
	}
}

func SetupDb() {
	tx := db.MustBegin()
	tx.Exec(`create table if not exists text (id integer primary KEY, text blob)`)
	tx.Exec(`create table if not exists revision (id integer primary key,
			pagetitle text, textid integer, comment text, userid int,
			usertext text, minor integer, deleted integer, len integer,
			parentid integer, timestamp integer, lendiff integer)`)
	tx.Exec(`create table if not exists page (title text,
			namespace text, nicetitle text, redirect integer, revisionid integer,
			len integer, PRIMARY KEY(title, namespace))`)
	tx.Exec(`CREATE TABLE IF NOT EXISTS user (id integer PRIMARY KEY, name text,
			realname text text default "", password text, registration int, email text default "",
			admin bool default false, UNIQUE(id, name))`)
	tx.Exec(`CREATE TABLE IF NOT EXISTS settings (name text PRIMARY KEY, value text)`)
	tx.Exec(`INSERT INTO settings (name, value) values ("anonediting", "true")`)
	tx.Exec(`INSERT INTO settings (name, value) values ("allowsignups", "true")`)
	tx.Exec(`INSERT INTO settings (name, value) values ("sessionsecret", $1)`, randString(20))

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("db setup", "err", err)
	}

	var texts int
	db.Get(&texts, `SELECT COUNT(*) as texts FROM text`)
	if texts == 0 {
		box, err := rice.FindBox("setup")
		if err != nil {
			logger.Error("can't find setup rice box", "err", err)
		}
		d, err := box.String("default.md")
		if err != nil {
			logger.Error("default.md file error", "err", err)
		}
		u := &User{
			ID:   0,
			Name: "Admin",
		}
		CreateOrUpdatePage(u, CreatePageOptions{
			Title:     "Main_Page",
			Namespace: "",
			Text:      d,
			Comment:   "",
			IsMinor:   false,
		})
	}
}

func randString(size int) string {
	const alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, size)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = alpha[v%byte(len(alpha))]
	}
	return string(bytes)
}
