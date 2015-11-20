package model

import "errors"

type Revision struct {
	ID        int64
	PageTitle string
	TextID    int64
	Comment   string
	UserID    int
	UserText  string
	Minor     bool
	Deleted   bool
	Len       int
	ParentID  int
}

func (r Revision) Verify() error {
	if r.PageTitle == "" || r.UserText == "" || r.Len < 1 {
		return errors.New("Invalid revision")
	}
	return nil
}
