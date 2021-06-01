package model

type User struct {
	ID             int
	Username       string
	HashedPassword []byte
}

type Jwt struct {
	ID int
}

func (this Jwt) Valid() error {
	return nil
}
