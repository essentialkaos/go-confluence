package confluence

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/base64"
	"errors"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Auth is interface for authorization method
type Auth interface {
	Validate() error
	Encode() string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AuthBasic is struct with data for basic authorization
type AuthBasic struct {
	User     string
	Password string
}

// AuthToken is struct with data for personal token authorization
type AuthToken struct {
	Token string
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrEmptyUser        = errors.New("User can't be empty")
	ErrEmptyPassword    = errors.New("Password can't be empty")
	ErrEmptyToken       = errors.New("Token can't be empty")
	ErrTokenWrongLength = errors.New("Token length must be equal to 44")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Validate validates authorization data
func (a AuthBasic) Validate() error {
	switch {
	case a.User == "":
		return ErrEmptyUser
	case a.Password == "":
		return ErrEmptyPassword
	}

	return nil
}

// Encode encodes data for authorization
func (a AuthBasic) Encode() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(a.User+":"+a.Password))
}

// Validate validates authorization data
func (a AuthToken) Validate() error {
	switch {
	case a.Token == "":
		return ErrEmptyToken
	case len(a.Token) != 44:
		return ErrTokenWrongLength
	}

	return nil
}

// Encode encodes data for authorization
func (a AuthToken) Encode() string {
	return "Bearer " + a.Token
}
