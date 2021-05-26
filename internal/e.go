package internal

type E struct {
	Err    error
	Status int
}

func (this E) Error() string {
	return this.Err.Error()
}
