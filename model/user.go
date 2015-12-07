package model

import "golang.org/x/crypto/bcrypt"

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

func (u *User) Verify() error {
	var dbUser User
	err := db.QueryRowx(`SELECT * FROM user WHERE name = $1`, u.Name).StructScan(&dbUser)
	if err != nil {
		return logger.Error("verify user db error", "err", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password))
	if err != nil {
		return logger.Error("password doesn't match", "err", err)
	}
	u = &dbUser
	return nil
}
