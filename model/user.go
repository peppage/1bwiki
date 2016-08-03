package model

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Name         string
	RealName     string
	Password     string
	Email        string
	Registration int64
	Anon         bool
	Admin        bool
	TimeZone     string
	DateFormat   string
}

var userFields = "name, registration, realname, timezone"

// CreateUser creates record of a new user
func CreateUser(u *User) (err error) {
	err = u.EncodePassword()
	if err != nil {
		return err
	}
	_, err = db.NamedExec(`INSERT INTO user (name, password, registration, realname, timezone)
							VALUES (:name, :password, :registration, '', :timezone)`, u)
	return err
}

func UpdateUserSettings(u *User) error {
	_, err := db.NamedExec(`UPDATE user SET timezone=:timezone, realname=:realname, dateformat=:dateformat WHERE id=:id`, u)
	return err
}

func UpdateUserPassword(u *User) error {
	_, err := db.NamedExec(`UPDATE user SET password=:password WHERE id=:id`, u)
	return err
}

func GetUserByName(name string) (*User, error) {
	if len(name) == 0 {
		return nil, errors.New("User doesn't exist")
	}

	u := &User{}
	err := db.QueryRowx(`SELECT * FROM user WHERE name = $1`, name).StructScan(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) IsAdmin() bool {
	return u.ID == 1 || u.Admin
}

func (u *User) EncodePassword() error {
	p, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.Password = string(p)
	return err
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsLoggedIn() bool {
	return u.Registration != 0
}

func GetUsers() ([]*User, error) {
	var users []*User
	err := db.Select(&users, `SELECT `+userFields+` FROM user ORDER BY name DESC`)
	if err != nil {
		return users, err
	}
	return users, nil

}
