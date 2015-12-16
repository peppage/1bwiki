package model

type setting struct {
	Name  string
	Value string
}

func getSettingValue(name string) string {
	var s setting
	err := db.QueryRowx(`SELECT * FROM settings WHERE name = $1`, name).StructScan(&s)
	if err != nil {
		logger.Error("verify user db error", "err", err)
		return ""
	}
	return s.Value
}

func AnonEditing() bool {
	if getSettingValue("anonediting") == "true" {
		return true
	}
	return false
}

func Signups() bool {
	if getSettingValue("allowsignups") == "true" {
		return true
	}
	return false
}
