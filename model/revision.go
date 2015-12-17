package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Revision struct {
	ID        int64
	PageTitle string
	TextID    int64
	Comment   string
	UserID    int64
	UserText  string
	Minor     bool
	Deleted   bool
	Len       int
	ParentID  int64
	TimeStamp int64
	LenDiff   int
}

type CreateRevOptions struct {
	Title   string
	Comment string
	IsMinor bool
	Txt     *Text
	Usr     *User
}

func createRevision(tx *sqlx.Tx, opts CreateRevOptions) (*Revision, error) {
	r := convertOptions(opts)
	oldRev := &Revision{}
	tx.QueryRowx(`SELECT * FROM revision WHERE pagetitle = $1 ORDER BY id desc LIMIT 1`, r.PageTitle).StructScan(oldRev)
	r.ParentID = oldRev.ID
	r.LenDiff = r.Len - oldRev.Len
	result, err := tx.Exec(`INSERT INTO revision (pagetitle, textid, comment, userid, usertext, minor,
		deleted, len, parentid, timestamp, lendiff) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		r.PageTitle, r.TextID, r.Comment, r.UserID, r.UserText, r.Minor, r.Deleted,
		r.Len, r.ParentID, r.TimeStamp, r.LenDiff)
	if err != nil {
		return nil, err
	}
	lastID, _ := result.LastInsertId()
	r.ID = lastID
	return r, nil
}

func CreateRevision(opts CreateRevOptions) (*Revision, error) {
	if opts.Title == "" {
		return nil, logger.Error("Invalid title")
	}
	tx := db.MustBegin()
	rev, err := createRevision(tx, opts)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return rev, tx.Commit()
}

func convertOptions(opts CreateRevOptions) *Revision {
	rev := &Revision{
		PageTitle: opts.Title,
		Comment:   opts.Comment,
		TextID:    opts.Txt.ID,
		UserID:    opts.Usr.ID,
		UserText:  opts.Usr.Name,
		Minor:     opts.IsMinor,
		Len:       len(opts.Txt.Text),
		TimeStamp: time.Now().Unix(),
	}
	return rev
}

func GetRevisions() ([]*Revision, error) {
	var revs []*Revision
	err := db.Select(&revs, `SELECT * FROM revision ORDER BY id DESC`)
	if err != nil {
		return revs, logger.Error("Unable to select all revisions", "err", err)
	}
	return revs, nil
}

func GetPageRevisions(title string) ([]*Revision, error) {
	var revs []*Revision
	err := db.Select(&revs, `SELECT * FROM revision WHERE pagetitle=$1 ORDER BY id DESC`, title)
	if err != nil {
		return revs, logger.Error("Unable to get revisions for page", "page", title, "err", err)
	}
	return revs, nil
}
