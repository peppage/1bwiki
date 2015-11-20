package model

import "fmt"

type Page struct {
	Title      string
	Namespace  string
	NiceTitle  string
	Redirect   bool
	RevisionID int64
	Len        int
}

func GetPage(namespace string, title string) *Page {
	var p *Page
	db.Select(&p, `SELECT * FROM pages WHERE title = $1 AND namespace  = $2`, title, namespace)
	return p
}

func (p Page) SavePage(t Text, r Revision) {
	tx := db.MustBegin()
	result := tx.MustExec(`INSERT INTO text (text) VALUES ($1)`, t.Text)
	lastID, _ := result.LastInsertId()
	r.TextID = lastID
	result = tx.MustExec(`INSERT INTO revision (pagetitle, textid, comment, userid, usertext, minor,
						deleted, len, parentid, sha1) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		r.PageTitle, r.TextID, r.Comment, r.UserID, r.UserText, r.Minor, r.Deleted, r.Len, r.ParentID, r.Sha1)
	lastID, _ = result.LastInsertId()
	r.ID = lastID
	p.RevisionID = r.ID
	tx.MustExec(`INSERT INTO page (title, namespace, nicetitle, redirect, revisionid, len)
						VALUES ($1, $2, $3, $4, $5, $6)`, p.Title, p.Namespace, p.NiceTitle, p.Redirect, p.RevisionID, p.Len)
	err := tx.Commit()
	if err != nil {
		fmt.Println(err)
	}
}
