package model

type Page struct {
	Title      string
	Namespace  string
	NiceTitle  string
	Redirect   bool
	RevisionId int
	Len        int
}

func GetPage(ID string) *Page {
	var p *Page
	db.Select(&p, `SELECT * FROM pages WHERE ID = $1`, ID)
	return p
}
