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

// View is a page with everything required for viewing it on the site
type View struct {
	NameSpace string
	Title     string
	NiceTitle string
	Text      string
	TimeStamp int64
	Deleted   bool
}

func (pv *View) Html() string {
	md := blackfriday.MarkdownCommon([]byte(pv.Text))
	return string(bluemonday.UGCPolicy().SanitizeBytes(md))
}

func (pv *View) Diff(pv2 *View) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(pv.Text, pv2.Text, false)
	dmp.DiffEditCost = 8
	diffs = dmp.DiffCleanupSemantic(diffs)
	return diffPretty(diffs)
}

func (pv *View) PrettyTime() string {
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
func GetPageView(namespace string, title string) *View {
	var p View
	db.Get(&p, `SELECT page.namespace, page.title, page.nicetitle, text.text,
				revision.timestamp, revision.deleted FROM page JOIN revision ON
				page.revisionid = revision.id JOIN text
				ON revision.textid = text.id WHERE title = $1
				AND namespace  = $2`, title, namespace)
	return &p
}

func GetPageVeiwByID(revID string) (*View, error) {
	var p View
	err := db.Get(&p, `SELECT page.namespace, page.title, page.nicetitle, text.text,
				revision.timestamp, revision.deleted
				FROM page JOIN revision on page.title = revision.pagetitle
				JOIN text on revision.textid = text.id WHERE revision.id = $1`, revID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetRandomPageViewTitle() string {
	title := []string{}
	db.Select(&title, `SELECT title FROM page JOIN revision on page.revisionid = revision.id
		WHERE namespace="" AND revision.deleted=0 ORDER BY title`)
	return title[rand.Intn(len(title))]
}

func GetPageViews() ([]*View, error) {
	var pages []*View
	err := db.Select(&pages, `SELECT page.namespace, page.title, page.nicetitle, text.text,
				revision.timestamp, revision.deleted FROM page JOIN revision ON
				page.revisionid = revision.id JOIN text
				ON revision.textid = text.id`)
	if err != nil {
		return pages, err
	}
	return pages, nil
}
