// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package jwtauth provides secure user payload signing.
package jwtauth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// ctxAuthKey is an key for context values.
type ctxAuthKey struct{}

var (
	errInit = errors.New("secret cannot be empty") // error for invalid Init()
	ctxKey  = &ctxAuthKey{}                        // ctxKey is a context key instance.
)

// UserClass indicates the access control class a user is registered to.
type UserClass string

const (
	// UserSuperuser is a user that has full root privileges
	UserSuperuser UserClass = "superuser"
	// UserAdmin is a user that has admin privileges
	UserAdmin UserClass = "admin"
	// UserStandard is a user that has regular privileges
	UserStandard UserClass = "standard"
)

var (
	// UserClassMap maps the string representation back to UserClass
	UserClassMap = map[string]UserClass{
		"superuser": UserSuperuser,
		"admin":     UserAdmin,
		"standard":  UserStandard,
	}
)

// Config contains the state (secret and key) function for signing and
// validating tokens.
type Config struct {
	secret  []byte                                // secret for signing
	keyFunc func(*jwt.Token) (interface{}, error) // key function
}

// Payload represents an authenticated user. All fields are application
// specific.
type Payload struct {
	UUID      string    `json:"uuid"`      // UUID is unique user identifier (email)
	Class     UserClass `json:"class"`     // Class is user class
	Name      string    `json:"name"`      // Name is the user's name
	Group     string    `json:"group"`     // Group is the user's group
	Activated bool      `json:"activated"` // Activated if account is active
	jwt.StandardClaims
}

// Init takes a secret and returns a new state or an error.
func Init(secret []byte) (*Config, error) {
	// disallow empty secret
	if bytes.Equal(secret, nil) {
		return nil, errInit
	}
	// keyFunc is required for JWT library
	config := &Config{
		secret: secret,
		keyFunc: func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	}
	return config, nil
}

// CreateToken takes a payload and expiry duration, returning a signed base64
// Payload token or an error.
func (c *Config) CreateToken(p *Payload, duration time.Duration) (string, error) {
	// set expiry time of now + duration (e.g. 5 minutes)
	p.ExpiresAt = time.Now().Add(duration).Unix()
	// generate token using HMAC with SHA-512
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, p)
	// sign string with secret, if error return error, else return token
	signedStr, err := token.SignedString(c.secret)
	if err != nil {
		return "", err
	}
	// convert the signed string to base64
	signedBytes := []byte(signedStr)
	str := base64.StdEncoding.EncodeToString(signedBytes)
	return str, nil
}

// ParseToken takes a base64 token string, validating the string and returning a
// Payload object. If the token string is expired or not valid, an error is
// returned.
func (c *Config) ParseToken(str string) (*Payload, error) {
	// convert base64 to string
	tokenBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	tokenStr := string(tokenBytes)
	// parse token string
	token, err := jwt.ParseWithClaims(tokenStr, &Payload{}, c.keyFunc)
	if err != nil {
		return nil, err
	}
	// generate payload from parsed token when validating
	if payload, ok := token.Claims.(*Payload); ok && token.Valid {
		return payload, nil
	}
	// not ok, parse could not parse/validate the string
	return nil, err
}

// Context injects the payload into the provided context, returning the new
// context.
func (p *Payload) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, p)
}

// FromContext returns the payload from context. If a Payload was not placed
// into context, an empty Payload will be returned.
func FromContext(ctx context.Context) *Payload {
	p, ok := ctx.Value(ctxKey).(*Payload)
	if !ok || p == nil {
		return &Payload{}
	}
	return p
}

// GenerateSeed will return a random byte array of the specified input length,
// or an error.
func GenerateSeed(length int) ([]byte, error) {
	// generate cryptographic random byte array
	buff := make([]byte, length)
	_, err := rand.Read(buff)
	if err != nil {
		// error generating buffer
		return nil, err
	}
	return buff, nil
}

// HashPassword takes a string and returns a salted hash or an error.
func HashPassword(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash := string(hashedBytes[:])
	return hash, nil
}

// HashCompare takes the hashed input and a string, returning true if the hash
// is equivalent to the plaintext.
func HashCompare(hash string, s string) bool {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming) == nil
}
