package utils

import (
	"dev.sum7.eu/genofire/yaja/model"
)

type DomainRegisterAllowed func(*model.JID) bool
