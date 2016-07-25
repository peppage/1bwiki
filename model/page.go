package model

import (
	"errors"
	"strings"
)

type Page struct {
	Title      string
	Namespace  string
	NiceTitle  string
	Redirect   bool
	RevisionID int64
	Len        int
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
		return nil, err
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
		return nil, err
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
		return pages, err
	}
	return pages, nil
}

func DeletePage(u *User, title string) error {
	tx := db.MustBegin()
	t := createText(tx, "")
	var count int
	tx.Get(&count, `SELECT COUNT(*) FROM page WHERE title=$1`, title)
	if count == 1 {
		rev, err := createRevision(tx, CreateRevOptions{
			Title:   title,
			Comment: "Page Deleted",
			Deleted: true,
			Usr:     u,
			Txt:     t,
		})

		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Exec(`UPDATE page SET revisionid=$1 WHERE title=$2`, rev.ID, title)

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}

	tx.Rollback()
	return errors.New("Page doesn't exist!")
}
