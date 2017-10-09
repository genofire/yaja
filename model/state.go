package model

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type State struct {
	Domains map[string]*Domain
	sync.Mutex
}

func ReadState(path string) (state *State, err error) {
	state = &State{}

	if f, err := os.Open(path); err == nil { // transform data to legacy meshviewer
		if err = json.NewDecoder(f).Decode(state); err == nil {
			return state, nil

		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *State) UpdateDomain(d *Domain) error {
	if d.FQDN == "" {
		return errors.New("No fqdn exists in domain")
	}
	s.Lock()
	s.Domains[d.FQDN] = d
	s.Unlock()
	return nil
}
