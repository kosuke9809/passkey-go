package user

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID          []byte `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Cred        []webauthn.Credential
	sessionData webauthn.SessionData
}

func (u *User) WebAuthnID() []byte { return u.ID }

func (u *User) WebAuthnName() string { return u.Name }

func (u *User) WebAuthnDisplayName() string { return u.DisplayName }

func (u *User) WebAuthnCred() []webauthn.Credential { return u.Cred }

func (u *User) SessionData() webauthn.SessionData { return u.sessionData }

func (u *User) SetSessionData(data webauthn.SessionData) {
	u.sessionData = data
}

func New(username string) *User {
	return &User{
		ID:          []byte(username),
		Name:        username,
		DisplayName: username,
	}
}
