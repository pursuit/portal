package internal

const EInvalidInput = 422_000
const EMutationRefNotFound = 422_001
const EDbProblem = 503_001
const EFailedHashPassword = 503_002
const EFailedStoreEvent = 503_003

type E struct {
	Err    error
	Status int
}

func (this E) Error() string {
	return this.Err.Error()
}
