package model

import (
	"errors"
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
	Deleted bool
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
		return nil, errors.New("Invalid Title")
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
		Deleted:   opts.Deleted,
		Minor:     opts.IsMinor,
		Len:       len(opts.Txt.Text),
		TimeStamp: time.Now().Unix(),
	}
	return rev
}

func GetRevisions(limit int) ([]*Revision, error) {
	var revs []*Revision
	err := db.Select(&revs, `SELECT * FROM revision ORDER BY id DESC LIMIT $1`, limit)
	if err != nil {
		return revs, err
	}
	return revs, nil
}

func GetLatestRevision(title string) (*Revision, error) {
	rev := Revision{}
	err := db.Get(&rev, `SELECT * FROM revision WHERE pagetitle = $1 ORDER BY id DESC`, title)
	if err != nil {
		return nil, err
	}
	return &rev, nil
}

func GetPageRevisions(title string, page int, limit int) ([]*Revision, error) {
	if limit == 0 {
		limit = 50
	}

	offset := 0
	if page != 0 {
		offset = page * limit
	}

	var revs []*Revision
	err := db.Select(&revs, `SELECT * FROM revision WHERE pagetitle=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`, title, limit, offset)
	if err != nil {
		return revs, err
	}
	return revs, nil
}

func GetAmountOfRevisionsForPage(title string) int {
	count := 0
	db.Get(&count, `SELECT COUNT(*) FROM revision WHERE pagetitle=$1`, title)
	return count
}

func (r *Revision) PrettyTime(timeZone, format string) string {
	l, _ := time.LoadLocation(timeZone)
	t := time.Unix(r.TimeStamp, 0).In(l)
	return t.Format(format)
}
