package model

type Page struct {
	Title      string
	Namespace  string
	NiceTitle  string
	Redirect   bool
	RevisionId int
	Len        int
}

func GetPage(namespace string, title string) *Page {
	var p *Page
	db.Select(&p, `SELECT * FROM pages WHERE title = $1 and namespace  = $2`, title, namespace)
	return p
}
