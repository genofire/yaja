package model

import (
	"bytes"
	"encoding/xml"
)

func XMLEscape(s string) string {
	var b bytes.Buffer
	xml.Escape(&b, []byte(s))

	return b.String()
}
