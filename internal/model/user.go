package model

type User struct {
	ID             int
	Username       string
	HashedPassword []byte
}
