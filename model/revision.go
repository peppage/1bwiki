package model

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
		return logger.Error("Invalid revision")
	}
	return nil
}

func GetRevisions() ([]*Revision, error) {
	var revs []*Revision
	err := db.Select(&revs, `SELECT * FROM revision ORDER BY id DESC`)
	if err != nil {
		return revs, logger.Error("Unable to select all revisions", "err", err)
	}
	return revs, nil
}
