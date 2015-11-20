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
