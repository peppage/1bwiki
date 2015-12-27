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

func saveSetting(name string, value string) error {
	_, err := db.Exec(`UPDATE settings SET value = $1 WHERE name=$2`, value, name)
	return err
}

func AnonEditing() bool {
	return getSettingValue("anonediting") == "true"
}

func SetAnonEditing(setting bool) error {
	return saveSetting("anonediting", strconv.FormatBool(setting))
}

func Signups() bool {
	return getSettingValue("allowsignups") == "true"
}

func SetSignups(setting bool) error {
	return saveSetting("allowsignups", strconv.FormatBool(setting))
}

func SessionSecret() string {
	return getSettingValue("sessionsecret")
}
