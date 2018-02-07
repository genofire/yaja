package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/server/utils"
)

type Extensions []Extension

type Extension interface {
	Process(*xml.StartElement, *utils.Client) bool
	Spaces() []string
}

func (ex Extensions) Spaces() (result []string) {
	for _, extension := range ex {
		result = append(result, extension.Spaces()...)
	}
	return result
}

func (ex Extensions) Process(element *xml.StartElement, client *utils.Client) {
	log := client.Log.WithField("extension", "all")

	// run every extensions
	count := 0
	for _, extension := range ex {
		if extension.Process(element, client) {
			count++
		}
	}

	// not extensions found
	if count != 1 {
		log.Debug(element)
	}
}
