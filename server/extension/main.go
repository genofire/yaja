package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/server/utils"
)

type Extension interface {
	Process(*xml.StartElement, *utils.Client) bool
}
