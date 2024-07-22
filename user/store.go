package user

import "errors"

var userDB = make(map[string]*User)

func Store(u *User) {
	userDB[u.Name] = u
}

func Get(username string) (*User, error) {
	u, ok := userDB[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}
