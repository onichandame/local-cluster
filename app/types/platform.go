package types

type Platform string

const (
	LINUX Platform = "linux"
)

func (*Platform) IsValid(raw Platform) bool {
	switch raw {
	case LINUX:
		return true
	}
	return false
}
