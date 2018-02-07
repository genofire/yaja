package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	"dev.sum7.eu/genofire/yaja/model"
)

// Cookie is used to give a unique identifier to each request.
type Cookie uint64

func CreateCookie() Cookie {
	var buf [8]byte
	if _, err := rand.Reader.Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return Cookie(binary.LittleEndian.Uint64(buf[:]))
}
func CreateCookieString() string {
	return fmt.Sprintf("%x", CreateCookie())
}

type DomainRegisterAllowed func(*model.JID) bool
