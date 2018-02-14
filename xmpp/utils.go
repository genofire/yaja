package xmpp

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"reflect"
)

func XMLStartElementToString(element *xml.StartElement) string {
	if element == nil {
		return "<nil>"
	}
	debug := fmt.Sprintf("<%s xmlns=\"%s\"", element.Name.Local, element.Name.Space)
	for _, attr := range element.Attr {
		debug = fmt.Sprintf("%s %s=\"%s\"", debug, attr.Name.Local, attr.Value)
	}
	debug += ">"
	return debug
}

func XMLChildrenString(o interface{}) (result string) {
	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Struct {
		first := true
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
				elm := valueField.Elem()
				if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
					valueField = elm
				}
			}

			if xmlElement, ok := valueField.Interface().(*xml.Name); ok && xmlElement != nil {
				if first {
					first = false
					result += xmlElement.Local
				} else {
					result += ", " + xmlElement.Local
				}
			}
		}
	}
	return
}

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
