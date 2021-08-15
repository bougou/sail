package models

type Target struct {
	Name     string
	PlayName string
}

type Common struct {
	SSHPort string
	SSHUser string
	SSHPass string

	SudoNopass bool

	InstallDir string
	DataDir    string
}
