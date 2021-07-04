package types

type CompressType string

const (
	ZIP CompressType = "zip"
	TGZ CompressType = "tgz"
)

func (ct *CompressType) IsValid(raw CompressType) bool {
	switch raw {
	case ZIP, TGZ:
		return true
	}
	return false
}
