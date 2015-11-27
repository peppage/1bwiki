package model

type User struct {
	ID           int64
	Name         string
	RealName     string
	Password     string
	Email        string
	Registration int64
}

func (u *User) Create() {
	db.NamedExec(`INSERT INTO user (name, password, registration)
				VALUES (:name, :password, :registration)`, u)
}
