package model

type Revision struct {
	ID        int
	PageTitle string
	TextID    int
	Comment   string
	UserId    int
	UserText  string
	Minor     bool
	Deleted   bool
	Len       int
	ParentId  int
	Sha1      string
}
