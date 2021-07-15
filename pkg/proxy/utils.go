package proxy

import "net"

func parseTCPAddr(raw *string) error {
	res, err := net.ResolveTCPAddr("tcp", *raw)
	if err == nil {
		*raw = res.String()
	}
	return err
}
