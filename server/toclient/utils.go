package toclient

import (
	"crypto/rand"
	"fmt"
)

func makeResource() string {
	var buf [16]byte
	if _, err := rand.Reader.Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return fmt.Sprintf("%x", buf)
}
