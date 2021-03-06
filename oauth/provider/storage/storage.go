package storage

import (
	"github.com/RangelReale/osin"
	"github.com/ory-am/common/rand/sequence"
	"net/url"
	"time"
)

type Storage interface {
	SaveStateData(stateData *StateData) error
	GetStateData(id string) (*StateData, error)
}

type StateData struct {
	ID          string
	ClientID    string
	RedirectURL string
	Scope       string
	State       string
	Type        string
	Provider    string
	ExpiresAt   time.Time
}

func (s *StateData) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *StateData) FromAuthorizeRequest(ar *osin.AuthorizeRequest, provider string) error {
	if seq, err := sequence.RuneSequence(25, sequence.AlphaNum); err != nil {
		return err
	} else {
		s.ID = string(seq)
	}
	s.ClientID = ar.Client.GetId()
	s.RedirectURL = ar.RedirectUri
	s.Scope = ar.Scope
	s.State = ar.State
	s.Type = string(ar.Type)
	s.Provider = provider
	return nil
}

func (s *StateData) ExpireInOneHour() {
	s.ExpiresAt = time.Now().Add(time.Hour)
}

func (s *StateData) ToURLValues() (v url.Values) {
	v = make(url.Values)
	v.Set("state", s.State)
	v.Set("response_type", s.Type)
	v.Set("client_id", s.ClientID)
	v.Set("redirect_uri", s.RedirectURL)
	v.Set("scope", s.Scope)
	v.Set("provider", s.Provider)
	return v
}
