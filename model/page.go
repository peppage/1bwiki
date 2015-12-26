package model

import (
	"bytes"
	"html"
	"math/rand"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/sergi/go-diff/diffmatchpatch"
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
	TimeStamp int64
}

func (pv *PageView) Html() string {
	md := blackfriday.MarkdownCommon([]byte(pv.Text))
	return string(bluemonday.UGCPolicy().SanitizeBytes(md))
}

func (pv *PageView) Diff(pv2 *PageView) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(pv.Text, pv2.Text, false)
	dmp.DiffEditCost = 8
	diffs = dmp.DiffCleanupSemantic(diffs)
	return diffPretty(diffs)
}

func (pv *PageView) PrettyTime() string {
	t := time.Unix(pv.TimeStamp, 0).UTC()
	return t.Format("15:04, 2 Jan 2006")
}

func diffPretty(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := strings.Replace(html.EscapeString(diff.Text), "\n", "<br>", -1)
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString("<ins>")
			buff.WriteString(text)
			buff.WriteString("</ins>")
		case diffmatchpatch.DiffDelete:
			buff.WriteString("<del>")
			buff.WriteString(text)
			buff.WriteString("</del>")
		case diffmatchpatch.DiffEqual:
			buff.WriteString(text)
		}
	}
	return buff.String()
}

// GetPageView gets all information to show a page to a user
func GetPageView(namespace string, title string) *PageView {
	var p PageView
	db.QueryRowx(`select page.namespace, page.title, page.nicetitle, text.text,
				revision.timestamp FROM page JOIN revision ON
				page.revisionid = revision.id JOIN text
				ON revision.textid = text.id WHERE title = $1
				AND namespace  = $2`, title, namespace).StructScan(&p)
	return &p
}

func GetPageVeiwByID(revID string) (*PageView, error) {
	var p PageView
	err := db.QueryRowx(`SELECT page.namespace, page.title, page.nicetitle, text.text,
				revision.timestamp
				FROM page JOIN revision on page.title = revision.pagetitle
				JOIN text on revision.textid = text.id WHERE revision.id = $1`, revID).StructScan(&p)
	if err != nil {
		return nil, logger.Error("unable to get page view by ID", "err", err)
	}
	return &p, nil
}

func GetRandomPageViewTitle() string {
	title := []string{}
	db.Select(&title, `SELECT title FROM page WHERE namespace="" ORDER BY title`)
	return title[rand.Intn(len(title))]
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

func DeletePage(u *User, title string) error {
	tx := db.MustBegin()
	t := Text{}
	rev, err := createRevision(tx, CreateRevOptions{
		Title:   title,
		Comment: "Page Deleted",
		Deleted: true,
		Usr:     u,
		Txt:     &t,
	})

	if err != nil {
		tx.Rollback()
		return logger.Error("Insert revision failed", "err", err)
	}

	tx.Exec(`UPDATE page SET revisionid=$1 WHERE title=$2`, rev.ID, title)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return logger.Error("Transaction failed", "err", err)
	}
	return nil
}
