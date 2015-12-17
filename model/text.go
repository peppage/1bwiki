package model

import (
	"github.com/jmoiron/sqlx"
)

type Text struct {
	ID   int64
	Text string
}

func createText(tx *sqlx.Tx, text string) *Text {
	result := tx.MustExec(`INSERT INTO text (text) VALUES ($1)`, text)
	lastID, _ := result.LastInsertId()
	t := &Text{Text: text, ID: lastID}
	return t
}

// CreateText creates a new text for a page
func CreateText(text string) (*Text, error) {
	if len(text) == 0 {
		return nil, logger.Error("Invalid Text")
	}
	tx := db.MustBegin()
	t := createText(tx, text)
	err := tx.Commit()
	if err != nil {
		return nil, logger.Error("Unable to commit text")
	}
	return t, nil
}
