package application

type ApplicationType string

const (
	LOCAL  ApplicationType = "LOCAL"
	STATIC ApplicationType = "STATIC"
	REMOTE ApplicationType = "REMOTE"
)
