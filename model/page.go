package model

import (
	"strings"

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
func GetPageVeiwByID(revID string) *PageView {
	var p PageView
	db.QueryRowx(`SELECT page.namespace, page.title, page.nicetitle, text.text
				 FROM page JOIN revision on page.title = revision.pagetitle
				 JOIN text on revision.textid = text.id WHERE revision.id = $1`, revID).StructScan(&p)
	return &p
}

type CreatePageOptions struct {
	Title     string
	Namespace string
	Text      string
	Comment   string
	IsMinor   bool
}

// CreateOrUpdatePage updates or creates a new page in the wiki
func CreateOrUpdatePage(u *User, opts CreatePageOptions) (*Page, error) {
	var err error

	tx := db.MustBegin()
	t := createText(tx, opts.Text)

	rev, err := createRevision(tx, CreateRevOptions{
		Title:   opts.Title,
		Comment: opts.Comment,
		IsMinor: opts.IsMinor,
		Usr:     u,
		Txt:     t,
	})

	if err != nil {
		tx.Rollback()
		return nil, logger.Error("Insert revision failed", "err", err)
	}

	p := &Page{
		Title:      opts.Title,
		Namespace:  opts.Namespace,
		NiceTitle:  NiceTitle(opts.Title),
		Len:        len(opts.Text),
		RevisionID: rev.ID,
	}
	tx.Exec(`INSERT OR REPLACE INTO page (title, namespace, nicetitle, redirect, revisionid, len)
						VALUES ($1, $2, $3, $4, $5, $6)`, p.Title, p.Namespace, p.NiceTitle, p.Redirect, p.RevisionID, p.Len)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, logger.Error("Transaction failed", "err", err)
	}

	return p, nil
}

// NiceTitle converts a title into a nice title
func NiceTitle(title string) string {
	return strings.Replace(title, "_", " ", -1)
}

func GetPages() ([]*Page, error) {
	var pages []*Page
	err := db.Select(&pages, `SELECT * FROM page ORDER BY title DESC`)
	if err != nil {
		return pages, logger.Error("Unable to get pages", "err", err)
	}
	return pages, nil
}
