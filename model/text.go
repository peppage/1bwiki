package model

import "errors"

type Text struct {
	ID   int64
	Text string
}

func (t Text) Verify() error {
	if t.Text == "" {
		return errors.New("Empty text")
	}
	return nil
}
