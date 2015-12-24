package model

import (
	"strconv"
)

type setting struct {
	Name  string
	Value string
}

func getSettingValue(name string) string {
	var s setting
	err := db.QueryRowx(`SELECT * FROM settings WHERE name = $1`, name).StructScan(&s)
	if err != nil {
		logger.Error("get settings value", "err", err)
		return ""
	}
	return s.Value
}

func AnonEditing() bool {
	return getSettingValue("anonediting") == "true"
}

func SetAnonEditing(setting bool) error {
	_, err := db.Exec(`UPDATE settings SET value = $1 WHERE name=$2`,
		strconv.FormatBool(setting), "anonediting")
	return err
}

func Signups() bool {
	return getSettingValue("allowsignups") == "true"
}

func SetSignups(setting bool) error {
	_, err := db.Exec(`UPDATE settings SET value = $1 WHERE name=$2`,
		strconv.FormatBool(setting), "allowsignups")
	return err
}

func SessionSecret() string {
	return getSettingValue("sessionsecret")
}
