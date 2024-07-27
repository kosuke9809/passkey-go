package auth

import "github.com/go-webauthn/webauthn/webauthn"

var Web *webauthn.WebAuthn

func InitWebAuthn() error {
	var err error
	Web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "WebAuthn Go",
		RPID:          "localhost",
		RPOrigin:      "http://localhost:8082",
		Debug:         true,
	})
	return err
}
