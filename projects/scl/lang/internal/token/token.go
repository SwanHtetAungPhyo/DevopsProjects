package token

const (
	SETTING   = "setting"
	TARGET    = "target"
	SUPERUSER = "superuser"
)

const (
	SetConfig = "configuration"
)
const (
	ONERROR  = iota
	ROLLBACK = iota
)

func GetAllowedSettingKeys() []string {
	return []string{
		SetConfig,
	}
}
