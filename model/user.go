package model

type User struct {
	ID           int64
	Name         string
	RealName     string
	Password     string
	Email        string
	Registration int64
	Anon         bool
	Admin        bool
}

// CreateUser creates record of a new user
func CreateUser(u *User) (err error) {
	db.NamedExec(`INSERT INTO user (name, password, registration, realname)
				VALUES (:name, :password, :registration, '')`, u)
	return nil
}

func GetUserByName(name string) (*User, error) {
	if len(name) == 0 {
		return nil, logger.Error("User doesn't exist")
	}

	u := &User{}
	err := db.QueryRowx(`SELECT * FROM user WHERE name = $1`, name).StructScan(u)
	if err != nil {
		return nil, logger.Error("User doesn't exist", "err", err)
	}
	return u, nil
}

func (u *User) IsAdmin() bool {
	if u.ID == 1 || u.Admin {
		return true
	}
	return false
}
