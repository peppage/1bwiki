package model

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type Page struct {
	Title      string
	Namespace  string
	NiceTitle  string
	Redirect   bool
	RevisionID int64
	Len        int
}

type PageView struct {
	NameSpace string
	Title     string
	NiceTitle string
	Text      string
}

func (pv *PageView) Html() string {
	md := blackfriday.MarkdownCommon([]byte(pv.Text))
	return string(bluemonday.UGCPolicy().SanitizeBytes(md))
}

// Need error handling here
func GetPageView(namespace string, title string) *PageView {
	var p PageView
	db.QueryRowx(`select page.namespace, page.title, page.nicetitle, text.text
				FROM page JOIN revision ON
				page.revisionid = revision.id JOIN text
				ON revision.textid = text.id WHERE title = $1
				AND namespace  = $2`, title, namespace).StructScan(&p)
	return &p
}

// Need error handling here
func GetOldPageView(revID string) *PageView {
	var p PageView
	db.QueryRowx(`SELECT page.namespace, page.title, page.nicetitle, text.text
				 FROM page JOIN revision on page.title = revision.pagetitle
				 JOIN text on revision.textid = text.id WHERE revision.id = $1`, revID).StructScan(&p)
	return &p
}

func (p Page) SavePage(t Text, r Revision) error {
	var err error
	tx := db.MustBegin()
	result := tx.MustExec(`INSERT INTO text (text) VALUES ($1)`, t.Text)
	lastID, _ := result.LastInsertId()
	r.TextID = lastID
	result, err = tx.Exec(`INSERT INTO revision (pagetitle, textid, comment, userid, usertext, minor,
						deleted, len, parentid, timestamp, lendiff) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		r.PageTitle, r.TextID, r.Comment, r.UserID, r.UserText, r.Minor, r.Deleted, r.Len, r.ParentID, r.TimeStamp, r.LenDiff)
	if err != nil {
		tx.Rollback()
		return logger.Error("Insert revision failed", "err", err)
	}
	lastID, _ = result.LastInsertId()
	r.ID = lastID
	p.RevisionID = r.ID
	tx.MustExec(`INSERT OR REPLACE INTO page (title, namespace, nicetitle, redirect, revisionid, len)
						VALUES ($1, $2, $3, $4, $5, $6)`, p.Title, p.Namespace, p.NiceTitle, p.Redirect, p.RevisionID, p.Len)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return logger.Error("Transaction failed", "err", err)
	}
	return nil
}
