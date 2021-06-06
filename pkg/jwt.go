package pkg

type Jwt struct {
	ID int `json:"id"`
}

func (this Jwt) Valid() error {
	return nil
}
