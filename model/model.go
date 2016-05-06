package model

import (
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Connect("sqlite3", "./1bwiki.db")
	if err != nil {
		panic("failed to connect to db " + err.Error())
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
	tx.Exec(`CREATE TABLE IF NOT EXISTS user (id integer PRIMARY KEY, name text UNIQUE,
			realname text text default "", password text, registration int, email text default "",
			admin bool default false, UNIQUE(id, name))`)
	tx.Exec(`CREATE TABLE IF NOT EXISTS settings (name text PRIMARY KEY, value text)`)
	tx.Exec(`INSERT INTO settings (name, value) values ("anonediting", "true")`)
	tx.Exec(`INSERT INTO settings (name, value) values ("allowsignups", "true")`)

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		log.WithError(err).Error("Failed initializing db")
	}

	var texts int
	db.Get(&texts, `SELECT COUNT(*) as texts FROM text`)
	if texts == 0 {
		box, err := rice.FindBox("setup")
		if err != nil {
			log.WithError(err).Error("can't find setup rice box")
		}
		d, err := box.String("default.md")
		if err != nil {
			log.WithError(err).Error("default.md file error")
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
