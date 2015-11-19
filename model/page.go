package model

type Page struct {
	ID int
}

func GetPage(ID string) *Page {
	var p *Page
	db.Select(&p, `SELECT * FROM pages WHERE ID = $1`, ID)
	return p
}
