package model

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

type dbrows struct {
	Cid     int            `db:"cid"`
	Name    string         `db:"name"`
	Type    string         `db:"type"`
	NotNull bool           `db:"notnull"`
	Default sql.NullString `db:"dflt_value"`
	Pk      int            `db:"pk"`
}

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
			realname text default "", password text, registration int, email text default "",
			admin bool default false, timezone text default "UTC", dateformat text default "15:04, 2 January 2006"
			UNIQUE(id, name))`)
	tx.Exec(`CREATE TABLE IF NOT EXISTS settings (name text PRIMARY KEY, value text)`)
	tx.Exec(`INSERT INTO settings (name, value) values ("anonediting", "true")`)
	tx.Exec(`INSERT INTO settings (name, value) values ("allowsignups", "true")`)

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		log.WithError(err).Error("Failed initializing db")
	}

	// Upgrade user's table to have timezone (from beta3)
	rows := []dbrows{}
	hasTimeZone := false
	err = db.Select(&rows, `PRAGMA table_info("user")`)
	if err != nil {
		log.WithError(err).Error("failed getting users table columns")
	}
	for _, r := range rows {
		if r.Name == "timezone" {
			hasTimeZone = true
		}
	}
	if !hasTimeZone {
		tx := db.MustBegin()
		tx.Exec(`create temporary table temp(id integer PRIMARY KEY, name text UNIQUE,
				realname text default "", password text, registration int, email text default "",
				admin bool default false, UNIQUE(id, name))`)
		tx.Exec(`insert into temp select id, name, realname, password, registration, email, admin from user`)
		tx.Exec(`drop table user`)
		tx.Exec(`CREATE TABLE IF NOT EXISTS user (id integer PRIMARY KEY, name text UNIQUE,
				realname text default "", timezone text default "UTC", dateformat text default "15:04, 2 January 2006", password text, registration int, email text default "",
				admin bool default false, UNIQUE(id, name))`)
		tx.Exec(`insert into user select id, name, realname, "UTC", "15:04, 2 January 2006", password, registration, email, admin from temp`)
		tx.Exec(`drop table temp`)
		tx.Exec(`UPDATE page SET namespace = "page"`)
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			log.WithError(err).Error("Failed upgrading user table")
		}
	}

	var texts int
	db.Get(&texts, `SELECT COUNT(*) as texts FROM text`)
	if texts == 0 {
		data, err := Asset("model/setup/default.md")
		if err != nil {
			log.WithError(err).Error("default.md file error")
		}
		u := &User{
			ID:   0,
			Name: "Admin",
		}
		CreateOrUpdatePage(u, CreatePageOptions{
			Title:     "Main_Page",
			Namespace: NameSpace[WikiPage],
			Text:      string(data),
			Comment:   "",
			IsMinor:   false,
		})
	}
}
